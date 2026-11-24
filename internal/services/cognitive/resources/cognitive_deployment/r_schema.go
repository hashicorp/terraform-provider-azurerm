// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Cognitive Services Account Deployment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Description:  "The ID of the Cognitive Services Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"dynamic_throttling_enabled": {
			Description: "Whether dynamic throttling is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"model": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"format": {
						Description:  "The format of the Cognitive Services Account Deployment model.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"name": {
						Description:  "The name of the Cognitive Services Account Deployment model.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Description: "The version of Cognitive Services Account Deployment model. If `version` is not specified, the default version of the model at the time will be assigned.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description: "The name of the SKU.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
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
						Description:  "This property is required only when multiple tiers are available with the SKU name.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(deployments.PossibleValuesForSkuTier(), false),
					},

					"size": {
						Description: "The SKU size. When the name field is the combination of tier and some other value, this would be the standalone code.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},

					"family": {
						Description: "If the service has different generations of hardware, for the same SKU, then that can be captured here.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},

					"capacity": {
						Description:  "Tokens-per-Minute (TPM). The unit of measure for this field is in the thousands of Tokens-per-Minute. If the resources SKU supports scale in/out then the capacity field should be included in the resources' configuration. If the scale in/out is not supported by the resources SKU then this field can be safely omitted. For more information about TPM please see the [product documentation](https://learn.microsoft.com/azure/ai-services/openai/how-to/quota?tabs=rest).",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"rai_policy_name": {
			Description: "The name of RAI policy.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C `raiPolicyName` has default value when `rai_policy_name` is not set. So, `O+C` is required otherwise it will incur difference.
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spillover_deployment_name": {
			Description:  "The name of the deployment that should serve requests when this deployment would have otherwise been throttled due to reaching its throughput limit.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_upgrade_option": {
			Description:  "Deployment model version upgrade option.",
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
