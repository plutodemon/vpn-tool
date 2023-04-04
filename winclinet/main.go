package main

import (
	"github.com/sirupsen/logrus"
	"ivs-net-winclinet/ui"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	//test1@sina.com
	//图形界面
	ui.UI()
	//建立连接链路
	//_, body := ui.GetIp()
	//netLink.Net(body)
	//net.Net()

}
