// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_template_deployment

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_group_template_deployment/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceGroupTemplateDeploymentResourceCreate,
		Read:   resourceGroupTemplateDeploymentResourceRead,
		Update: resourceGroupTemplateDeploymentResourceUpdate,
		Delete: resourceGroupTemplateDeploymentResourceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := deployments.ParseResourceGroupProviderDeploymentID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			// v0 -> v1 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			0: migration.ResourceGroupTemplateDeploymentV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		// (@jackofallops - lintignore needed as we need to make sure the JSON is usable in `output_content`)

		// lintignore:S033
		Schema: resourceGroupTemplateDeploymentSchema(),

		// this is needed to fix https://github.com/hashicorp/terraform-provider-azurerm/issues/12828
		// On a change to `template_content` or `parameters_content`, we'll set `output_content` to empty
		// The adverse effect of this is that any change to `template_content` will also cause any resource referencing `output_content` to update
		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			if d.HasChange("template_content") {
				o, n := d.GetChange("template_content")

				// the json has to be normalized and then compared against to see if a change has occurred
				if !strings.EqualFold(o.(string), pluginsdk.NormalizeJson(n)) {
					return d.SetNewComputed("output_content")
				}
			}

			if d.HasChange("parameters_content") {
				o, n := d.GetChange("parameters_content")

				// the json has to be normalized and then compared against to see if a change has occurred
				if !strings.EqualFold(o.(string), pluginsdk.NormalizeJson(n)) {
					return d.SetNewComputed("output_content")
				}
			}

			return nil
		},
	}
}
