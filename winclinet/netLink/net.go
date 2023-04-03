package netLink

import (
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
	"golang.zx2c4.com/wireguard/tun"
	"ivs-net-winclinet/configure"
	"ivs-net-winclinet/login"
	"ivs-net-winclinet/winipcfg"
	"net"
	"net/netip"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// Net 连接
var link winipcfg.LUID
var ip netip.Prefix
var dev tun.Device
var err error
var conn net.Conn
var dhcpBody login.DhcpBody
var address string

func Net(body login.DhcpBody) {
	dhcpBody = body
	output, _ := exec.Command("ipconfig").Output()
	ifName := "MyTun"
	if strings.Contains(string(output), ifName) {
		logrus.Println("创建网卡失败，已包含此网卡")
		return
	}
	dev, err = tun.CreateTUN(ifName, 0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// 保存原始设备句柄
	nativeTunDevice := dev.(*tun.NativeTun)

	// 获取LUID用于配置网络
	link = winipcfg.LUID(nativeTunDevice.LUID())

	toIp := Long2IPString(body.Ip)
	mask := find32One(int(body.Mask))
	ipAddress := toIp + "/" + strconv.Itoa(mask)

	ip, err = netip.ParsePrefix(ipAddress)
	if err != nil {
		panic(err)
	}
	err = link.SetIPAddresses([]netip.Prefix{ip})
	if err != nil {
		panic(err)
	}

	address = configure.Config.Get("servers.tcp_address").(string) + ":" + configure.Config.Get("servers.tcp_port").(string)
	conn, err = net.Dial("tcp", address)
	//defer conn.Close()
	if err != nil {
		fmt.Printf("-----create conn failed, err:%v", err)
	}
	wg.Add(2)
	go Send(conn, dev, body)
	go Rec(conn, dev)
	wg.Wait()
}

// Send 发送
func Send(conn net.Conn, dev tun.Device, body login.DhcpBody) {
	defer wg.Done()
	n := 65535
	buf := make([]byte, n)

	mask := body.Mask
	srcIp := body.Ip
	ipNet := mask & srcIp

	//rc4加密
	bytes := make([]byte, 64)
	for i, b := range login.Rc4Key {
		bytes[i] = b
	}
	cipher, _ := rc4.NewCipher(bytes)
	var out []byte
	// 读取发送
	for {
		n, _ = dev.Read(buf[32:], 0)
		header, err := ipv4.ParseHeader(buf[32 : n+32])
		if err != nil {
			continue
		}
		netPart := uint32(buf[48])<<24 | uint32(buf[49])<<16 | uint32(buf[50])<<8 | uint32(buf[51])
		airPart := uint32(buf[48])<<24 | uint32(buf[49])<<16 | uint32(buf[50])<<8 | uint32(0xff)
		if ipNet == netPart&mask && netPart != airPart {
			timeNow := time.Now().Unix()
			binary.BigEndian.PutUint64(buf[0:], uint64(timeNow))
			binary.BigEndian.PutUint64(buf[8:], uint64(body.Id))
			binary.BigEndian.PutUint64(buf[16:], uint64(body.Ip))
			binary.BigEndian.PutUint64(buf[24:], uint64(body.Mask))
			logrus.Println("Src:", header.Src, " dst:", header.Dst, " 协议：", header.Protocol)

			// 加密后的数据直接覆盖到buf[:n+32]中
			out = buf[:n+32]
			cipher.XORKeyStream(out, out)
			//if header.Protocol == 6
			_, err = conn.Write(out)
			if err != nil {
				fmt.Println("写入数据失败：", err)
			}

		}
	}
}

var Em string
var IIp uint32

// Rec 接收
func Rec(conn net.Conn, dev tun.Device) {
	defer wg.Done()
	n := 65535
	buffer := make([]byte, n)

	var frame []byte
	//rc4解密
	bytes := make([]byte, 64)
	for i, b := range login.Rc4Key {
		bytes[i] = b
	}
	cipher, _ := rc4.NewCipher(bytes)

	for {
		/*reader := bufio.NewReader(conn)
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println("读取数据失败:", err)
			return
		}*/
		n, err = conn.Read(buffer)
		if err != nil {
			//断线重连
			//str := "An existing connection was forcibly closed by the remote host."
			fmt.Printf("read from conn failed, err:%v", err)
			//conn, err = net.Dial("tcp", address)
			//defer conn.Close()
			return
		}

		//解密
		frame = buffer[:n]
		cipher.XORKeyStream(frame, frame)

		//_, err = dev.Write(buf[32:], n)
		n, err := dev.Write(frame, 0)
		if err != nil {
			fmt.Println("写入网卡失败：", err)
			return
		}
		//fmt.Println("写入的数据：", ln)
		header, err := ipv4.ParseHeader(frame)
		if err != nil {
			continue
		}
		fmt.Println("这里接收到了：", header, "\n写入的长度：", n)
	}
}

func UnNet() {
	f := login.DisConn(dhcpBody)
	if f && conn != nil {
		err := conn.Close()
		if err != nil {
			fmt.Printf("conn close err:%v", err)
			return
		}
		err = link.DeleteIPAddress(ip)
		if err != nil {
			fmt.Printf("link close err:%v", err)
			return
		}
		err = dev.Close()
		if err != nil {
			fmt.Printf("dev close err:%v", err)
			return
		}
	} else {
		return
	}
}
func Long2IPString(i uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String()
}
func find32One(n int) int {
	MASK1 := 0x55555555
	MASK2 := 0x33333333
	MASK4 := 0x0f0f0f0f
	MASK8 := 0x00ff00ff
	MASK16 := 0x0000ffff

	n = (n & MASK1) + ((n >> 1) & MASK1)
	n = (n & MASK2) + ((n >> 2) & MASK2)
	n = (n & MASK4) + ((n >> 4) & MASK4)
	n = (n & MASK8) + ((n >> 8) & MASK8)
	n = (n & MASK16) + ((n >> 16) & MASK16)
	return n
}
