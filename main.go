package main

import (
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lokiClient LokiClient

func init() {
	lokiClient = LokiClient{
		PushIntveralSeconds: 10,  // Threshhold of 10s
		MaxBatchSize:        500, //Threshold of 500 events
		LokiEndpoint:        "http://localhost:3100",
	}

	go lokiClient.bgRun()
	log.Logger = log.Hook(LokiHook{})
}

func main() {

	for {
		log.Info().Msg("Sample log message")
		time.Sleep(1 * time.Second)
	}
}

type LokiHook struct {
}

func (h LokiHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	lokiClient.Values = append(lokiClient.Values, []string{strconv.FormatInt(time.Now().UnixNano(), 10), msg})
}
