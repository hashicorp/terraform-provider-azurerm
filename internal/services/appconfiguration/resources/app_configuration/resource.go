// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppConfigurationCreate,
		Read:   resourceAppConfigurationRead,
		Update: resourceAppConfigurationUpdate,
		Delete: resourceAppConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurationstores.ParseConfigurationStoreID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// sku cannot be downgraded from a production tier (`premium` or `standard`) to a non-production tier (`developer` or `free`), or downgraded from `developer` to `free`
			// https://learn.microsoft.com/azure/azure-app-configuration/faq#can-i-upgrade-or-downgrade-an-app-configuration-store
			pluginsdk.ForceNewIfChange("sku", func(ctx context.Context, old, new, meta any) bool {
				return ((old == "premium" || old == "standard") && new == "developer") || new == "free"
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
				authMode := d.Get("data_plane_proxy_authentication_mode").(string)
				privLinkDelegation := d.Get("data_plane_proxy_private_link_delegation_enabled").(bool)

				if authMode == string(configurationstores.AuthenticationModeLocal) && privLinkDelegation {
					return errors.New("`data_plane_proxy_private_link_delegation_enabled` cannot be set to `true` when `data_plane_proxy_authentication_mode` is `Local`")
				}

				return nil
			}),
		),

		Schema: appConfigurationSchema(),
	}
}
