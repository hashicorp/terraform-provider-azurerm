package util

import log "github.com/sirupsen/logrus"

func InitLogger(debug bool) {
	textFmt := &log.TextFormatter{
		DisableLevelTruncation: true,
		ForceQuote:             true,
	}

	log.SetFormatter(textFmt)
	log.SetLevel(log.WarnLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
