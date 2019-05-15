//go:generate go run embed.go

package main

import (
	"github.com/shurcooL/vfsgen"
	"net/http"
)

func main() {
	scriptDir := http.Dir("scripts")

	err := vfsgen.Generate(scriptDir, vfsgen.Options{
		Filename:     "internal/enum/data.go",
		PackageName:  "enum",
		VariableName: "scripts",
	})
	if err != nil {
		panic(err)
	}

	keyDir := http.Dir("configs")
	err = vfsgen.Generate(keyDir, vfsgen.Options{
		Filename:     "internal/sshocks/data.go",
		PackageName:  "sshocks",
		VariableName: "configs",
	})
	if err != nil {
		panic(err)
	}
}
