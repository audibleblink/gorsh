package sitrep

import (
	"fmt"
	"net"
	"os"
	"os/user"

	"github.com/fatih/structs"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

type Printable interface {
	String() string
}

type Process struct {
	PID     int32
	PPID    int32
	Name    string
	Owner   string
	exe     string
	Cmdline string
}

func (p *Process) String() string {
	template := "%8v | %-4v\n"
	out := _printer(template, p)
	return out
}

type Host struct {
	Hostname string
	Procs    uint64
	OS       string
	Platform string
	Family   string
	Version  string
	Kernel   string
}

func (h *Host) String() string {
	template := "%8v | %-4v\n"
	out := _printer(template, h)
	return out
}

type Iface struct {
	Name      string
	Addresses string
}

func (i *Iface) String() string {
	template := "%10v | %-4v\n"
	out := _printer(template, i)
	return out
}

type User struct {
	Username string
	Uid      string
	Gid      string
	Homedir  string
}

func (u *User) String() string {
	template := "%8v | %-4v\n"
	out := _printer(template, u)
	return out
}

func All() string {
	var output string
	output += Header("Host")
	output += Gethost()
	output += Header("User")
	output += Getuser()
	output += Header("Network")
	output += Getnetwork()
	dir, _ := os.Getwd()
	output += fmt.Sprintf("Current Directory: %s", dir)
	return output
}

func Header(name string) string {
	var output string
	output += "\n========================================\n"
	output += name + "\n"
	output += "========================================\n"
	return output
}

func Gethost() string {
	info, err := host.Info()
	if err != nil {
		return err.Error()
	}
	// "%8s | %4s\n",
	host := Host{
		info.Hostname,
		info.Procs,
		info.OS,
		info.Platform,
		info.PlatformFamily,
		info.PlatformVersion,
		info.KernelVersion}

	return host.String()
}

func Processes() string {
	var output string

	procs, err := process.Processes()
	if err != nil {
		fmt.Printf(err.Error())
	}
	for _, p := range procs {
		pid := p.Pid
		ppid, _ := p.Ppid()
		name, _ := p.Name()
		user, _ := p.Username()
		exe, _ := p.Exe()
		cmd, _ := p.Cmdline()

		proc := Process{pid, ppid, name, user, exe, cmd}
		output += proc.String()
	}
	return output
}

func Getnetwork() string {
	var output string

	networks, _ := net.Interfaces()
	for _, net := range networks {
		addrs, _ := net.Addrs()
		for _, addr := range addrs {
			network := Iface{net.Name, addr.String()}
			output += network.String()
		}
	}
	return output
}

func Getuser() string {
	userData, _ := user.Current()
	user := User{userData.Username, userData.Uid, userData.Gid, userData.HomeDir}
	return user.String()
}

func _printer(template string, p Printable) string {
	sMap := structs.Map(p)
	var output string
	for k, v := range sMap {
		output += fmt.Sprintf(template, k, v)
	}
	output += "\n"
	return output
}
