package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type LokiClient struct {
	PushIntveralSeconds int
	// This will also trigger the send event
	MaxBatchSize int
	Values       map[string][][]string
	LokiEndpoint string
	BatchCount   int
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type lokiLogEvent struct {
	Streams []lokiStream
}

func (l *LokiClient) bgRun() {
	lastRunTimestamp := 0
	isWorking := true
	for {
		if time.Now().Second()-lastRunTimestamp > l.PushIntveralSeconds || l.BatchCount > l.MaxBatchSize {
			// Loop over all log levels and send them
			for k, _ := range l.Values {
				if len(l.Values) > 0 {
					prevLogs := l.Values[k]
					l.Values[k] = [][]string{}
					err := pushToLoki(prevLogs, l.LokiEndpoint, k)
					if err != nil && isWorking {
						isWorking = false
						log.Error().Msgf("Logs are currently not being forwarded to loki due to an error: %v", err)
					}
					if err == nil && !isWorking {
						isWorking = true
						// I will not accept PR comments about this log message tyvm
						log.Info().Msgf("Logs are now being published again. The loki instance seems to be reachable once more! May the logeth collecteth'r beest did bless with our logs")
					}
				}
			}
			lastRunTimestamp = time.Now().Second()
			l.BatchCount = 0
		}
	}
}

/*
This function contains *no* error handling/logging because this:
a) should not crash the application
b) would mean that every run of this creates further logs that cannot be published
=> The error will be returned and the problem will be logged ONCE by the handling function
*/
func pushToLoki(logs [][]string, lokiEndpoint string, logLevel string) error {

	lokiPushPath := "/loki/api/v1/push"

	data, err := json.Marshal(lokiLogEvent{
		Streams: []lokiStream{
			{
				Stream: map[string]string{
					"service": "demo",
					"level":   logLevel,
				},
				Values: logs,
			},
		},
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v%v", lokiEndpoint, lokiPushPath), bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(req.Context(), 100*time.Millisecond)

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
