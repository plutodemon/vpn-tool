package main

import (
	"github.com/sirupsen/logrus"
	_ "ivs-net-server/auth/configure"
	"ivs-net-server/auth/https"
	_ "ivs-net-server/auth/models"
	"ivs-net-server/server"
	"sync"
)

var WG sync.WaitGroup

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	WG.Add(2)
	go https.RunHttpsServer(&WG)
	go server.Server(&WG)
	WG.Wait()
}
