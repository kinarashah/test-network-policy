package framework

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewRancherServerFromEnvVars(t *testing.T) {
	_, err := NewRancherServerFromEnvVars()
	if err != nil {
		logrus.Errorf("err: %v", err)
		t.Fail()
	}
}
