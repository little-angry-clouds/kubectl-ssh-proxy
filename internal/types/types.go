package types

import (
	"fmt"
	"os"
)

// KubeSSHProxyConfig is the Kubeconfig section that stores SSHProxy's stuff
type KubeSSHProxyConfig struct {
	KubeSSHProxy struct {
		SSH struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
			User string `yaml:"user"`
			// TODO añadir soporte para ssh-agent y contraseñas
			KeyPath string `yaml:"key_path"`
		} `yaml:"ssh"`
		BindPort int `yaml:"bind_port"`
	}
}

// Kubeconfig stores the relevant Kubeconfig information
type Kubeconfig struct {
	CurrentCluster string
	CurrentContext string `yaml:"current-context"`
	Contexts       []struct {
		Name string `yaml:"name"`
	} `yaml:"context"`
	KubeSSHProxyConfig
}

// UnmarshalYAML unmarshals yaml to match kubeconfig config
func (k *Kubeconfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var aux map[string]interface{}
	if unmarshal(&aux); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if aux["current-context"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing a value.")
		os.Exit(1)
	}
	k.CurrentContext = aux["current-context"].(string)
	// Search the name of the cluster of the current context
	for key := range aux {
		if key == "contexts" {
			c := aux[key].([]interface{})
			for _, v := range c {
				n := v.(map[interface{}]interface{})
				if n["name"] == k.CurrentContext {
					l := n["context"].(map[interface{}]interface{})
					k.CurrentCluster = l["cluster"].(string)
					break
				}
			}
		}
	}
	return nil
}

// UnmarshalYAML unmarshals yaml to match kube-ssh-proxy config
func (k *KubeSSHProxyConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var aux map[string]interface{}
	if unmarshal(&aux); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if aux["kube-ssh-proxy"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy` value.")
		os.Exit(1)
	}
	kubeSSHProxyConfig := aux["kube-ssh-proxy"].(map[interface{}]interface{})
	if kubeSSHProxyConfig["ssh"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh` value.")
		os.Exit(1)
	}
	kubeSSHProxyConfigSSH := kubeSSHProxyConfig["ssh"].(map[interface{}]interface{})
	if kubeSSHProxyConfigSSH["host"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	if kubeSSHProxyConfigSSH["user"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	if kubeSSHProxyConfigSSH["port"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	if kubeSSHProxyConfig["bind_port"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	k.KubeSSHProxy.SSH.Host = kubeSSHProxyConfigSSH["host"].(string)
	k.KubeSSHProxy.SSH.User = kubeSSHProxyConfigSSH["user"].(string)
	k.KubeSSHProxy.SSH.Port = kubeSSHProxyConfigSSH["port"].(int)
	k.KubeSSHProxy.BindPort = kubeSSHProxyConfig["bind_port"].(int)
	if kubeSSHProxyConfigSSH["key_path"] != nil {
		k.KubeSSHProxy.SSH.KeyPath = kubeSSHProxyConfigSSH["key_path"].(string)
	}
	return nil
}
