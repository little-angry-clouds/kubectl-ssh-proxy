package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/helpers"
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/types"
	"gopkg.in/yaml.v2"
)

// SSHProxy is the main object
type SSHProxy struct {
	kubeconfig Kubeconfig
	pidPath    string
}

// Init initializes the SSHProxy object
func (proxy *SSHProxy) Init() {
	proxy.getKubeconfig()
	proxy.getPidPath()
}

// Start starts the SSHProxy
func (proxy *SSHProxy) Start() {
	var err error
	pidPath := proxy.pidPath
	pidDir := path.Dir(pidPath)
	CheckActiveProcess(pidPath)
	if _, err := os.Stat(pidDir); os.IsNotExist(err) {
		err = os.MkdirAll(pidDir, 0755)
		CheckGenericError(err)
	}
	args := proxy.createArgs()
	cmd := exec.Command("kubectl-ssh-proxy-ssh-bin", args...)
	err = cmd.Start()
	CheckGenericError(err)
	// Capture the state of the subcommand
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:
		message := fmt.Sprintf("The ssh proxy failed. The error is: %s", err)
		fmt.Println(message)
		message = fmt.Sprintf("You may debug the error executing the ssh binary manually: %s", cmd.String())
		fmt.Println(message)
		os.Exit(1)
	case <-time.After(10 * time.Millisecond):
		// timeout, which means that all works correctly
	}
	pid := []byte(strconv.Itoa(cmd.Process.Pid))
	err = ioutil.WriteFile(pidPath, pid, 0644)
	CheckGenericError(err)
	fmt.Println("Proxy started!")
	fmt.Println("Eval the next: \nexport HTTPS_PROXY=socks5://localhost:8080")
}

// Stop stops the SSHProxy
func (proxy *SSHProxy) Stop() {
	pidPath := proxy.pidPath
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
	fmt.Println("SSH proxy stopped!")
	os.Remove(pidPath)
}

// Status gets the SSHProxy status
func (proxy *SSHProxy) Status() {
	pidPath := proxy.pidPath
	if _, err := os.Stat(pidPath); err == nil {
		fmt.Println("SSH proxy activated!")
	} else {
		fmt.Println("SSH proxy stopped!")
	}
}

func (proxy *SSHProxy) getPidPath() {
	// TODO probar con mac y winsux
	pidDir := fmt.Sprintf("%s/kubectl-ssh-proxy/%s", os.Getenv("XDG_RUNTIME_DIR"), proxy.kubeconfig.CurrentCluster)
	pidPath := fmt.Sprintf("%s/PID", pidDir)
	proxy.pidPath = pidPath
}

func (proxy *SSHProxy) createArgs() []string {
	args := []string{
		"-H", proxy.kubeconfig.SSHProxy.SSH.Host,
		"-p", strconv.Itoa(proxy.kubeconfig.SSHProxy.SSH.Port),
		"-u", proxy.kubeconfig.SSHProxy.SSH.User,
		"-b", strconv.Itoa(proxy.kubeconfig.SSHProxy.BindPort),
	}
	if proxy.kubeconfig.SSHProxy.SSH.KeyPath != "" {
		args = append(args, "-k", proxy.kubeconfig.SSHProxy.SSH.KeyPath)
	}
	return args
}

func (proxy *SSHProxy) getKubeconfig() {
	var kubeconfig Kubeconfig
	yamlFile, err := ioutil.ReadFile(os.Getenv("KUBECONFIG"))
	CheckGenericError(err)
	err = yaml.Unmarshal(yamlFile, &kubeconfig)
	CheckGenericError(err)
	proxy.kubeconfig = kubeconfig
}

func main() {
	start := flag.Bool("start", false, "start ssh proxy")
	stop := flag.Bool("stop", false, "stop ssh proxy")
	status := flag.Bool("status", true, "status of the ssh proxy")
	flag.Parse()
	SSHProxy := SSHProxy{}
	SSHProxy.Init()
	if *start {
		SSHProxy.Start()
	} else if *stop {
		SSHProxy.Stop()
	} else if *status {
		SSHProxy.Status()
	} else {
		SSHProxy.Status()
	}
}
