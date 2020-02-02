package sshproxy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"

	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/types"
	"gopkg.in/yaml.v2"
)

// SSHProxy is the main object
type SSHProxy struct {
	kubeconfig Kubeconfig
	pidPath    string
}

// Init initializes the SSHProxy object
func (proxy *SSHProxy) Init() error {
	var err error
	err = proxy.getKubeconfig()
	if err != nil {
		return err
	}
	proxy.getPidPath()
	return nil
}

// Start starts the SSHProxy
func (proxy *SSHProxy) Start() error {
	var err error
	pidPath := proxy.pidPath
	pidDir := path.Dir(pidPath)
	if _, err := os.Stat(pidPath); err == nil {
		return err
	}
	if _, err := os.Stat(pidDir); os.IsNotExist(err) {
		err = os.MkdirAll(pidDir, 0755)
		return err
	}
	args := proxy.createArgs()
	cmd := exec.Command("kube-ssh-proxy-ssh-bin", args...)
	err = cmd.Start()
	if err != nil {
		return err
	}
	// Capture the state of the subcommand. To do it it's necessary to add a little sleep
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:
		message := fmt.Sprintf("# The ssh proxy failed. The error is: %s", err)
		fmt.Println(message)
		message = fmt.Sprintf("# You may debug the error by executing the ssh binary manually: %s", cmd.String())
		fmt.Println(message)
		os.Exit(1)
	case <-time.After(1000 * time.Millisecond):
		// timeout, which means that all works correctly
	}
	pid := []byte(strconv.Itoa(cmd.Process.Pid))
	err = ioutil.WriteFile(pidPath, pid, 0644)
	if err != nil {
		return err
	}
	fmt.Println("# The SSH Proxy started!")
	if runtime.GOOS == "windows" {
		fmt.Println("# Eval the next: \nset HTTPS_PROXY=socks5://localhost:8080")
	} else {
		fmt.Println("# Eval the next: \nexport HTTPS_PROXY=socks5://localhost:8080")
	}
	return nil
}

// Stop stops the SSHProxy
func (proxy *SSHProxy) Stop() error {
	pidPath := proxy.pidPath
	if _, err := os.Stat(pidPath); err != nil {
		fmt.Println("# The ssh proxy is already stopped!")
		os.Exit(0)
	}
	file, err := os.Open(pidPath)
	if err != nil {
		return err
	}
	defer file.Close()
	pid, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	p, err := strconv.Atoi(string(pid))
	if err != nil {
		return err
	}
	process, err := os.FindProcess(p)
	if err != nil {
		return err
	}
	process.Signal(os.Interrupt)
	if runtime.GOOS == "windows" {
		fmt.Println("# The SSH Proxy is already stopped! Eval the next:\nset HTTPS_PROXY=")
	} else {
		fmt.Println("# The SSH Proxy is already stopped! Eval the next:\nunset HTTPS_PROXY")
	}
	os.Remove(pidPath)
	return nil
}

// Status gets the SSHProxy status
func (proxy *SSHProxy) Status() string {
	var message string
	pidPath := proxy.pidPath
	if _, err := os.Stat(pidPath); err == nil {
		message = "# The SSH Proxy is active."
	} else {
		message = "# The SSH Proxy is not active."
	}
	return message
}

func (proxy *SSHProxy) getPidPath() {
	var path string
	if runtime.GOOS == "linux" {
		path = os.Getenv("XDG_RUNTIME_DIR")
	} else if runtime.GOOS == "darwin" {
		// Not sure if XDG_RUNTIME_DIR exists, so just to be sure let's set it
		path = "/tmp/"
	} else if runtime.GOOS == "windows" {
		path = "%USERPROFILE%\\AppData\\Local\\Temp"
	}
	pidDir := fmt.Sprintf("%s/kubectl-ssh-proxy/%s", path, proxy.kubeconfig.CurrentCluster)
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

func (proxy *SSHProxy) getKubeconfig() error {
	var kubeconfig Kubeconfig
	var kubeSSHProxyConfig KubeSSHProxyConfig
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		kubeconfigPath = fmt.Sprintf("%s/.kube/config", home)
	}
	yamlFile, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &kubeconfig)
	if err != nil {
		return err
	}
	proxy.kubeconfig = kubeconfig
	kubeSSHProxyConfigPath := os.Getenv("KUBECONFIG-SSH-PROXY")
	if kubeSSHProxyConfigPath == "" {
		kubeSSHProxyConfigPath = fmt.Sprintf("%s-ssh-proxy", kubeconfigPath)
	}
	yamlFile, err = ioutil.ReadFile(kubeSSHProxyConfigPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &kubeSSHProxyConfig)
	if err != nil {
		return err
	}
	proxy.kubeconfig.KubeSSHProxyConfig = kubeSSHProxyConfig
	return nil
}
