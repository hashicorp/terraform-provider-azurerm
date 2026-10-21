// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/resources/frontdoor_custom_https_configuration/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorCustomHTTPSConfigurationCreateUpdate,
		Read:   resourceFrontDoorCustomHTTPSConfigurationRead,
		Update: resourceFrontDoorCustomHTTPSConfigurationCreateUpdate,
		Delete: resourceFrontDoorCustomHTTPSConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.CustomHttpsConfigurationID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Frontdoor.FrontDoorsClient

			// validate that the passed ID is a valid custom HTTPS configuration ID
			custom, err := parse.CustomHttpsConfigurationID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing Custom HTTPS Configuration ID %q for import: %v", d.Id(), err)
			}

			// convert the passed custom HTTPS configuration ID to a frontend endpoint ID
			frontend := frontdoors.NewFrontendEndpointID(custom.SubscriptionId, custom.ResourceGroup, custom.FrontDoorName, custom.CustomHttpsConfigurationName)

			// validate that the frontend endpoint ID exists in the Frontdoor resource
			if _, err = client.FrontendEndpointsGet(ctx, frontend); err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving the Custom HTTPS Configuration(ID: %q) for the frontend endpoint (ID: %q): %s", custom.ID(), frontend.ID(), err)
			}

			// set the new values for the custom HTTPS configuration resource
			d.Set("id", custom.ID())
			d.Set("frontend_endpoint_id", frontend.ID())

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: frontdoorCustomHttpsConfigurationSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(customizeHttpsConfigurationCustomizeDiff),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CustomHttpsConfigurationV0ToV1{},
		}),
	}
}
