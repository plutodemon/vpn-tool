package https

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"ivs-net-server/auth/models"
	"net"
	"net/http"
	"sync"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var Regular *cron.Cron = nil
var ws *websocket.Conn = nil
var wsC *websocket.Conn = nil

// Speed 服务端速率
type Speed struct {
	UserNum     int     `json:"userNum"`
	ClientNum   uint32  `json:"clientNum"`
	ClientNumSp uint32  `json:"clientNumSp"`
	Pbps        uint32  `json:"pbps"`
	Kbps        uint64  `json:"kbps"`
	CpuPer      float64 `json:"cpuPer"`
	MemPer      float64 `json:"memPer"`
	DiskPer     float64 `json:"diskPer"`
}

var SerSpeed Speed
var numSp uint32

func GetSpeed(c *gin.Context) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		ws, _ = upGrader.Upgrade(c.Writer, c.Request, nil)
		Regular = cron.New()
		_, err := Regular.AddFunc("@every 1s", func() {
			SerSpeed.CpuPer = GetCpuPercent()
			SerSpeed.MemPer = GetMemPercent()
			SerSpeed.DiskPer = GetDiskPercent()
			SerSpeed.ClientNumSp = SerSpeed.ClientNum - numSp
			SerSpeed.UserNum = userNum
			message, _ := json.Marshal(SerSpeed)
			//写入ws数据
			err := ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("写入失败：", err)
				return
			}
			numSp = SerSpeed.ClientNum
			SerSpeed.Pbps = 0
			SerSpeed.Kbps = 0
		})
		if err != nil {
			fmt.Println("定时器失败：", err)
			return
		}
		Regular.Start()
	}
}
func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

var userNum int

func init() {
	numSp = 0
	SerSpeed.ClientNum = 0
	SerSpeed.Pbps = 0
	SerSpeed.Kbps = 0
	var user models.User
	userNum = len(user.FindUi())
	Clients.Client = make(map[uint32]*ClientSpeed, userNum)
	UpSort = make([]uint32, 0)
}

type ClientSpeed struct {
	Lose      uint
	Email     string
	Id        uint
	Ip        uint32
	UserLevel uint8
	Pbps      uint32
	Kbps      uint64
}
type CliSpeed struct {
	Ip   net.IP `json:"ip"`
	Pbps uint32 `json:"pbps"`
	Kbps uint64 `json:"kbps"`
}

var UpSort []uint32

var Clients SafeClient

type SafeClient struct {
	sync.RWMutex
	Client map[uint32]*ClientSpeed
}

func (m *SafeClient) Add(item *ClientSpeed) {
	m.Lock()
	defer m.Unlock()
	m.Client[item.Ip] = item
}

func (m *SafeClient) Get(ip uint32) *ClientSpeed {
	m.RLock()
	defer m.RUnlock()
	return m.Client[ip]
}
func (m *SafeClient) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.Client)
}

func (m *SafeClient) Del(ip uint32) {
	m.Lock()
	defer m.Unlock()
	delete(m.Client, ip)
}

func GetClientSpeed(c *gin.Context) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		wsC, _ = upGrader.Upgrade(c.Writer, c.Request, nil)
		Regular = cron.New()
		_, err := Regular.AddFunc("@every 3s", func() {

			//写入ws数据
			f := 0
			clients := make([]CliSpeed, Clients.Len())
			for _, ip := range UpSort {
				clients[f].Ip = UInt32ToIP(Clients.Get(ip).Ip)
				clients[f].Pbps = Clients.Get(ip).Pbps
				clients[f].Kbps = Clients.Get(ip).Kbps
				f++
			}
			message, _ := json.Marshal(clients)
			err := wsC.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("写入失败：", err)
				return
			}
			for i, _ := range Clients.Client {
				get := Clients.Get(i)
				get.Pbps = 0
				get.Kbps = 0
				Clients.Add(get)
			}
		})

		if err != nil {
			fmt.Println("定时器失败：", err)
			return
		}
		Regular.Start()
	}
}
func GetClientSpeeds(c *gin.Context) {
	f := 0
	clients := make([]CliSpeed, Clients.Len())
	for _, ip := range UpSort {
		clients[f].Ip = UInt32ToIP(Clients.Get(ip).Ip)
		clients[f].Pbps = Clients.Get(ip).Pbps
		clients[f].Kbps = Clients.Get(ip).Kbps
		f++
	}
	c.JSON(200, clients)
}
