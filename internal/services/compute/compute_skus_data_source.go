// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ComputeSkusDataSource struct{}

var _ sdk.DataSource = ComputeSkusDataSource{}

type ComputeSkusDataSourceModel struct {
	Name                string                `tfschema:"name"`
	Location            string                `tfschema:"location"`
	IncludeCapabilities bool                  `tfschema:"include_capabilities"`
	Skus                []ComputeSkusSkuModel `tfschema:"skus"`
}

type ComputeSkusSkuModel struct {
	Name                 string            `tfschema:"name"`
	ResourceType         string            `tfschema:"resource_type"`
	Size                 string            `tfschema:"size"`
	Tier                 string            `tfschema:"tier"`
	LocationRestrictions []string          `tfschema:"location_restrictions"`
	ZoneRestrictions     []string          `tfschema:"zone_restrictions"`
	Capabilities         map[string]string `tfschema:"capabilities"`
	Zones                []string          `tfschema:"zones"`
}

func (ds ComputeSkusDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"location": commonschema.Location(),
		"include_capabilities": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (ds ComputeSkusDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"skus": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"resource_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"size": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tier": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"location_restrictions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"zone_restrictions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"capabilities": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"zones": commonschema.ZonesMultipleComputed(),
				},
			},
		},
	}
}

func (ds ComputeSkusDataSource) ModelObject() interface{} {
	return &ComputeSkusDataSourceModel{}
}

func (ds ComputeSkusDataSource) ResourceType() string {
	return "azurerm_compute_skus"
}

func (ds ComputeSkusDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ComputeSkusDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			name := state.Name
			loc := location.Normalize(state.Location)
			availableSkus := make([]ComputeSkusSkuModel, 0)
			id := parse.NewSkusID(subscriptionId)

			resp, err := metadata.Client.Compute.SkusClient.ResourceSkusList(ctx, commonids.NewSubscriptionID(subscriptionId), skus.DefaultResourceSkusListOperationOptions())
			if err != nil {
				return err
			}

			if model := resp.Model; model != nil {
				for _, sku := range *model {
					// the API does not allow filtering by name
					if name != "" {
						if !strings.EqualFold(*sku.Name, name) {
							continue
						}
					}

					// while the API accepts OData filters, the location filter is currently
					// not working, thus we need to filter the results manually
					locationsNormalized := make([]string, len(*sku.Locations))
					for _, v := range *sku.Locations {
						locationsNormalized = append(locationsNormalized, location.Normalize(v))
					}
					if !slices.Contains(locationsNormalized, loc) {
						continue
					}

					var zones []string
					var locationRestrictions []string
					var zoneRestrictions []string
					capabilities := make(map[string]string)

					if sku.Restrictions != nil && len(*sku.Restrictions) > 0 {
						for _, restriction := range *sku.Restrictions {
							restrictionType := *restriction.Type

							switch restrictionType {
							case skus.ResourceSkuRestrictionsTypeLocation:
								restrictedLocationsNormalized := make([]string, 0)
								for _, v := range *restriction.RestrictionInfo.Locations {
									restrictedLocationsNormalized = append(restrictedLocationsNormalized, location.Normalize(v))
								}
								locationRestrictions = restrictedLocationsNormalized

							case skus.ResourceSkuRestrictionsTypeZone:
								zoneRestrictions = *restriction.RestrictionInfo.Zones
							}
						}
					}

					if sku.LocationInfo != nil && len(*sku.LocationInfo) > 0 {
						for _, locationInfo := range *sku.LocationInfo {
							if location.Normalize(*locationInfo.Location) == loc {
								zones = *locationInfo.Zones
							}
						}
					}

					if state.IncludeCapabilities {
						if sku.Capabilities != nil && len(*sku.Capabilities) > 0 {
							for _, capability := range *sku.Capabilities {
								capabilities[*capability.Name] = *capability.Value
							}
						}
					}

					availableSkus = append(availableSkus, ComputeSkusSkuModel{
						Name:                 pointer.From(sku.Name),
						ResourceType:         pointer.From(sku.ResourceType),
						Size:                 pointer.From(sku.Size),
						Tier:                 pointer.From(sku.Tier),
						LocationRestrictions: locationRestrictions,
						ZoneRestrictions:     zoneRestrictions,
						Zones:                zones,
						Capabilities:         capabilities,
					})
				}

				state.Skus = availableSkus
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
