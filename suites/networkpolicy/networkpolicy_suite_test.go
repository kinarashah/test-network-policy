package networkpolicy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rancher/test-network-policy/framework"
	"github.com/onsi/ginkgo/config"
)

var RancherServer *framework.RancherServer

func TestNetworkpolicy(t *testing.T) {
	RegisterFailHandler(Fail)
	config.DefaultReporterConfig.SlowSpecThreshold = 60
	config.DefaultReporterConfig.Verbose = true
	RunSpecs(t, "Networkpolicy Suite")
}

var _ = BeforeSuite(func() {
	//logrus.Infof("BeforeSuite")
	var err error

	RancherServer, err = framework.NewRancherServerFromEnvVars()
	Expect(err).NotTo(HaveOccurred(), "while creating rancher server")
})

var _ = AfterSuite(func() {
	//logrus.Infof("AfterSuite")
})
