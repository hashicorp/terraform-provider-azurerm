package trafficmanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
		return fmt.Errorf("Error loading Traffic Manager Geographical Hierarchies: %+v", err)
	}

	name := d.Get("name").(string)

	var result *trafficmanager.Region
	if props := results.GeographicHierarchyProperties; props != nil {
		if topLevelRegion := props.GeographicHierarchy; topLevelRegion != nil {
			result = topLevelRegion
			if !geographicalRegionIsMatch(topLevelRegion, name) {
				result = filterGeographicalRegions(topLevelRegion.Regions, name)
			}
		}
	}

	if result == nil {
		return fmt.Errorf("Couldn't find a Traffic Manager Geographic Location with the name %q", name)
	}

	d.SetId(*result.Code)
	return nil
}

func filterGeographicalRegions(input *[]trafficmanager.Region, name string) *trafficmanager.Region {
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

func geographicalRegionIsMatch(input *trafficmanager.Region, name string) bool {
	return strings.EqualFold(*input.Name, name)
}
