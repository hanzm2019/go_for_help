package affect

import (
	"bytes"
	"encoding/binary"
	"github.com/blind-oracle/go-radius"
	"go_for_help/redisUtil"
	"nac/concurrence"
	"nac/nacConfig"
	"net"
	"runtime"
	"strconv"
	"strings"
)

var (
	//srcIpAddr                 string
	//disconnectDstIpAddr       string
	//serverIpAddr              string

	reauthDstIpAddr           string //重认证，发送请求给AC是AC的ip+重连接端口1700
	secret                    string //AC 和服务端的共享秘钥
	acIp                      string //AC 的ip
	ciscoReauthenticateAVPair string //思科 重连接 ACL
	ciscoLastAVPair           string //思科 重连接 ACL

	preAuthAcl    string // 定义的 认证重定向 ACL
	redirectUrl   string // 定义的 认证重定向 url
	permissionAcl string // 定义的 放行 ACL
	letGoMAcs     string // Mac 白名单

	preAuthAttribute              radius.Attribute // 封装好的 认证重定向ACL
	permissionAttribute           radius.Attribute // 封装好的 放行 ACL
	ciscoReauthenticateAVPairAttr radius.Attribute // 封装好的 重连接ACL
	ciscoLastAVPairAttr           radius.Attribute // 封装好的 重连接ACL
	newIpAttr                     radius.Attribute // 封装好的 AC的ip

	expireTime int //会话超时时间

	workerPool *concurrence.WorkerPool
)

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	config := nacConfig.GetConfig()
	nac := config.Nac

	redirectUrl = nac.RedirectUrl
	reauthDstIpAddr = nac.ReauthDstIpAddr
	acIp = nac.AcIp
	secret = nac.Secret
	//serverIpAddr = nac.ServerIpAddr
	ciscoReauthenticateAVPair = nac.CiscoReauthenticateAVPair
	ciscoLastAVPair = nac.CiscoLastAVPair

	//srcIpAddr = nac.SrcIpAddr
	//disconnectDstIpAddr = nac.DisconnectDstIpAddr

	preAuthAcl = nac.PreAuthAcl
	permissionAcl = nac.PermissionAcl
	letGoMAcs = nac.LetGoMAcs

	expireTime = config.TimeConfig.ExpireTime

	preAuthAttribute = EncodeAVPairByte(9, 1, []byte(preAuthAcl))
	permissionAttribute = EncodeAVPairByte(9, 1, []byte(permissionAcl))

	ciscoReauthenticateAVPairAttr = EncodeAVPairByte(9, 1, []byte(ciscoReauthenticateAVPair))
	ciscoLastAVPairAttr = EncodeAVPairByte(9, 1, []byte(ciscoLastAVPair))

	attrAcIp := strToNetIp(acIp)
	newIpAttr, _ = radius.NewIPAddr(attrAcIp)

}

func strToNetIp(ip string) net.IP {
	ips := strings.Split(ip, ".")
	ip0, err := strconv.Atoi(ips[0])
	if err != nil {
		return nil
	}
	ip1, err := strconv.Atoi(ips[1])
	if err != nil {
		return nil
	}
	ip2, err := strconv.Atoi(ips[2])
	if err != nil {
		return nil
	}
	ip3, err := strconv.Atoi(ips[3])
	if err != nil {
		return nil
	}
	netIP := net.IPv4(byte(ip0), byte(ip1), byte(ip2), byte(ip3))
	return netIP
}

func EncodeAVPairByte(vendorID uint32, typeID uint8, value []byte) (vsa []byte) {
	var b bytes.Buffer
	bv := make([]byte, 4)
	binary.BigEndian.PutUint32(bv, vendorID)

	// Vendor-Id(4) + Type-ID(1) + Length(1)
	b.Write(bv)
	b.Write([]byte{byte(typeID), byte(len(value) + 2)})

	// Append attribute value pair
	b.Write(value)

	vsa = b.Bytes()
	return
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
	}
	isExist := redisUtil.IsUserSessionExist(mac)
	if isExist {
		return 1
	} else {
		return 0
	}

}
