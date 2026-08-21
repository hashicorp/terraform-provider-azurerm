// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package timeouts

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestForReadPrefersRegisteredDataSourceReadTimeout(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})

	RegisterDataSourceReadTimeout(d, 5*time.Minute)
	defer DeregisterDataSourceReadTimeout(d)

	ctx, cancel := ForRead(context.Background(), d)
	defer cancel()

	assertDeadlineApprox(t, ctx, 5*time.Minute)
}

func TestForReadFallsBackToResourceDataTimeout(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})

	ctx, cancel := ForRead(context.Background(), d)
	defer cancel()

	// with nothing registered `d.Timeout` returns the SDK-global default of 20 minutes
	assertDeadlineApprox(t, ctx, 20*time.Minute)
}

func TestForReadAfterDeregisterFallsBack(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})

	RegisterDataSourceReadTimeout(d, 5*time.Minute)
	DeregisterDataSourceReadTimeout(d)

	ctx, cancel := ForRead(context.Background(), d)
	defer cancel()

	assertDeadlineApprox(t, ctx, 20*time.Minute)
}

func assertDeadlineApprox(t *testing.T, ctx context.Context, expected time.Duration) {
	t.Helper()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline to be set on the context but none was")
	}
	remaining := time.Until(deadline)
	if remaining > expected || remaining < expected-time.Minute {
		t.Fatalf("expected a deadline of ~%s from now but got %s", expected, remaining)
	}
}
