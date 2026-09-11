// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"io"
	"log"
	"time"
)

var (
	defaultTimeout = 10 * time.Second
	discardLogger  = log.New(io.Discard, "", 0)
)

type fileCheckFunc func(path string) error
