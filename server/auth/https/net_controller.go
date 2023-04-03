package https

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ivs-net-server/auth/models"
	"net"
	"strconv"
	"strings"
)

type CreatBody struct {
	Ip         net.IP `json:"ip"`
	NetMask    net.IP `json:"netMask"`
	NetGateway net.IP `json:"netGateway"`
	DnsAddress net.IP `json:"dnsAddress"`
	MaxAllow   uint64 `json:"maxAllow"`
}

// IPToUInt32 转uint32
func IPToUInt32(ip net.IP) uint32 {
	bits := strings.Split(ip.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

// CreateNet 创建网络
func CreateNet(c *gin.Context) {
	var creatBody CreatBody
	var netWork models.NetWork
	err := c.ShouldBindJSON(&creatBody)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "失败",
			"data":    creatBody,
			"error":   err,
		})
	} else {
		fmt.Println(creatBody)
		netWork.Ip = IPToUInt32(creatBody.Ip)
		netWork.Mask = IPToUInt32(creatBody.NetMask)
		netWork.NetGateway = IPToUInt32(creatBody.NetGateway)
		netWork.DnsAddress = IPToUInt32(creatBody.DnsAddress)
		netWork.MaxAllow = creatBody.MaxAllow
		netWork.CreateNet()
		c.JSON(200, gin.H{
			"message": "成功",
			"data":    creatBody,
		})
	}
}
func UInt32ToIP(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// FindNet 查询
func FindNet(c *gin.Context) {
	var netWork models.NetWork
	nets := netWork.FindNet()
	netBody := make([]NetBody, len(nets))
	for i := 0; i < len(nets); i++ {
		netBody[i].Id = nets[i].ID
		netBody[i].Ip = UInt32ToIP(nets[i].Ip)
		netBody[i].Mask = UInt32ToIP(nets[i].Mask)
		netBody[i].NetGateway = UInt32ToIP(nets[i].NetGateway)
		netBody[i].DnsAddress = UInt32ToIP(nets[i].NetGateway)
		netBody[i].MaxAllow = strconv.FormatUint(nets[i].MaxAllow, 10)
	}
	c.JSON(200, netBody)
}

// DeleteNet 删除
func DeleteNet(c *gin.Context) {
	var netWork models.NetWork
	netWork.DeleteNet(c.Query("id"))
	c.JSON(200, gin.H{"mes": "success"})
}

// NetBody 更新结构体
type NetBody struct {
	Id uint `json:"ID"`
	//ip
	Ip net.IP `json:"ip"`
	//掩码
	Mask net.IP `json:"mask"`
	//网关
	NetGateway net.IP `json:"netGateway"`
	//dns服务器地址
	DnsAddress net.IP `json:"dnsAddress"`
	//最大允许设备数
	MaxAllow string `json:"maxAllow"`
}

// EditNet 更新
func EditNet(c *gin.Context) {
	var netWork models.NetWork
	var netBody NetBody

	err := c.ShouldBindJSON(&netBody)
	netWork.ID = netBody.Id
	netWork.Ip = IPToUInt32(netBody.Ip)
	netWork.Mask = IPToUInt32(netBody.Mask)
	netWork.NetGateway = IPToUInt32(netBody.NetGateway)
	netWork.DnsAddress = IPToUInt32(netBody.DnsAddress)
	max, _ := strconv.ParseUint(netBody.MaxAllow, 10, 8)
	netWork.MaxAllow = max

	if err != nil {
		c.JSON(400, gin.H{
			"mes": "error",
		})

	} else {
		netWork.UpdateNet()
		c.JSON(200, gin.H{
			"mes": "success",
		})
	}
}
