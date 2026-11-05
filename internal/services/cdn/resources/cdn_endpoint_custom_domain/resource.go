// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_endpoint_custom_domain/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmCdnEndpointCustomDomainCreate,
		Read:   resourceArmCdnEndpointCustomDomainRead,
		Update: resourceArmCdnEndpointCustomDomainUpdate,
		Delete: resourceArmCdnEndpointCustomDomainDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := customdomains.ParseEndpointCustomDomainID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnEndpointCustomDomainV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(12 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(24 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(12 * time.Hour),
		},

		Schema: cdnEndpointCustomDomainSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			if helpers.IsCdnFullyRetired() {
				return fmt.Errorf("%s", helpers.FullyRetiredMessage)
			}

			if helpers.IsCdnDeprecatedForCreation() && d.HasChanges("name", "cdn_endpoint_id", "host_name") {
				return fmt.Errorf("%s", helpers.CreateDeprecationMessage)
			}

			return nil
		}),
	}
}
