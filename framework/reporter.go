package framework

import (
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters/stenographer"
)

func NewReporter () reporters.Reporter {
	config.DefaultReporterConfig.SlowSpecThreshold = 60
	c := config.DefaultReporterConfig
	c.SlowSpecThreshold = 60
	s:= stenographer.New(!config.DefaultReporterConfig.NoColor, config.GinkgoConfig.FlakeAttempts > 1)
	return reporters.NewDefaultReporter(c, s)
}
