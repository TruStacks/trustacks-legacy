package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LokiLogger struct {
	Project  string
	Revision string
	RunID    string
	Host     string
}

type LogMessage struct {
	Streams []LogMessageStream `json:"streams"`
}

type LogMessageStream struct {
	Stream LogMessageStreamStream `json:"stream"`
	Values [][]string             `json:"values"`
}

type LogMessageStreamStream struct {
	Revision string `json:"revision"`
	Project  string `json:"project"`
	RunID    string `json:"runId"`
}

// Write pushes the log message to loki.
func (l *LokiLogger) Write(message []byte) (int, error) {
	// ignore empty(ish) lines
	// if string(message) == "" || string(message) == "\n" {
	// 	return -1, nil
	// }
	logMessage, err := json.Marshal(&LogMessage{
		Streams: []LogMessageStream{
			{
				Stream: LogMessageStreamStream{
					Revision: l.Revision,
					Project:  l.Project,
					RunID:    l.RunID,
				},
				Values: [][]string{
					{
						strconv.FormatInt(time.Now().UnixNano(), 10),
						strings.ReplaceAll(string(message), "\n", ""),
					},
				},
			},
		},
	})
	if err != nil {
		return -1, err
	}
	body := bytes.NewBuffer(logMessage)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/loki/api/v1/push", l.Host), body)
	if err != nil {
		return -1, err
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}
	return len(message), nil
}

// NewLokiLogger .
func NewLokiLogger(project, revision, runID, host string) *LokiLogger {
	return &LokiLogger{
		Project:  project,
		Revision: revision,
		RunID:    runID,
		Host:     host,
	}
}
