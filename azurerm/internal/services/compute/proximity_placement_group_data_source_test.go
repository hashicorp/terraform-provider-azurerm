package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ProximityPlacementGroupDataSource struct {
}

func TestAccProximityPlacementGroupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_proximity_placement_group", "test")
	r := ProximityPlacementGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
	})
}

func (ProximityPlacementGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_proximity_placement_group" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_proximity_placement_group.test.name
}
`, ProximityPlacementGroupResource{}.withTags(data))
}
