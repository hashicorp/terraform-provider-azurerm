// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_application_rule_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_policy_rule_collection_group"
)

var (
	AzureFirewallPolicyResourceName = firewall_policy_rule_collection_group.AzureFirewallPolicyResourceName
	AzureFirewallResourceName       = firewall_application_rule_collection.AzureFirewallResourceName
)
