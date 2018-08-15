package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2017-05-01/trafficmanager"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmTrafficManagerGeographicalLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmTrafficManagerGeographicalLocationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceArmTrafficManagerGeographicalLocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).trafficManagerGeographialHierarchiesClient
	ctx := meta.(*ArmClient).StopContext

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
