package main

import (
	"flag"
	"fmt"

	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/helpers"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/pkg/sshproxy"
)

func main() {
	var err error
	start := flag.Bool("start", false, "start ssh proxy")
	stop := flag.Bool("stop", false, "stop ssh proxy")
	status := flag.Bool("status", true, "status of the ssh proxy")
	flag.Parse()
	SSHProxy := SSHProxy{}
	err = SSHProxy.Init()
	if err != nil {
		CheckGenericError(err)
	}
	if *start {
		err = SSHProxy.Start()
		if err != nil {
			CheckGenericError(err)
		}
	} else if *stop {
		err = SSHProxy.Stop()
		if err != nil {
			CheckGenericError(err)
		}
	} else if *status {
		message := SSHProxy.Status()
		fmt.Println(message)
	} else {
		message := SSHProxy.Status()
		fmt.Println(message)
	}
}
