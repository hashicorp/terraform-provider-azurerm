package automanage_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"testing"
)

type AutomanageConfigurationProfileAssignmentDataSource struct{}

func TestAccAutomanageConfigurationProfileAssignmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automanage_configuration_profile_assignment", "test")
	r := AutomanageConfigurationProfileAssignmentDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}
func (AutomanageConfigurationProfileAssignmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automanage_configuration_profile_assignment" "test" {
  name = azurerm_automanage_configuration_profile_assignment.test.name
  resource_group_name = azurerm_automanage_configuration_profile_assignment.test.resource_group_name
  vm_name = azurerm_automanage_configuration_profile_assignment.test.vm_name
}
`, AutomanageConfigurationProfileAssignmentResource{}.basic(data))
}
