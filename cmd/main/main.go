package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
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
	cmd := exec.Command("kube-ssh-proxy-ssh-bin", args...)
	err = cmd.Start()
	CheckGenericError(err)
	// Capture the state of the subcommand. To do it it's necessary to add a little sleep
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:
		message := fmt.Sprintf("# The ssh proxy failed. The error is: %s", err)
		fmt.Println(message)
		message = fmt.Sprintf("# You may debug the error by executing the ssh binary manually: %s %s", "kube-ssh-proxy-ssh-bin", strings.Join(args[:], " "))
		fmt.Println(message)
		os.Exit(1)
	case <-time.After(1000 * time.Millisecond):
		// timeout, which means that all works correctly
	}
	pid := []byte(strconv.Itoa(cmd.Process.Pid))
	err = ioutil.WriteFile(pidPath, pid, 0644)
	CheckGenericError(err)
	fmt.Println("# The SSH Proxy started!")
	fmt.Println("# Eval the next: \nexport HTTPS_PROXY=socks5://localhost:8080")
}

// Stop stops the SSHProxy
func (proxy *SSHProxy) Stop() {
	pidPath := proxy.pidPath
	if _, err := os.Stat(pidPath); err != nil {
		fmt.Println("# The ssh proxy is already stopped!")
		os.Exit(0)
	}
	file, err := os.Open(pidPath)
	CheckGenericError(err)
	defer file.Close()
	pid, err := ioutil.ReadAll(file)
	CheckGenericError(err)
	p, _ := strconv.Atoi(string(pid))
	process, _ := os.FindProcess(p)
	process.Signal(os.Interrupt)
	// TODO probar con mac y winsux
	fmt.Println("# The SSH Proxy is already stopped! Eval the next:\nunset HTTPS_PROXY")
	os.Remove(pidPath)
}

// Status gets the SSHProxy status
func (proxy *SSHProxy) Status() {
	pidPath := proxy.pidPath
	if _, err := os.Stat(pidPath); err == nil {
		fmt.Println("# The SSH Proxy is active.")
	} else {
		fmt.Println("# The SSH Proxy is not active.")
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
		"-H", proxy.kubeconfig.KubeSSHProxy.SSH.Host,
		"-p", strconv.Itoa(proxy.kubeconfig.KubeSSHProxy.SSH.Port),
		"-u", proxy.kubeconfig.KubeSSHProxy.SSH.User,
		"-b", strconv.Itoa(proxy.kubeconfig.KubeSSHProxy.BindPort),
	}
	if proxy.kubeconfig.KubeSSHProxy.SSH.KeyPath != "" {
		args = append(args, "-k", proxy.kubeconfig.KubeSSHProxy.SSH.KeyPath)
	}
	return args
}

func (proxy *SSHProxy) getKubeconfig() {
	var kubeconfig Kubeconfig
	var kubeSSHProxyConfig KubeSSHProxyConfig
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		// TODO probar con mac y winsux
		kubeconfigPath = fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
	}
	yamlFile, err := ioutil.ReadFile(kubeconfigPath)
	CheckGenericError(err)
	err = yaml.Unmarshal(yamlFile, &kubeconfig)
	CheckGenericError(err)
	proxy.kubeconfig = kubeconfig

	kubeSSHProxyConfigPath := os.Getenv("KUBECONFIG-SSH-PROXY")
	if kubeSSHProxyConfigPath == "" {
		// TODO probar con mac y winsux
		kubeSSHProxyConfigPath = fmt.Sprintf("%s-ssh-proxy", kubeconfigPath)
	}
	yamlFile, err = ioutil.ReadFile(kubeSSHProxyConfigPath)
	CheckGenericError(err)
	err = yaml.Unmarshal(yamlFile, &kubeSSHProxyConfig)
	CheckGenericError(err)
	proxy.kubeconfig.KubeSSHProxyConfig = kubeSSHProxyConfig
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
