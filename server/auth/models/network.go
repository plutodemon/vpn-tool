package models

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// NetWork 网络表
type NetWork struct {
	//网络号
	Ip uint32 `json:"ip"`
	//掩码
	Mask uint32 `json:"mask"`
	//网关
	NetGateway uint32 `json:"netGateway"`
	//dns服务器地址
	DnsAddress uint32 `json:"dnsAddress"`
	//最大允许设备数
	MaxAllow uint64 `json:"maxAllow" gorm:"default:0"`

	gorm.Model
}

// CreateNet 创建网络
func (n *NetWork) CreateNet() {
	DB.Create(n) // 通过数据的指针来创建
	fmt.Println(n)
	logrus.Println("新增网络成功:", n.ID)
	return
}

// FindNet 查找
func (n *NetWork) FindNet() (nets []NetWork) {
	DB.Model(n).Find(&nets)
	return
}

// FindNetById 查找网络
func (n *NetWork) FindNetById(id uint) {
	DB.First(n, id)
}

// DeleteNet 删除
func (n *NetWork) DeleteNet(id string) {
	DB.Delete(n, id)
}

// UpdateNet 更新
func (n *NetWork) UpdateNet() {
	DB.Model(n).Updates(map[string]interface{}{"Ip": n.Ip, "Mask": n.Mask, "NetGateway": n.NetGateway, "DnsAddress": n.DnsAddress, "MaxAllow": n.MaxAllow})
}

// DhcpInfo 随机的ip
func (n *NetWork) DhcpInfo(id int64) NetWork {
	DB.First(n, id)
	return *n
}
