package main

import (
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/jessevdk/go-flags"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/helpers"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/types"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type options struct {
	SSHHost    string `short:"H" long:"ssh-host" description:"ssh's host to login" optional:"no"`
	SSHPort    int    `short:"p" long:"ssh-port" description:"ssh's port to use in login" optional:"no"`
	SSHUser    string `short:"u" long:"ssh-user" description:"ssh's user to use in login" optional:"no"`
	BindPort   int    `short:"b" long:"bind-port" description:"local port to bind the proxy" optional:"no"`
	SSHKeyPath string `short:"k" long:"ssh-key-path" description:"ssh's key to use in login" optional:"no"`
}

func getAuthType(keyPath string) []ssh.AuthMethod {
	var auth []ssh.AuthMethod
	if keyPath == "" {
		// Will probably not work with windows
		socket := os.Getenv("SSH_AUTH_SOCK")
		conn, err := net.Dial("unix", socket)
		CheckGenericError(err)
		agentClient := agent.NewClient(conn)
		auth = []ssh.AuthMethod{ssh.PublicKeysCallback(agentClient.Signers)}
	} else {
		buff, _ := ioutil.ReadFile(keyPath)
		key, err := ssh.ParsePrivateKey(buff)
		CheckGenericError(err)
		auth = []ssh.AuthMethod{ssh.PublicKeys(key)}
	}
	return auth
}

func createSSHTunnel(configuration KubeSSHProxyConfig) {
	var auth []ssh.AuthMethod
	auth = getAuthType(configuration.KubeSSHProxy.SSH.KeyPath)
	con := &connect{}
	err := con.CreateClient(
		configuration.KubeSSHProxy.SSH.Host,
		strconv.Itoa(configuration.KubeSSHProxy.SSH.Port),
		configuration.KubeSSHProxy.SSH.User,
		auth,
	)
	CheckGenericError(err)
	err = con.TCPDynamicForward("localhost", strconv.Itoa(configuration.KubeSSHProxy.BindPort))
	CheckGenericError(err)
}

func main() {
	var conf KubeSSHProxyConfig
	var opt options
	var parser = flags.NewParser(&opt, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	conf.KubeSSHProxy.SSH.Host = opt.SSHHost
	conf.KubeSSHProxy.SSH.Port = opt.SSHPort
	conf.KubeSSHProxy.SSH.User = opt.SSHUser
	conf.KubeSSHProxy.SSH.KeyPath = opt.SSHKeyPath
	conf.KubeSSHProxy.BindPort = opt.BindPort
	createSSHTunnel(conf)
}
