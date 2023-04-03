package configure

import "github.com/pelletier/go-toml"

var Config *toml.Tree = nil

func init() {
	if Config == nil {
		Config, _ = toml.LoadFile("conf/app.toml")
	}
}
