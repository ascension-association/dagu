// Binary dagu is a gokrazy wrapper program that runs the bundled dagu
// executable in /usr/local/bin/dagu after doing any necessary runtime system
// setup.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"errors"
	"net"
	"io/ioutil"

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

func run(logging bool, command string) {
	cmd := exec.Command("/usr/local/bin/busybox", "sh", "-c", command)

	if logging {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
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
	run(false, "/usr/local/bin/busybox mkdir -p /perm/dagu/.config/dagu")

	// enable basic auth
	config := "/perm/dagu/.config/dagu/config.yaml"
	if _, err = os.Stat(config); os.IsNotExist(err) {
		pw, _ := ioutil.ReadFile("/etc/gokr-pw.txt")
		f, _ := os.OpenFile(config, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		defer f.Close()
		_, err = f.WriteString("auth:\n  basic:\n    username: gokrazy\n    password: ")
		if err != nil {
			fmt.Errorf("Error: %v", err)
		} else {
			_, err = f.WriteString(string(pw))
			if err != nil {
				fmt.Errorf("Error: %v", err)
			}
		}
	}

	// run Dagu
	command := "export PATH=/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin ; export DAGU_HOME=/perm/dagu ; /usr/local/bin/busybox nohup /usr/local/bin/dagu server --dagu-home /perm/dagu --config /perm/dagu/.config/dagu/config.yaml --host " + ipAddress + " --port " + port + " ; /usr/local/bin/busybox nohup /usr/local/bin/dagu scheduler --dagu-home /perm/dagu ; /usr/local/bin/busybox nohup /usr/local/bin/dagu coordinator --dagu-home /perm/dagu"
	run(true, command)
}