// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
)

const EnvGracefulStop = "ARM_TEST_GRACEFUL_STOP"

var (
	gracefulStopOnce      sync.Once
	gracefulStopRequested atomic.Bool
)

func initGracefulStop(t *testing.T) {
	gracefulStopOnce.Do(func() {
		if !gracefulStopEnabled() {
			return
		}

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM)

		go func() {
			sig := <-ch
			gracefulStopRequested.Store(true)
			signal.Reset(syscall.SIGTERM)
			t.Logf("[WARN] SIGTERM received: finishing in-progress tests and skipping the rest; send it again to force stop", sig)
		}()
	})
}

func gracefulStopEnabled() bool {
	enabled, err := strconv.ParseBool(os.Getenv(EnvGracefulStop))
	return err == nil && enabled
}

func skipIfGracefulStopRequested(t *testing.T) {
	if gracefulStopRequested.Load() {
		t.Skipf("graceful stop requested: skipping before provisioning resources for test %s", t.Name())
	}
}
