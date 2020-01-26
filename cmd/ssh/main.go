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

type options struct {
	SSHHost    string `short:"H" long:"ssh-host" description:"ssh's host to login" optional:"no"`
	SSHPort    int    `short:"p" long:"ssh-port" description:"ssh's port to use in login" optional:"no"`
	SSHUser    string `short:"u" long:"ssh-user" description:"ssh's user to use in login" optional:"no"`
	BindPort   int    `short:"b" long:"bind-port" description:"local port to bind the proxy" optional:"no"`
	SSHKeyPath string `short:"k" long:"ssh-key-path" description:"ssh's key to use in login" optional:"no"`
}

func createSSHTunnel(configuration SSHProxyConfig) {
	con := &sshlib.Connect{}
	buff, _ := ioutil.ReadFile(configuration.SSHProxy.SSH.KeyPath)
	key, err := ssh.ParsePrivateKey(buff)
	CheckGenericError(err)
	keyContent := []ssh.AuthMethod{ssh.PublicKeys(key)}
	err = con.CreateClient(
		configuration.SSHProxy.SSH.Host,
		strconv.Itoa(configuration.SSHProxy.SSH.Port),
		configuration.SSHProxy.SSH.User,
		keyContent,
	)
	CheckGenericError(err)
	err = con.TCPDynamicForward("localhost", strconv.Itoa(configuration.SSHProxy.BindPort))
	CheckGenericError(err)
}

func main() {
	var conf SSHProxyConfig
	var opt options
	var parser = flags.NewParser(&opt, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	conf.SSHProxy.SSH.Host = opt.SSHHost
	conf.SSHProxy.SSH.Port = opt.SSHPort
	conf.SSHProxy.SSH.User = opt.SSHUser
	conf.SSHProxy.SSH.KeyPath = opt.SSHKeyPath
	conf.SSHProxy.BindPort = opt.BindPort
	createSSHTunnel(conf)
}
