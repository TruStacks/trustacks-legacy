package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type LokiLogger struct {
	Project  string
	Revision string
	Host     string
}

func (l *LokiLogger) Write(data []byte) (int, error) {
	// ignore empty(ish) lines
	if string(data) == "" || string(data) == "\n" {
		return 0, nil
	}
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
		"streams": [
		  {
			"stream": {
				"revision": "%s",
				"project": "%s"
			},
			"values": [
			  ["%d", "%s"]
			]
		  }
		]
	}`, l.Revision, l.Project, time.Now().UnixNano(), strings.ReplaceAll(string(data), "\n", ""))))
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/loki/api/v1/push", l.Host), body)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func NewLokiLogger(project, revision, host string) *LokiLogger {
	return &LokiLogger{project, revision, host}
}
