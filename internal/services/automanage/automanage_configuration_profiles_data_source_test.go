package automanage_test

import (
	"testing"
)

type AutomanageConfigurationProfileDataSource struct{}

func TestAccAutomanageConfigurationProfileDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("tags").Exists(),
			),
		},
	})
}
func (AutomanageConfigurationProfileDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automanage_configuration_profile" "test" {
  name = azurerm_automanage_configuration_profile.test.name
  resource_group_name = azurerm_automanage_configuration_profile.test.resource_group_name
}
`, AutomanageConfigurationProfileResource{}.basic(data))
}
