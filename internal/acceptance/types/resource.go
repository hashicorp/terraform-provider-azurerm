// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TestResource interface {
	Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error)
}

type TestResourceVerifyingRemoved interface {
	TestResource
	Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error)
}

// TestResourceVerifyingRemovedWithCustomDestroyTimeout is an optional interface which can be implemented
// by a TestResourceVerifyingRemoved to override the default acceptance-test delete deadline used by
// internal/acceptance/helpers.DeleteResourceFunc.
//
// If not implemented (or if DestroyTimeout returns <= 0), the default deadline remains 1 hour.
type TestResourceVerifyingRemovedWithCustomDestroyTimeout interface {
	// DestroyTimeout returns the maximum duration allowed for the Destroy call.
	DestroyTimeout() time.Duration
}
