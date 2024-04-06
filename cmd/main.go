// go run ./cmd
// go run -tags prime ./cmd
// tinygo flash -target xxx ./cmd

package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/runner"
	"github.com/merliot/temp"
)

var (
	id           = dean.GetEnv("ID", "temp01")
	name         = dean.GetEnv("NAME", "Temp/Humidity")
	deployParams = dean.GetEnv("DEPLOY_PARAMS", "target=nano-rp2040&gpio=D2")
	wsScheme     = dean.GetEnv("WS_SCHEME", "ws://")
	port         = dean.GetEnv("PORT", "8000")
	portPrime    = dean.GetEnv("PORT_PRIME", "8001")
	user         = dean.GetEnv("USER", "")
	passwd       = dean.GetEnv("PASSWD", "")
	dialURLs     = dean.GetEnv("DIAL_URLS", "")
	ssids        = dean.GetEnv("WIFI_SSIDS", "")
	passphrases  = dean.GetEnv("WIFI_PASSPHRASES", "")
)

func main() {
	temp := temp.New(id, "temp", name).(*temp.Temp)
	temp.SetDeployParams(deployParams)
	temp.SetWifiAuth(ssids, passphrases)
	temp.SetDialURLs(dialURLs)
	temp.SetWsScheme(wsScheme)
	runner.Run(temp, port, portPrime, user, passwd, dialURLs, wsScheme)
}
