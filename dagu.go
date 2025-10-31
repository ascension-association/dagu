// Binary dagu is a gokrazy wrapper program that runs the bundled dagu
// executable in /usr/local/bin/dagu after doing any necessary runtime system
// setup.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"net"
)

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
        return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
    }
    return ipv4Addr.String(), nil
}

func isMounted(mountpoint string) (bool, error) {
	b, err := ioutil.ReadFile("/proc/self/mountinfo")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // platform does not have /proc/self/mountinfo, fall back to not verifying
		}
		return false, err
	}

	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		if parts[4] == mountpoint {
			return true, nil
		}
	}

	return false, nil
}

func makeWritable(dir string) error {
	mounted, err := isMounted(dir)
	if err != nil {
		return err
	}
	if mounted {
		// Nothing to do, directory is already mounted.
		return nil
	}

	// Read all regular files in this directory.
	regularFiles := make(map[string]string)
	fis, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		b, err := os.ReadFile(filepath.Join(dir, fi.Name()))
		if err != nil {
			return err
		}
		regularFiles[fi.Name()] = string(b)
	}

	if err := syscall.Mount("tmpfs", dir, "tmpfs", 0, ""); err != nil {
		return fmt.Errorf("tmpfs on %s: %v", dir, err)
	}

	// Write all regular files from memory back to new tmpfs.
	for name, contents := range regularFiles {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0644); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ipAddress, err := GetInterfaceIpv4Addr("eth0")
	if err != nil {
		return err
	}
	if err := syscall.Exec("/usr/local/bin/dagu", []string{"start-all"}, expandPath(append(os.Environ(), "DAGU_HOST=" + ipAddress))); err != nil {
		log.Fatal(err)
	}
}

// expandPath returns env, but with PATH= modified or added
// such that both /user and /usr/local/bin are included, which containerd needs.
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