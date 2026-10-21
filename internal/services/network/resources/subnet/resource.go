// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var VirtualNetworkResourceName = "azurerm_virtual_network"

var SubnetResourceName = "azurerm_subnet"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceSubnetCreate,
		Read:     resourceSubnetRead,
		Update:   resourceSubnetUpdate,
		Delete:   resourceSubnetDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.SubnetId{}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				// Validate `sharing_scope` cannot be set when `default_outbound_access_enabled` is true.
				if diff.Get("sharing_scope").(string) != "" && diff.Get("default_outbound_access_enabled").(bool) {
					return fmt.Errorf("`sharing_scope` cannot be set if `default_outbound_access_enabled` is set to `true`")
				}
				return nil
			}),
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.SubnetId{}),
		},

		Schema: subnetSchema(),
	}
}
