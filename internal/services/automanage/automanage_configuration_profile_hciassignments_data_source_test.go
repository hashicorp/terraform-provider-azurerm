package automanage_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"testing"
)

type AutomanageConfigurationProfileHCIAssignmentDataSource struct{}

func TestAccAutomanageConfigurationProfileHCIAssignmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automanage_configuration_profile_hciassignment", "test")
	r := AutomanageConfigurationProfileHCIAssignmentDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}
func (AutomanageConfigurationProfileHCIAssignmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automanage_configuration_profile_hciassignment" "test" {
  name = azurerm_automanage_configuration_profile_hciassignment.test.name
  resource_group_name = azurerm_automanage_configuration_profile_hciassignment.test.resource_group_name
  cluster_name = azurerm_automanage_configuration_profile_hciassignment.test.cluster_name
}
`, AutomanageConfigurationProfileHCIAssignmentResource{}.basic(data))
}
