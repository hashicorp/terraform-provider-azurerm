// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name network_manager_routing_rule -properties "name" -compare-values "subscription_id:rule_collection_id,resource_group_name:rule_collection_id,network_manager_name:rule_collection_id,routing_configuration_name:rule_collection_id,rule_collection_name:rule_collection_id"

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
)

type Resource struct{}

// Identity implements [sdk.ResourceWithIdentity].
func (r Resource) Identity() resourceids.ResourceId {
	return &routingrules.RuleId{}
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return routingrules.ValidateRuleID
}

func (Resource) ResourceType() string {
	return "azurerm_network_manager_routing_rule"
}
