// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceComputeSkus() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceComputeSkusRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"zones": commonschema.ZonesMultipleComputed(),
					},
				},
			},
		},
	}
}

func dataSourceComputeSkusRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SkusClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.ResourceSkusList(ctx, commonids.NewSubscriptionID(subscriptionId), skus.DefaultResourceSkusListOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving SKUs: %+v", err)
	}

	name := d.Get("name").(string)
	loc := location.Normalize(d.Get("location").(string))
	availableSkus := make([]map[string]interface{}, 0)

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

			if d.Get("include_capabilities").(bool) {
				if sku.Capabilities != nil && len(*sku.Capabilities) > 0 {
					for _, capability := range *sku.Capabilities {
						capabilities[*capability.Name] = *capability.Value
					}
				}
			}

			availableSkus = append(availableSkus, map[string]interface{}{
				"name":                  sku.Name,
				"resource_type":         sku.ResourceType,
				"size":                  sku.Size,
				"tier":                  sku.Tier,
				"location_restrictions": locationRestrictions,
				"zone_restrictions":     zoneRestrictions,
				"zones":                 zones,
				"capabilities":          capabilities,
			})
		}
		d.SetId(uuid.New().String())
		d.Set("skus", availableSkus)
	}

	return nil
}
