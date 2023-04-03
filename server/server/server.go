package server

import (
	"crypto/rc4"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"ivs-net-server/auth/https"
	"strconv"
	"sync"
)

// tcp
type tcpServer struct {
	*gnet.EventServer
}

// 客户端连接信息
type ipKey struct {
	c      gnet.Conn
	cipher *rc4.Cipher
}
type cKey struct {
	ip     uint32
	cipher *rc4.Cipher
}

// OnOpened 连接创建
func (es *tcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	dhcp := https.DhcpGet
	var bytes = make([]byte, 64)
	for i, b := range https.Rc4Key {
		bytes[i] = b
	}
	cipher, _ := rc4.NewCipher(bytes)
	ik := ipKey{
		c:      c,
		cipher: cipher,
	}
	session.Set(strconv.Itoa(int(dhcp.Ip)), ik, cache.NoExpiration)

	cip, _ := rc4.NewCipher(bytes)
	ck := cKey{
		ip:     dhcp.Ip,
		cipher: cip,
	}
	c.SetContext(ck)
	session.Set(fmt.Sprintf("%v", c.Context()), ck, cache.NoExpiration)
	https.SerSpeed.ClientNum++
	fmt.Println("已连接：", dhcp.Ip)
	return
}

// OnClosed 断开
func (es *tcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	err = c.Close()
	if err != nil {
		return 0
	}
	session.Delete(fmt.Sprintf("%v", c.Context()))
	https.SerSpeed.ClientNum--
	fmt.Println(c.Context().(cKey).ip, "断开连接！")
	return
}

func (es *tcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	//解密
	get, b := session.Get(fmt.Sprintf("%v", c.Context()))
	if b {
		get.(cKey).cipher.XORKeyStream(frame, frame)
	} else {
		fmt.Println("发送端c没找到")
	}
	/*
		id := binary.BigEndian.Uint64(frame[8:])
		mask := binary.BigEndian.Uint64(frame[16:])
		src_ip := binary.BigEndian.Uint64(frame[24:])
	*/

	//fmt.Println("原本发送的：", frame)

	ip := uint32(frame[48])<<24 | uint32(frame[49])<<16 | uint32(frame[50])<<8 | uint32(frame[51])

	cli := https.Clients.Get(get.(cKey).ip)
	cli.Pbps++
	cli.Kbps += uint64(len(frame))
	https.Clients.Add(cli)

	https.SerSpeed.Pbps++
	https.SerSpeed.Kbps += uint64(len(frame))

	conn, flag := session.Get(strconv.Itoa(int(ip)))
	if flag {
		client := https.Clients.Get(ip)

		if client.UserLevel == 0 {
			if client.Lose > 5 {
				client.Lose = 0
				https.Clients.Add(client)
				return nil, 0
			} else {
				client.Lose++
			}
		}
		out1 := frame[32:]
		conn.(ipKey).cipher.XORKeyStream(out1, out1)

		https.SerSpeed.Pbps++
		https.SerSpeed.Kbps += uint64(len(out1))

		client.Pbps++
		client.Kbps += uint64(len(out1))
		https.Clients.Add(client)

		err := conn.(ipKey).c.AsyncWrite(out1)
		if err != nil {
			fmt.Println("转发发生错误：", err)
			return nil, 0
		}
	} else {
		fmt.Println("目标ip未在线")
	}
	return
}

//var list *SkipList

var session *cache.Cache

func Server(WG *sync.WaitGroup) {
	defer WG.Done()
	session = cache.New(cache.NoExpiration, cache.NoExpiration)
	echo := new(tcpServer)
	err := gnet.Serve(echo, "tcp://0.0.0.0:9000", gnet.WithMulticore(true))
	if err != nil {
		logrus.Error("gnet启动失败：", err)
		return
	}
}
