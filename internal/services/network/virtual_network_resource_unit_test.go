// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	network "github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
)

func TestVirtualNetwork_AddressSpaceDiffSuppress_IpamUnknown(t *testing.T) {
	rawCfg := cty.ObjectVal(map[string]cty.Value{
		"ip_address_pool": cty.UnknownVal(cty.List(cty.DynamicPseudoType)),
	})

	if !network.ShouldSuppressVnetAddressSpaceDiffForTest(rawCfg) {
		t.Fatalf("expected diff to be suppressed when ip_address_pool is unknown-but-present")
	}
}
