package automanage_test

import (
	"testing"
)

type AutomanageConfigurationProfileHCRPAssignmentDataSource struct{}

func TestAccAutomanageConfigurationProfileHCRPAssignmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automanage_configuration_profile_hcrpassignment", "test")
	r := AutomanageConfigurationProfileHCRPAssignmentDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}
func (AutomanageConfigurationProfileHCRPAssignmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automanage_configuration_profile_hcrpassignment" "test" {
  name = azurerm_automanage_configuration_profile_hcrpassignment.test.name
  resource_group_name = azurerm_automanage_configuration_profile_hcrpassignment.test.resource_group_name
  machine_name = azurerm_automanage_configuration_profile_hcrpassignment.test.machine_name
}
`, AutomanageConfigurationProfileHCRPAssignmentResource{}.basic(data))
}
