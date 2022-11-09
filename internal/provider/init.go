package provider

import (
	"io"
	"log"
	"os"
)

func init() {
	// Disable logging unless debugging, otherwise resource configuration is written to the logs
	val, ok := os.LookupEnv("TF_LOG")
	if !ok || val == "" {
		log.SetOutput(io.Discard)
	}
}
