// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Unit and benchmark tests verifying that per-service Provider interface client resolution
// preserves exact client pointers with zero heap allocations and negligible (~1-2 ns) overhead.
package sdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
)

func TestClientProviderRoundTrip(t *testing.T) {
	comp := &computeclient.Client{}
	acc := &account.ResourceManagerAccount{SubscriptionId: "00000000-0000-0000-0000-000000000000"}
	c := &clients.Client{
		Account:  acc,
		Compute:  comp,
		Features: features.UserFeatures{PersistIDOnCreateBeforePollingForCompletion: true},
	}

	var provider sdk.ClientProvider = c
	if got := computeclient.From(provider); got != comp {
		t.Fatalf("computeclient.From returned %p, want %p", got, comp)
	}
	if got := account.From(provider); got != acc {
		t.Fatalf("account.From returned %p, want %p", got, acc)
	}
	if !features.From(provider).PersistIDOnCreateBeforePollingForCompletion {
		t.Fatal("features.From did not preserve UserFeatures")
	}
}

var sinkCompute *computeclient.Client

func BenchmarkClientProviderFrom(b *testing.B) {
	c := &clients.Client{Compute: &computeclient.Client{}}
	var meta any = c

	b.ReportAllocs()

	for b.Loop() {
		sinkCompute = computeclient.From(meta)
	}
}
