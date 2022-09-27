package automanage_test

import (
	"testing"
)

type AutomanageConfigurationProfilesVersionDataSource struct{}

func TestAccAutomanageConfigurationProfilesVersionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionDataSource{}
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
func (AutomanageConfigurationProfilesVersionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automanage_configuration_profiles_version" "test" {
  name = azurerm_automanage_configuration_profiles_version.test.name
  resource_group_name = azurerm_automanage_configuration_profiles_version.test.resource_group_name
  configuration_profile_name = azurerm_automanage_configuration_profiles_version.test.configuration_profile_name
}
`, AutomanageConfigurationProfilesVersionResource{}.basic(data))
}
