// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

const AzureFirewallPolicyResourceName = "azurerm_firewall_policy"

var AzureFirewallResourceName = "azurerm_firewall"

func RegisterResource() *pluginsdk.Resource {
	resource := pluginsdk.Resource{
		Create:   resourceFirewallCreateUpdate,
		Read:     resourceFirewallRead,
		Update:   resourceFirewallCreateUpdate,
		Delete:   resourceFirewallDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&azurefirewalls.AzureFirewallId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&azurefirewalls.AzureFirewallId{}),
		},

		Schema: firewallSchema(),
	}

	return &resource
}
