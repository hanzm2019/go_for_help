package affect

import (
	"fmt"
	"github.com/blind-oracle/go-radius"
	"nac/dbutil"
	"nac/logger"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	AttrNASIPAddress       = 4
	AttrAcctTerminateCause = 49
)

func getNetIp(ip string) *net.IP {
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
	return &netIP
}

func getUDPAddr(ipAddr string) *net.UDPAddr {
	ipAddrs := strings.Split(ipAddr, ":")
	ip := ipAddrs[0]
	port, err := strconv.Atoi(ipAddrs[1])
	if err != nil {
		logger.Error(err)
		return nil
	}
	netIP := getNetIp(ip)
	udpaddr := net.UDPAddr{IP: *netIP, Port: port}
	return &udpaddr
}

/**
关闭设备连接wifi，radius服务器发送关闭packet给AC
*/
func DisconnectRequest(callingStationId, network string) bool {
	mac := strings.ReplaceAll(callingStationId, "-", "")
	logger.Info("mac:", mac)
	isOk := dbutil.UpdateUserOnlineStatus("mac", "network", mac, network, 0)
	if !isOk {
		return isOk
	}

	srcAddr := getUDPAddr(srcIpAddr)
	if srcAddr == nil {
		logger.Error("get local addr failed")
		return false
	}
	dstAddress := getUDPAddr(disconnectDstIpAddr)
	if dstAddress == nil {
		logger.Error("get src addr failed")
		return false
	}
	c := radius.Client{LocalAddr: srcAddr, Timeout: 1 * time.Second, Retries: 1}
	params := radius.RequestParams{Secret: []byte(secret), DstAddressPort: dstAddress, SrcAddress: srcAddr}
	dstAddrs := strings.Split(disconnectDstIpAddr, ":")
	nasIpAddress := getNetIp(dstAddrs[0])
	reply := c.CiscoRequest(
		&params,
		radius.CodeDisconnectRequest,
		&radius.Attribute{
			Type:  AttrAcctTerminateCause,
			Value: uint32(6),
		},
		&radius.Attribute{
			Type:  radius.AttrCallingStationID,
			Value: callingStationId,
		},
		&radius.Attribute{
			Type:  AttrNASIPAddress,
			Value: *nasIpAddress,
		},
	)

	return reply.Success
}

/**
连接设备wifi，radius服务器发送准入设备连接wifi的packet到AC
callingStationId : 5c-23-fe-de-4f-cd
*/
func ReAuthConnectRequest(callingStationId string) bool {

	if callingStationId == "" {
		return false
	}

	srcAddr := getUDPAddr(srcIpAddr)
	if srcAddr == nil {
		logger.Error("get local addr failed")
		return false
	}
	dstAddress := getUDPAddr(reauthDstIpAddr)
	if dstAddress == nil {
		logger.Error("get src addr failed")
		return false
	}
	c := radius.Client{LocalAddr: srcAddr, Timeout: 3 * time.Second, Retries: 6}
	param := radius.RequestParams{Secret: []byte(secret), DstAddressPort: dstAddress, SrcAddress: srcAddr}
	dstAddrs := strings.Split(reauthDstIpAddr, ":")
	nasIpAddress := getNetIp(dstAddrs[0])

	reply := c.CiscoRequest(
		&param,
		radius.CodeCoARequest,
		&radius.Attribute{
			Type:  AttrAcctTerminateCause,
			Value: uint32(6),
		},
		&radius.Attribute{
			Type:  radius.AttrCallingStationID,
			Value: callingStationId,
		},
		&radius.Attribute{
			Type:  AttrNASIPAddress,
			Value: *nasIpAddress,
		},
		&radius.Attribute{
			Type:  radius.AttrVendorSpecific,
			Value: radius.EncodeAVPairCisco(ciscoReauthenticateAVPair),
		},
		&radius.Attribute{
			Type:  radius.AttrVendorSpecific,
			Value: radius.EncodeAVPairCisco(ciscoLastAVPair),
		},
	)

	return reply.Success
}

func main() {
	//ReAuthConnectRequest()
	//reply := DisconnectRequest("d0-81-7a-b6-e8-56")
	reply := DisconnectRequest("d0-81-7a-b6-e8-56", "xmly-aqtest")
	fmt.Println(reply)
}
