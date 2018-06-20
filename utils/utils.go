package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/base64"
	"github.com/gorilla/websocket"
)

func GetWSURL(url, clusterID, podNS, podName, containerName, command string) string {
	s := strings.Replace(url, "https", "wss", 1)
	fmtCmd := GetFormattedCommand(command)
	wsURLTemplate := "%v/k8s/clusters/%v/api/v1/namespaces/%v/pods/%v/exec?container=%v&stdout=1&stdin=1&stderr=1&tty=0%v"
	wsURL := fmt.Sprintf(wsURLTemplate, s, clusterID, podNS, podName, containerName, fmtCmd)
	return wsURL
}

func GetFormattedCommand(command string) string {
	var result string
	splits := strings.Split(command, " ")
	for _, split := range splits {
		result = result + "&command=" + split
	}
	return result
}

func RunExecCommand(wsURL, username, password, token string) (string, error) {
	var data []byte
	var credentials string

	if username != "" && password != "" {
		credentials = username + ":" + password
	} else if token != "" {
		credentials = token
	} else {
		return "", fmt.Errorf("login credentials not provided")
	}
	h := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))}}

	d := &websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
	c, _, err := d.Dial(wsURL, h)
	if err != nil {
		return "", err
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			data = append(data, message...)
		}
	}()

	// TODO: Timeout
	for {
		select {
		case <-done:
			return string(data), nil
		}
	}
}
