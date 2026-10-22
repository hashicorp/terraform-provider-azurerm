// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_hub_local_rulestack

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallresources.ValidateFirewallID
}

func (r Resource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack"
}
