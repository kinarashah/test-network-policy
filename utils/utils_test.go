package utils

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGetWSURL(t *testing.T) {
	expected := "wss://192.168.236.1:8443/k8s/clusters/cluster-k4kxr/api/v1/namespaces/ns-in-non-def-proj/pods/test1-69bc79587b-pk9vt/exec?container=test1&stdout=1&stdin=1&stderr=1&tty=0&command=bash"
	actual := GetWSURL(
		"https://192.168.236.1:8443",
		"cluster-k4kxr",
		"ns-in-non-def-proj",
		"test1-69bc79587b-pk9vt",
		"test1",
		"bash",
	)
	logrus.Infof("e= %v", expected)
	logrus.Infof("a= %v", actual)
	if actual != expected {
		t.Fail()
	}
}

func TestGetFormattedCommand(t *testing.T) {
	expected := "&command=curl&command=--max-time&command=10&command=-s&command=http://test1"
	input := "curl --max-time 10 -s http://test1"
	actual := GetFormattedCommand(input)

	if actual != expected {
		t.Fail()
	}
}
