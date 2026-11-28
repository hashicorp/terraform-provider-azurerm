// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_routing_intent

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name virtual_hub_routing_intent -properties "name" -compare-values "subscription_id:virtual_hub_id,resource_group_name:virtual_hub_id,virtual_hub_name:virtual_hub_id"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &virtualwans.RoutingIntentId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_hub_routing_intent"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualwans.ValidateRoutingIntentID
}
