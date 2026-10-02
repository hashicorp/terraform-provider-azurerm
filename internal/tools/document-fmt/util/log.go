// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package util

import "github.com/sirupsen/logrus"

func InitLogger(debug bool) {
	textFmt := &logrus.TextFormatter{
		DisableLevelTruncation: true,
		ForceQuote:             true,
	}

	logrus.SetFormatter(textFmt)
	logrus.SetLevel(logrus.WarnLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
