package types

import (
	"fmt"
	"os"
)

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

type Kubeconfig struct {
	CurrentCluster string
	CurrentContext string `yaml:"current-context"`
	Contexts       []struct {
		Name string `yaml:"name"`
	} `yaml:"context"`
	SSHProxyConfig
}

func (k *Kubeconfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var aux map[string]interface{}
	if unmarshal(&aux); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	k.CurrentContext = aux["current-context"].(string)
	// Search the name of the cluster of the current context
	for key, _ := range aux {
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
	kubeSshProxyConfig := aux["kube-ssh-proxy"].(map[interface{}]interface{})
	kubeSshProxyConfigSsh := kubeSshProxyConfig["ssh"].(map[interface{}]interface{})
	k.SSHProxy.BindPort = kubeSshProxyConfig["bind_port"].(int)
	k.SSHProxy.SSH.Host = kubeSshProxyConfigSsh["host"].(string)
	k.SSHProxy.SSH.User = kubeSshProxyConfigSsh["user"].(string)
	k.SSHProxy.SSH.Port = kubeSshProxyConfigSsh["port"].(int)
	k.SSHProxy.SSH.KeyPath = kubeSshProxyConfigSsh["key_path"].(string)
	return nil
}
