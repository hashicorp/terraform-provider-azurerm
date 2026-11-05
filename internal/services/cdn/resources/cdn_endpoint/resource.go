// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_endpoint/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnEndpointCreate,
		Read:   resourceCdnEndpointRead,
		Update: resourceCdnEndpointUpdate,
		Delete: resourceCdnEndpointDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnEndpointV0ToV1{},
			// v1 -> v2 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			1: migration.CdnEndpointV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := endpoints.ParseEndpointID(id)
			return err
		}),

		Schema: cdnEndpointSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			if helpers.IsCdnFullyRetired() {
				return fmt.Errorf("%s", helpers.FullyRetiredMessage)
			}

			if helpers.IsCdnDeprecatedForCreation() && d.HasChanges("name", "profile_name", "origin") {
				return fmt.Errorf("%s", helpers.CreateDeprecationMessage)
			}

			return nil
		}),
	}
}
