package sshocks

import (
	_ "embed"
	"encoding/json"

	"github.com/audibleblink/holeysocks/pkg/holeysocks"
)

var (
	Config holeysocks.MainConfig
)

//go:embed conf/ssh.json
var configBytes []byte

//go:embed conf/id_ed25519
var privKeyBytes []byte

func init() {
	err := json.Unmarshal(configBytes, &Config)
	if err != nil {
		panic(err)
	}

	Config.SSH.SetKey(privKeyBytes)
}
