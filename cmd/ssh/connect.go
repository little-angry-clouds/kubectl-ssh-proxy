package main

import (
	"context"
	"github.com/armon/go-socks5"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"net"
	"time"
)

// All the credit goes to https://github.com/blacknon/go-sshlib/
// This file has a different license, the one you can see in the projects github

type connect struct {
	Client *ssh.Client
}

func (c *connect) CreateClient(host, port, user string, authMethods []ssh.AuthMethod) (err error) {
	uri := net.JoinHostPort(host, port)

	timeout := 20

	// Create new ssh.ClientConfig{}
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(timeout) * time.Second,
	}

	proxyDialer := proxy.Direct

	// Dial to host:port
	netConn, err := proxyDialer.Dial("tcp", uri)
	if err != nil {
		return
	}

	// Create new ssh connect
	sshCon, channel, req, err := ssh.NewClientConn(netConn, uri, config)
	if err != nil {
		return
	}

	// Create *ssh.Client
	c.Client = ssh.NewClient(sshCon, channel, req)

	return
}

// socks5Resolver prevents DNS from resolving on the local machine, rather than over the SSH connection.
type socks5Resolver struct{}

func (socks5Resolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	return ctx, nil, nil
}

func (c *connect) TCPDynamicForward(address, port string) (err error) {
	// Create Socks5 config
	conf := &socks5.Config{
		Dial: func(ctx context.Context, n, addr string) (net.Conn, error) {
			return c.Client.Dial(n, addr)
		},
		Resolver: socks5Resolver{},
	}

	// Create Socks5 server
	s, err := socks5.New(conf)
	if err != nil {
		return
	}

	// Listen
	err = s.ListenAndServe("tcp", net.JoinHostPort(address, port))

	return
}
