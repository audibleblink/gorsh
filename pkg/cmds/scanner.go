package cmds

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/scanner"
)

func Scanner(c *ishell.Context) {
	if len(c.Args) != 1 {
		c.Println(c.Cmd.Help)
		return
	}

	var (
		net, lvl string
	)

	net = c.Args[0]
	if net == "localhost" {
		net = "127.0.0.1/32"
	}

	if !strings.Contains(net, "/") {
		net = net + "/32"
	}

	if !validCIDR(net) {
		c.Printf("invalid cidr: %s\n", net)
		return
	}

	if len(c.Args) > 1 {
		lvl = c.Args[1]
	}

	var portRange []int
	switch lvl {
	case "500":
		portRange = scanner.TOP_500
	default:
		portRange = scanner.TOP_250
	}

	now := time.Now()

	c.Printf("Ping sweeping %s...", net)
	hosts, err := scanner.Sweep(net)
	if err != nil {
		c.Println("not admin. skipping ping sweep")
	} else {
		c.Println("")
		for _, h := range hosts {
			c.Printf("Host %s : %s : %s alive\n", h.Hostname, h.MAC, h.IP.String())
		}
	}

	allResults := make(chan string, 20)
	go func(cn chan string, c *ishell.Context) {
		for res := range cn {
			c.Println(res)
		}
	}(allResults, c)

	for _, host := range hosts {
		tcpScanner := scanner.NewTCPScanner(host.IP.String())
		tcpScanner.Start(portRange, allResults)
	}

	log.Print("time spent: " + time.Now().Sub(now).String())

}

func validCIDR(ip string) bool {
	_, _, err := net.ParseCIDR(ip)
	return err == nil
}
