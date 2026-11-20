// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"context"
	"fmt"
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

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
			if v, ok := d.GetOk("autoscale"); ok {
				autoscale := v.([]any)
				if len(autoscale) == 1 && autoscale[0] != nil {
					configuration := autoscale[0].(map[string]any)
					if configuration["maximum_capacity"].(int) < configuration["minimum_capacity"].(int) {
						return fmt.Errorf("`autoscale.maximum_capacity` (%d) must be greater than or equal to `autoscale.minimum_capacity` (%d)", configuration["maximum_capacity"].(int), configuration["minimum_capacity"].(int))
					}
				}
			}
			return nil
		}),

		Schema: firewallSchema(),
	}

	return &resource
}
