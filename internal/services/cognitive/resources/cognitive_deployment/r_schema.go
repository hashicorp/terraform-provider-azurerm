// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_deployment

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"dynamic_throttling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"model": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"format": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Standard",
							"DataZoneBatch",
							"DataZoneProvisionedManaged",
							"DataZoneStandard",
							"GlobalBatch",
							"GlobalProvisionedManaged",
							"GlobalStandard",
							"ProvisionedManaged",
						}, false),
					},

					"tier": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(deployments.PossibleValuesForSkuTier(), false),
					},

					"size": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"family": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"rai_policy_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C `raiPolicyName` has default value when `rai_policy_name` is not set. So, `O+C` is required otherwise it will incur difference.
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spillover_deployment_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_upgrade_option": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(deployments.DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable),
			ValidateFunc: validation.StringInSlice(deployments.PossibleValuesForDeploymentModelVersionUpgradeOption(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config cognitiveDeploymentModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if config.SpilloverDeploymentName == "" || len(config.Sku) == 0 {
				return nil
			}

			spilloverSupportedSkuNames := []string{
				"DataZoneProvisionedManaged",
				"GlobalProvisionedManaged",
				"ProvisionedManaged",
			}

			if skuName := config.Sku[0].Name; !slices.Contains(spilloverSupportedSkuNames, skuName) {
				return fmt.Errorf("`spillover_deployment_name` can only be set when `sku.0.name` is `DataZoneProvisionedManaged`, `GlobalProvisionedManaged`, or `ProvisionedManaged`, got `%s`", skuName)
			}

			return nil
		},
	}
}
