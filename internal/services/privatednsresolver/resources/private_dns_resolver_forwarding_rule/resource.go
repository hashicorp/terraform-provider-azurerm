// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_forwarding_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:dns_forwarding_ruleset_id,resource_group_name:dns_forwarding_ruleset_id,dns_forwarding_ruleset_name:dns_forwarding_ruleset_id"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &forwardingrules.ForwardingRuleId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_private_dns_resolver_forwarding_rule"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return forwardingrules.ValidateForwardingRuleID
}
