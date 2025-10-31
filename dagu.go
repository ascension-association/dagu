// Binary dagu is a gokrazy wrapper program that runs the bundled dagu
// executable in /usr/local/bin/dagu after doing any necessary runtime system
// setup.
package main

import (
	"fmt"
	"log"
	"context"
	"errors"
	"net"

	execute "github.com/alexellis/go-execute/v2"
	"github.com/gokrazy/gokrazy"
)

var port = "8080"

// https://gist.github.com/schwarzeni/f25031a3123f895ff3785970921e962c
func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
    var (
        ief      *net.Interface
        addrs    []net.Addr
        ipv4Addr net.IP
    )
    if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
        return
    }
    if addrs, err = ief.Addrs(); err != nil { // get addresses
        return
    }
    for _, addr := range addrs { // get ipv4 address
        if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
            break
        }
    }
    if ipv4Addr == nil {
        return "", errors.New(fmt.Sprintf("interface %s does not have an ipv4 address\n", interfaceName))
    }
    return ipv4Addr.String(), nil
}

func run(logging bool, exe string, args ...string) {
	var cmd execute.ExecTask

	if logging {
		cmd = execute.ExecTask{
			Command:     exe,
			Args:        args,
			StreamStdio: true,
		}
	} else {
		cmd = execute.ExecTask{
			Command:     exe,
			Args:        args,
			StreamStdio: false,
			DisableStdioBuffer: true,
		}
	}

	res, err := cmd.Execute(context.Background())

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	if res.ExitCode != 0 {
		fmt.Errorf("Error: %v", res.Stderr)
	}
}

func main() {
	// wait for local network
	gokrazy.WaitFor("net-route")

	// get local IP address
	ipAddress, err := GetInterfaceIpv4Addr("eth0")
	if err != nil {
		ipAddress = "127.0.0.1"
	}
	log.Println("Local IP Address: " + ipAddress)

	// create mount point and use for Dagu storage
	run(false, "/usr/local/bin/busybox", "mkdir", "-p", "/perm/dagu")
	run(false, "export", "DAGU_HOME=/perm/dagu")

	// run Dagu
	run(false, "export", "PATH=/user:/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin")
	run(true, "/usr/local/bin/dagu", "server", "--host", ipAddress, "--port", port)
	run(true, "/usr/local/bin/dagu", "scheduler")
	run(true, "/usr/local/bin/dagu", "coordinator")
}