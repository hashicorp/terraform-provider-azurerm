// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/resources/frontdoor_firewall_policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorFirewallPolicyCreateUpdate,
		Read:   resourceFrontDoorFirewallPolicyRead,
		Update: resourceFrontDoorFirewallPolicyCreateUpdate,
		Delete: resourceFrontDoorFirewallPolicyDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebApplicationFirewallPolicyV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: frontdoorFirewallPolicySchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			if helpers.IsFrontDoorFullyRetired() {
				return fmt.Errorf("%s", helpers.FullyRetiredMessage)
			}

			// New resources are not supported, and since these fields are 'ForceNew' we also need to block changing them as
			// the re-create would fail with the create error from the service API...
			if helpers.IsFrontDoorDeprecatedForCreation() && d.HasChanges("name", "resource_group_name") {
				return fmt.Errorf("%s", helpers.CreateDeprecationMessage)
			}

			return nil
		}),
	}
}
