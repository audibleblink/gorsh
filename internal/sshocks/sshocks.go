package sshocks

import (
	"encoding/json"

	"github.com/audibleblink/HoleySocks/pkg/holeysocks"
	"github.com/gobuffalo/packr"
)

var Config holeysocks.MainConfig

func init() {
	box := packr.NewBox("../../configs")
	configBytes, err := box.Find("ssh.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configBytes, &Config)
	if err != nil {
		panic(err)
	}

	privKeyBytes, err := box.Find("id_ed25519")
	if err != nil {
		panic(err)
	}
	Config.SSH.SetKey(privKeyBytes)
}
