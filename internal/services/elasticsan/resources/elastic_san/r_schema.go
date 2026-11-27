// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san

import (
	"context"
	"errors"
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"base_size_in_tib": {
			Description:  "The base size of the Elastic SAN resource in TiB.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"extended_size_in_tib": {
			Description:  "The extended size of the Elastic SAN resource in TiB.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The SKU name.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						ValidateFunc: validation.StringInSlice(
							elasticsans.PossibleValuesForSkuName(),
							false,
						),
					},

					"tier": {
						Description: "The SKU tier.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"total_mbps": {
			Description: "Total Provisioned MBps of the Elastic SAN resource.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"total_size_in_tib": {
			Description: "Total size of the Elastic SAN resource in TB.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"total_volume_size_in_gib": {
			Description: "Total size of the provisioned Volumes in GiB.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"volume_group_count": {
			Description: "Total number of volume groups in this Elastic SAN resource.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
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
				return errors.New("new base_size_in_tib should be greater than the existing one")
			}

			if oldVal, newVal := metadata.ResourceDiff.GetChange("extended_size_in_tib"); newVal.(int) < oldVal.(int) {
				return errors.New("new extended_size_in_tib should be greater than the existing one")
			}

			return nil
		},
	}
}
