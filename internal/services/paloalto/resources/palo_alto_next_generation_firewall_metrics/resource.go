// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_metrics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id firewall_id -test-expect-non-empty

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate               = Resource{}
	_ sdk.ResourceWithIdentityTypeOverride = Resource{}
)

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return metricsobjectfirewallresources.ValidateFirewallID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &metricsobjectfirewallresources.FirewallId{}
}

func (r Resource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityVirtual
}

func (r Resource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_metrics"
}
