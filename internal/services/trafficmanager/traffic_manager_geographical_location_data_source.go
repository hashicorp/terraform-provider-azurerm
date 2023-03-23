package trafficmanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/geographichierarchies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmTrafficManagerGeographicalLocation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmTrafficManagerGeographicalLocationRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceArmTrafficManagerGeographicalLocationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.GeographialHierarchiesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	results, err := client.GetDefault(ctx)
	if err != nil {
		return fmt.Errorf("loading Traffic Manager Geographical Hierarchies: %+v", err)
	}

	name := d.Get("name").(string)

	var result *geographichierarchies.Region
	if model := results.Model; model != nil {
		if props := model.Properties; props != nil {
			if topLevelRegion := props.GeographicHierarchy; topLevelRegion != nil {
				result = topLevelRegion
				if !geographicalRegionIsMatch(topLevelRegion, name) {
					result = filterGeographicalRegions(topLevelRegion.Regions, name)
				}
			}
		}
	}

	if result == nil || result.Code == nil {
		return fmt.Errorf("Couldn't find a Traffic Manager Geographic Location with the name %q", name)
	}

	// NOTE: @tombuildsstuff: this is a unique data source that outputs the location as the ID, so this is fine
	id := *result.Code
	d.SetId(id)
	return nil
}

func filterGeographicalRegions(input *[]geographichierarchies.Region, name string) *geographichierarchies.Region {
	if regions := input; regions != nil {
		for _, region := range *regions {
			if geographicalRegionIsMatch(&region, name) {
				return &region
			}

			result := filterGeographicalRegions(region.Regions, name)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

func geographicalRegionIsMatch(input *geographichierarchies.Region, name string) bool {
	return strings.EqualFold(*input.Name, name)
}
