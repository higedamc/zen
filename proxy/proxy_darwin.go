package proxy

import (
	"fmt"
	"os/exec"
	"strings"
)

var exclusionListURLs = []string{
	"https://raw.githubusercontent.com/anfragment/zen/main/proxy/exclusions/common.txt",
	"https://raw.githubusercontent.com/anfragment/zen/main/proxy/exclusions/apple.txt",
}

var interfaceName string

// setSystemProxy sets the system proxy to the proxy address
func (p *Proxy) setSystemProxy() error {
	cmd := exec.Command("sh", "-c", "networksetup -listnetworkserviceorder | grep `route -n get 0.0.0.0 | grep 'interface' | cut -d ':' -f2` -B 1 | head -n 1 | cut -d ' ' -f2")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("get default interface: %v\n%s", err, out)
	}

	interfaceName = strings.TrimSpace(string(out))
	cmd = exec.Command("networksetup", "-setwebproxy", interfaceName, "127.0.0.1", fmt.Sprint(p.port))
	if out, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("set web proxy (interface: %s, port: %d): %v\n%s", interfaceName, p.port, err, out)
	}

	cmd = exec.Command("networksetup", "-setsecurewebproxy", interfaceName, "127.0.0.1", fmt.Sprint(p.port))
	if out, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("set secure web proxy (interface: %s, port: %d): %v\n%s", interfaceName, p.port, err, out)
	}

	return nil
}

func (p *Proxy) unsetSystemProxy() error {
	if interfaceName == "" {
		return fmt.Errorf("trying to unset system proxy without setting it first")
	}

	cmd := exec.Command("networksetup", "-setwebproxystate", interfaceName, "off")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unset web proxy (interface: %s): %v\n%s", interfaceName, err, out)
	}

	cmd = exec.Command("networksetup", "-setsecurewebproxystate", interfaceName, "off")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("unset secure web proxy (interface: %s): %v\n%s", interfaceName, err, out)
	}

	return nil
}
