package sshocks

import (
	"encoding/json"

	"github.com/audibleblink/gorsh/internal/utils"

	"github.com/audibleblink/HoleySocks/pkg/holeysocks"
)

var (
	Config holeysocks.MainConfig
)

func init() {

	configBytes, err := utils.GetBytes(configs, "ssh.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configBytes, &Config)
	if err != nil {
		panic(err)
	}

	privKeyBytes, err := utils.GetBytes(configs, "id_ed25519")
	if err != nil {
		panic(err)
	}
	Config.SSH.SetKey(privKeyBytes)
}
