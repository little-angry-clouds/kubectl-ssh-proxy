package sshproxy

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
	os.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	suite.sshProxy.getPidPath()
	os.Setenv("XDG_RUNTIME_DIR", "")
	assert.Equal(suite.T(), "/run/user/1000/kubectl-ssh-proxy/MyCluster/PID", suite.sshProxy.pidPath, "they should be equal")
}

func (suite *Suite) TestSSHProxyStatusActive() {
	os.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	expectedMessage := "# The SSH Proxy is active."
	os.MkdirAll("/run/user/1000/kubectl-ssh-proxy/MyCluster/", os.ModePerm)
	emptyFile, _ := os.Create("/run/user/1000/kubectl-ssh-proxy/MyCluster/PID")
	defer emptyFile.Close()
	suite.sshProxy.getPidPath()
	message := suite.sshProxy.Status()
	os.Remove("/run/user/1000/kubectl-ssh-proxy/MyCluster/PID")
	os.Setenv("XDG_RUNTIME_DIR", "")
	assert.Equal(suite.T(), expectedMessage, message, "they should be equal")
}

func (suite *Suite) TestSSHProxyStatusNotActive() {
	os.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	expectedMessage := "# The SSH Proxy is not active."
	message := suite.sshProxy.Status()
	os.Setenv("XDG_RUNTIME_DIR", "")
	assert.Equal(suite.T(), expectedMessage, message, "they should be equal")
}
