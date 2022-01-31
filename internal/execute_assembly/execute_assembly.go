package execute_assembly

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io/ioutil"

	clr "github.com/Ne0nd0g/go-clr"
)

//go:embed assemblies/*
var Assemblies embed.FS

type CLR struct {
	Active string

	loaded      bool
	runtimeHost *clr.ICORRuntimeHost
	assembly    *clr.MethodInfo
}

func NewCLR() (CLR, error) {
	c := &CLR{
		loaded: false,
	}
	return c.LoadCLR()
}

func (e CLR) LoadCLR() (CLR, error) {
	err := clr.RedirectStdoutStderr()
	if err != nil {
		return e, err
	}

	e.runtimeHost, err = clr.LoadCLR("v4")
	if err != nil {
		return e, err
	}

	e.loaded = true
	return e, err
}

func (e CLR) IsLoaded() bool {
	return e.loaded
}

func (e *CLR) Load(name string) (err error) {
	assemblyBytes, err := unzippedBytes(name)
	if err != nil {
		return
	}

	e.assembly, err = clr.LoadAssembly(e.runtimeHost, assemblyBytes)
	if err != nil {
		return
	}

	e.Active = name
	return
}

func (e CLR) Execute(args []string) (stdout, stderr string) {
	stdout, stderr = clr.InvokeAssembly(e.assembly, args)
	return
}

func unzippedBytes(name string) (asmBytes []byte, err error) {
	file := fmt.Sprintf("assemblies/%s.exe.gz", name)
	gZipFile, err := Assemblies.Open(file)
	if err != nil {
		return
	}
	defer gZipFile.Close()

	gzReader, err := gzip.NewReader(gZipFile)
	if err != nil {
		return
	}
	defer gzReader.Close()

	asmBytes, err = ioutil.ReadAll(gzReader)
	if err != nil {
		return
	}

	if len(asmBytes) == 0 {
		err = fmt.Errorf("null read size for %s", name)
	}
	return
}
