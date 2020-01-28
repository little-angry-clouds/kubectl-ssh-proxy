package main

import (
	. "github.com/little-angry-clouds/kubectl-ssh-proxy/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

// Set base suite
type Suite struct {
	suite.Suite
	sshProxy SSHProxy
}

func (suite *Suite) SetupTest() {
	os.Setenv("KUBECONFIG", "./test_data/example.yml")
	suite.sshProxy = SSHProxy{}
	suite.sshProxy.getKubeconfig()
	os.Setenv("KUBECONFIG", "")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// Start tests
func (suite *Suite) TestGetKubeconfig() {
	var kubeconfig Kubeconfig
	kubeconfig.CurrentCluster = "MyCluster"
	kubeconfig.CurrentContext = "default"
	kubeconfig.KubeSSHProxy.SSH.Host = "anywhere"
	kubeconfig.KubeSSHProxy.SSH.Port = 22
	kubeconfig.KubeSSHProxy.SSH.User = "nobody"
	kubeconfig.KubeSSHProxy.SSH.KeyPath = "/home/nobody/.ssh/nobody"
	kubeconfig.KubeSSHProxy.BindPort = 8080
	assert.Equal(suite.T(), kubeconfig, suite.sshProxy.kubeconfig, "they should be equal")
}

func (suite *Suite) TestCreateArgs() {
	args := suite.sshProxy.createArgs()
	expectedArgs := "-H anywhere -p 22 -u nobody -b 8080 -k /home/nobody/.ssh/nobody"
	assert.Equal(suite.T(), expectedArgs, strings.Join(args[:], " "), "they should be equal")
}

func (suite *Suite) TestGetPidPath() {
	suite.sshProxy.getPidPath()
	assert.Equal(suite.T(), "/run/user/1000/kubectl-ssh-proxy/MyCluster/PID", suite.sshProxy.pidPath, "they should be equal")
}

func (suite *Suite) TestSSHProxyStopNotActive() {
	suite.sshProxy.getPidPath()
	pidPath := suite.sshProxy.pidPath
	emptyFile, _ := os.Create(pidPath)
	defer emptyFile.Close()
	suite.sshProxy.Stop()
	_, err := os.Stat(pidPath)
	assert.NotEqual(suite.T(), err, nil, "they should not be equal")
}

func (suite *Suite) TestSSHProxyStatusNotActive() {
	suite.sshProxy.getPidPath()
	pidPath := suite.sshProxy.pidPath
	emptyFile, _ := os.Create(pidPath)
	defer emptyFile.Close()
	suite.sshProxy.Status()
	os.Remove(pidPath)
}

func (suite *Suite) TestSSHProxyStatusNotActivated() {
	suite.sshProxy.Status()
}
