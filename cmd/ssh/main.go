package main

import (
	"io/ioutil"
	"os"
	"strconv"

	sshlib "github.com/blacknon/go-sshlib"
	"github.com/jessevdk/go-flags"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/helpers"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/types"
	"golang.org/x/crypto/ssh"
)

type Options struct {
	SSHHost    string `short:"H" long:"ssh-host" description:"ssh's host to login" optional:"no"`
	SSHPort    int    `short:"p" long:"ssh-port" description:"ssh's port to use in login" optional:"no"`
	SSHUser    string `short:"u" long:"ssh-user" description:"ssh's user to use in login" optional:"no"`
	BindPort   int    `short:"b" long:"bind-port" description:"local port to bind the proxy" optional:"no"`
	SSHKeyPath string `short:"k" long:"ssh-key-path" description:"ssh's key to use in login" optional:"no"`
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, _ := ioutil.ReadFile(keyPath)
	return ssh.ParsePrivateKey(buff)
}

func createSshTunnel(configuration SSHProxyConfig) {
	con := &sshlib.Connect{}
	key, err := parsePrivateKey(configuration.SSHProxy.SSH.KeyPath)
	CheckGenericError(err)
	key_content := []ssh.AuthMethod{ssh.PublicKeys(key)}
	err = con.CreateClient(
		configuration.SSHProxy.SSH.Host,
		strconv.Itoa(configuration.SSHProxy.SSH.Port),
		configuration.SSHProxy.SSH.User,
		key_content,
	)
	CheckGenericError(err)
	con.TCPDynamicForward("localhost", strconv.Itoa(configuration.SSHProxy.BindPort))
}

func main() {
	var conf SSHProxyConfig
	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	conf.SSHProxy.SSH.Host = options.SSHHost
	conf.SSHProxy.SSH.Port = options.SSHPort
	conf.SSHProxy.SSH.User = options.SSHUser
	conf.SSHProxy.SSH.KeyPath = options.SSHKeyPath
	conf.SSHProxy.BindPort = options.BindPort
	createSshTunnel(conf)
}
