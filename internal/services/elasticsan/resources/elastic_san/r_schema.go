// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Elastic SAN resource.",
			Required:     true,
			ForceNew:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validate.ElasticSanName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"base_size_in_tib": {
			Description:  "The base size of the Elastic SAN resource in TiB.",
			Required:     true,
			Type:         pluginsdk.TypeInt,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"extended_size_in_tib": {
			Description:  "The extended size of the Elastic SAN resource in TiB.",
			Optional:     true,
			Type:         pluginsdk.TypeInt,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"sku": {
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The SKU name.",
						Required:    true,
						ForceNew:    true,
						Type:        pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							elasticsans.PossibleValuesForSkuName(),
							false,
						),
					},

					"tier": {
						Description: "The SKU tier.",
						Optional:    true,
						Type:        pluginsdk.TypeString,
						Default:     string(elasticsans.SkuTierPremium),
						ValidateFunc: validation.StringInSlice(
							elasticsans.PossibleValuesForSkuTier(),
							false,
						),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"total_iops": {
			Description: "Total Provisioned IOPS of the Elastic SAN resource.",
			Computed:    true,
			Type:        pluginsdk.TypeInt,
		},

		"total_mbps": {
			Description: "Total Provisioned MBps of the Elastic SAN resource.",
			Computed:    true,
			Type:        pluginsdk.TypeInt,
		},

		"total_size_in_tib": {
			Description: "Total size of the Elastic SAN resource in TB.",
			Computed:    true,
			Type:        pluginsdk.TypeInt,
		},

		"total_volume_size_in_gib": {
			Description: "Total size of the provisioned Volumes in GiB.",
			Computed:    true,
			Type:        pluginsdk.TypeInt,
		},

		"volume_group_count": {
			Description: "Total number of volume groups in this Elastic SAN resource.",
			Computed:    true,
			Type:        pluginsdk.TypeInt,
		},
	}
}

func (k Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config ElasticSANResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(config.Zones) > 0 && len(config.Sku) > 0 && config.Sku[0].Name == string(elasticsans.SkuNamePremiumZRS) {
				return fmt.Errorf("zones are not supported for the %s SKU", elasticsans.SkuNamePremiumZRS)
			}

			if oldVal, newVal := metadata.ResourceDiff.GetChange("base_size_in_tib"); newVal.(int) < oldVal.(int) {
				return fmt.Errorf("new base_size_in_tib should be greater than the existing one")
			}

			if oldVal, newVal := metadata.ResourceDiff.GetChange("extended_size_in_tib"); newVal.(int) < oldVal.(int) {
				return fmt.Errorf("new extended_size_in_tib should be greater than the existing one")
			}

			return nil
		},
	}
}
