package sitrep

import (
	"fmt"
	"net"
	"os"
	"os/user"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

func All() string {
	var output string
	output += Header("Host")
	output += HostInfo()
	output += Header("User")
	output += User()
	output += Header("Network")
	output += Network()
	return output
}

func Header(name string) string {
	var output string
	output += "\n========================================\n"
	output += name + "\n"
	output += "========================================\n"
	return output
}

func HostInfo() string {
	info, err := host.Info()
	if err != nil {
		return err.Error()
	}
	infoBlock := fmt.Sprintf("%8s | %4s\n"+
		"%8s | %4d\n"+
		"%8s | %4s\n"+
		"%8s | %4s\n"+
		"%8s | %4s\n"+
		"%8s | %4s\n"+
		"%8s | %4s\n",
		"Hostname", info.Hostname,
		"Procs", info.Procs,
		"OS", info.OS,
		"Platform", info.Platform,
		"Family", info.PlatformFamily,
		"Version", info.PlatformVersion,
		"Kernel", info.KernelVersion)
	return infoBlock + "\n"
}

func Processes() string {
	var output string

	procs, err := process.Processes()
	if err != nil {
		fmt.Printf(err.Error())
	}
	for _, proc := range procs {
		pid := proc.Pid
		ppid, _ := proc.Ppid()
		name, _ := proc.Name()
		exe, _ := proc.Exe()
		cmd, _ := proc.Cmdline()
		user, _ := proc.Username()
		procBlock := fmt.Sprintf("%8s | %4d\n"+
			"%8s | %4d\n"+
			"%8s | %4s\n"+
			"%8s | %4s\n"+
			"%8s | %4s\n"+
			"%8s | %4s\n",
			"PID", pid,
			"PPID", ppid,
			"Name", name,
			"Owner", user,
			"exe", exe,
			"Cmdline", cmd)
		output += (procBlock + "\n")
	}
	return output

}

func Network() string {
	var output string

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		output += fmt.Sprintf("%8s | ", iface.Name)
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			output += fmt.Sprintf("%-18s\t", addr.String())
		}
		output += "\n"
	}
	return output
}

func User() string {
	var output string
	user, _ := user.Current()
	userBlock := fmt.Sprintf("%8s | %4s\n"+
		"%8s | %4s\n",
		"User", user.Username, "ID", user.Uid)
	output += userBlock

	dir, _ := os.Getwd()
	dirBlock := fmt.Sprintf("%8s | %4s\n"+
		"%8s | %4s\n",
		"CurrDir", dir, "Home", user.HomeDir)

	output += dirBlock
	return output
}
