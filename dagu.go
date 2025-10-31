// Binary dagu is a gokrazy wrapper program that runs the bundled dagu
// executable in /usr/local/bin/dagu after doing any necessary runtime system
// setup.
package main

import (
	"fmt"
	"log"
	"errors"
	"net"
	"os"
	"strings"
	"syscall"

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
	cmd := []string{"/usr/local/bin/busybox", "mkdir", "-p", "/perm/dagu"}
	err := syscall.Exec(cmd[0], cmd, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// run Dagu
	cmd = []string{"/usr/local/bin/dagu", "server", "--host", ipAddress, "--port", port}
	err = syscall.Exec(cmd[0], cmd, expandPath(append(os.Environ(), "DAGU_HOME=/perm/dagu")))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	cmd = []string{"/usr/local/bin/dagu", "scheduler"}
	err = syscall.Exec(cmd[0], cmd, expandPath(append(os.Environ(), "DAGU_HOME=/perm/dagu")))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	cmd = []string{"/usr/local/bin/dagu", "coordinator"}
	err = syscall.Exec(cmd[0], cmd, expandPath(append(os.Environ(), "DAGU_HOME=/perm/dagu")))
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// expandPath returns env, but with PATH= modified or added such that both /user and /usr/local/bin are included.
func expandPath(env []string) []string {
	extra := "/user:/usr/local/bin"
	found := false
	for idx, val := range env {
		parts := strings.Split(val, "=")
		if len(parts) < 2 {
			continue // malformed entry
		}
		key := parts[0]
		if key != "PATH" {
			continue
		}
		val := strings.Join(parts[1:], "=")
		env[idx] = fmt.Sprintf("%s=%s:%s", key, extra, val)
		found = true
	}
	if !found {
		const busyboxDefaultPATH = "/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin"
		env = append(env, fmt.Sprintf("PATH=%s:%s", extra, busyboxDefaultPATH))
	}
	return env
}
