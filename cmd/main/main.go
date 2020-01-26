package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/helpers"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/types"
	"gopkg.in/yaml.v2"
)

func createCommandArgs(kubeconfig Kubeconfig) []string {
	args := []string{
		"-H", kubeconfig.SSHProxy.SSH.Host,
		"-p", strconv.Itoa(kubeconfig.SSHProxy.SSH.Port),
		"-u", kubeconfig.SSHProxy.SSH.User,
		"-b", strconv.Itoa(kubeconfig.SSHProxy.BindPort),
	}
	if kubeconfig.SSHProxy.SSH.KeyPath != "" {
		args = append(args, "-k", kubeconfig.SSHProxy.SSH.KeyPath)
	}
	return args
}

func main() {
	var err error
	var conf Kubeconfig

	yamlFile, err := ioutil.ReadFile(os.Getenv("KUBECONFIG"))
	err = yaml.Unmarshal(yamlFile, &conf)
	CheckGenericError(err)
	// TODO mac y windows?
	pidDir := fmt.Sprintf("%s/kubectl-ssh-proxy/%s", os.Getenv("XDG_RUNTIME_DIR"), conf.CurrentCluster)
	pidPath := fmt.Sprintf("%s/PID", pidDir)
	start := flag.Bool("start", false, "start ssh proxy")
	stop := flag.Bool("stop", false, "stop ssh proxy")
	status := flag.Bool("status", true, "status of the ssh proxy")
	flag.Parse()

	if *start {
		if _, err := os.Stat(pidPath); err == nil {
			fmt.Println("There's already an active process!")
			os.Exit(1)
		}
		if _, err := os.Stat(pidDir); os.IsNotExist(err) {
			err = os.MkdirAll(pidDir, 0755)
			CheckGenericError(err)
		}
		args := createCommandArgs(conf)
		cmd := exec.Command("./bin/kubectl-ssh-proxy-ssh-bin", args...)
		err = cmd.Start()
		CheckGenericError(err)
		pid := []byte(strconv.Itoa(cmd.Process.Pid))
		err = ioutil.WriteFile(pidPath, pid, 0644)
		CheckGenericError(err)
		fmt.Println("Proxy started!")
		fmt.Println("Eval the next: \neval HTTPS_PROXY=socks5://localhost:8080")

	} else if *stop {
		if _, err := os.Stat(pidPath); err != nil {
			fmt.Println("The ssh proxy is already stopped!")
			os.Exit(1)
		}
		file, err := os.Open(pidPath)
		CheckGenericError(err)
		defer file.Close()
		pid, err := ioutil.ReadAll(file)
		CheckGenericError(err)
		p, _ := strconv.Atoi(string(pid))
		process, _ := os.FindProcess(p)
		process.Signal(os.Interrupt)
		fmt.Println("Proxy stopped!")
		os.Remove(pidPath)
	} else if *status {
		if _, err := os.Stat(pidPath); err == nil {
			fmt.Println("Proxy activated!")
		} else {
			fmt.Println("Proxy stopped!")
		}
	} else {
		if _, err := os.Stat(pidPath); err == nil {
			fmt.Println("Proxy activated!")
		} else {
			fmt.Println("Proxy stopped!")
		}
	}
}
