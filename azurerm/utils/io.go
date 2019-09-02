package utils

import (
	"io"
	"log"
)

// TODO: deprecate / remove me

func IoCloseAndLogError(c io.Closer, message string) {
	if err := c.Close(); err != nil {
		log.Printf("%s: %v", message, err)
	}
}
