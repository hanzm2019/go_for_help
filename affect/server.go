package affect

import (
	"github.com/blind-oracle/go-radius"
	"nac/basicUtil"
	"nac/logger"
	"nac/redisUtil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	EXPIRE_TIME = 300
)

func Service() {
	logger.Info("radius服务开始启动....")

	server := radius.Server{
		Addr:       serverIpAddr,
		Handler:    radius.HandlerFunc(handleRadius),
		Secret:     []byte(secret),
		Dictionary: radius.Builtin,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error)

	logger.Info("waiting for packets...")

	if err := server.ListenAndServe(); err != nil {
		logger.Error("radius 服务启动失败...")
		logger.Error(err)
		errChan <- err
	}
	select {
	case <-signalChan:
		logger.Info("stopping server...")
		server.Close()
	case err := <-errChan:
		logger.Error("[ERR] %v", err.Error())
	}
}

func handleRadius(response radius.ResponseWriter, request *radius.Packet) {
	switch request.Code {
	case radius.CodeAccessRequest:
		handleAccessRequest(response, request)
	case radius.CodeAccountingRequest:
		response.AccountingACK()
	default:
		response.AccessReject()
	}
}

func handleAccessRequest(response radius.ResponseWriter, request *radius.Packet) {
	mac := request.String("User-Name")
	callingStationId := request.String("Called-Station-Id")

	if mac == "" || callingStationId == "" {
		logger.Error("radius连接请求中用户名/callingStationId为空")
		response.AccessReject()
		return
	}
	csis := strings.Split(callingStationId, ":")
	if len(csis) < 2 {
		logger.Error("radius请求中，callingstationId不含网络名称network")
		response.AccessReject()
		return
	}
	ssid_from := csis[0]

	status := IsUserOnLine(mac, ssid_from)
	uid := basicUtil.GetUid(mac)

	acl := &radius.Attribute{
		Type:  radius.AttrVendorSpecific,
		Value: radius.EncodeAVPairCisco(preAuthAcl),
	}

	redirect := &radius.Attribute{
		Type:  radius.AttrVendorSpecific,
		Value: radius.EncodeAVPairCisco(redirectUrl + "?uid=" + uid),
	}
	authpass := &radius.Attribute{
		Type:  radius.AttrVendorSpecific,
		Value: radius.EncodeAVPairCisco(permissionAcl),
	}
	switch status {
	case -1:
		logger.Error("[" + mac + "]信息操作异常...")
		response.AccessReject()
	case 0:
		logger.Info("[" + mac + "]-> status: 不在线(0)...")
		redisUtil.SetRedis(uid, mac)
		response.AccessAccept(acl, redirect)
	case 1:
		logger.Info("[" + mac + "]-> status: 在线(1)...")
		response.AccessAccept(authpass)
	default:
		logger.Info("[" + mac + "]-> status: 其他(default)...")
		response.AccessReject()
	}
}

/**
判断用户是否在线，先从redis中拿取用户信息，存在，获取用户状态
不存在，获取用户在线时间，如果在线时间获取失败，返回-1拒绝连接，
如果允许在线时间为-1，则为哑终端，始终在线，并将终端写入redis有限期无限制
如果允许在线时间为其他，查询用户，用户不存在，自定义用户信息写入redis状态为0，写入数据库
如果用户存在，获取用户状态，将用户写入redis中
如果添加用到数据失败，返回-1拒绝连接

*/
func IsUserOnLine(mac, callingStationId string) int {
	if mac == "" || callingStationId == "" {
		return -1
	}
	if callingStationId[:4] != "1#7F" { //非七楼的AP直接放行
		return 1
	}
	if strings.Contains(letGoMAcs, mac) {
		return 1
	} else {
		userMap := redisUtil.GetSessionIdByUser(mac)
		if userMap == "" {
			redisUtil.SetUserSessionIdWithExpire(mac, "0", EXPIRE_TIME)
			return 0
		} else if userMap == "0" {
			return 0
		} else {
			return 1
		}
	}
}
