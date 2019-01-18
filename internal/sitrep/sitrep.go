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

type printable interface {
	String() string
}

type Process struct {
	PID     int32
	PPID    int32
	Name    string
	Owner   string
	Exe     string
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

type Interface struct {
	Name      string
	Addresses []string
}

func (i *Interface) String() string {
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

func SysInfo() string {
	var output string
	output += header("Host")
	if host, err := HostInfo(); err == nil {
		output += host.String()
	}
	output += header("User")
	if user, err := UserInfo(); err == nil {
		output += user.String()
	}
	output += header("Network")
	if networks, err := NetworkInfo(); err == nil {
		var netString string
		for _, net := range networks {
			netString += net.String()
		}
		output += netString
	}
	dir, _ := os.Getwd()
	output += fmt.Sprintf("Current Directory: %s", dir)
	return output
}

func Processes() string {
	var output string
	processes, err := ProcessInfo()
	if err != nil {
		return err.Error()
	}
	for _, process := range processes {
		output += process.String()
	}
	return output
}

func HostInfo() (Host, error) {
	info, err := host.Info()
	if err != nil {
		return Host{}, err
	}

	host := Host{
		info.Hostname,
		info.Procs,
		info.OS,
		info.Platform,
		info.PlatformFamily,
		info.PlatformVersion,
		info.KernelVersion}

	return host, err
}

// ProcessInfo returns the list of interface for the current host
func ProcessInfo() ([]Process, error) {
	var results []Process
	procs, err := process.Processes()
	if err != nil {
		return []Process{}, err
	}
	for _, p := range procs {
		pid := p.Pid
		ppid, _ := p.Ppid()
		name, _ := p.Name()
		user, _ := p.Username()
		exe, _ := p.Exe()
		cmd, _ := p.Cmdline()

		proc := Process{pid, ppid, name, user, exe, cmd}
		results = append(results, proc)
	}
	return results, err
}

// NetworkInfo returns the list of interface for the current host
func NetworkInfo() ([]Interface, error) {
	var (
		err     error
		results []Interface
	)

	networks, err := net.Interfaces()
	for _, net := range networks {
		addrs, err := net.Addrs()
		if err != nil {
			continue
		}
		var ips []string
		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}
		network := Interface{net.Name, ips}
		results = append(results, network)
	}
	return results, err
}

// UserInfo returns data about the current user like their UID/GID and home directory
func UserInfo() (User, error) {
	userData, err := user.Current()
	if err != nil {
		return User{}, err
	}
	user := User{userData.Username, userData.Uid, userData.Gid, userData.HomeDir}
	return user, err
}

func _printer(template string, p printable) string {
	sMap := structs.Map(p)
	var output string
	for k, v := range sMap {
		output += fmt.Sprintf(template, k, v)
	}
	output += "\n"
	return output
}

func header(name string) string {
	var output string
	output += "\n========================================\n"
	output += name + "\n"
	output += "========================================\n"
	return output
}

func Environ() string {
	data := os.Environ()
	var output string
	for _, el := range data {
		output += el + "\n"
	}
	return output
}
