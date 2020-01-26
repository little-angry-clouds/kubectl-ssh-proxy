package types

import (
	"fmt"
	"os"
)

// SSHProxyConfig is the Kubeconfig section that stores SSHProxy's stuff
type SSHProxyConfig struct {
	SSHProxy struct {
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
	SSHProxyConfig
}

// UnmarshalYAML unmarshals yaml to get the unexistent key CurrentCluster
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
	if aux["kube-ssh-proxy"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy` value.")
		os.Exit(1)
	}
	kubeSshProxyConfig := aux["kube-ssh-proxy"].(map[interface{}]interface{})
	if kubeSshProxyConfig["ssh"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh` value.")
		os.Exit(1)
	}
	kubeSshProxyConfigSsh := kubeSshProxyConfig["ssh"].(map[interface{}]interface{})
	if kubeSshProxyConfigSsh["host"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	if kubeSshProxyConfigSsh["user"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	if kubeSshProxyConfigSsh["port"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	if kubeSshProxyConfig["bind_port"] == nil {
		fmt.Println("Your configuration is incorrect, you're missing the `kube-ssh-proxy.ssh.host` value.")
		os.Exit(1)
	}
	k.SSHProxy.SSH.Host = kubeSshProxyConfigSsh["host"].(string)
	k.SSHProxy.SSH.User = kubeSshProxyConfigSsh["user"].(string)
	k.SSHProxy.SSH.Port = kubeSshProxyConfigSsh["port"].(int)
	k.SSHProxy.BindPort = kubeSshProxyConfig["bind_port"].(int)
	if kubeSshProxyConfigSsh["key_path"] != nil {
		k.SSHProxy.SSH.KeyPath = kubeSshProxyConfigSsh["key_path"].(string)
	}
	return nil
}
