package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type HealthCareServiceDataSource struct {
}

func TestAccHealthCareServiceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_service", "test")
	r := HealthCareServiceDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("kind").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_throughput").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_key_vault_key_versionless_id").HasValue(""),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccHealthCareServiceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_service", "test")
	r := HealthCareServiceDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cosmosdb_key_vault_key_versionless_id").Exists(),
			),
		},
	})
}

func (HealthCareServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_service" "test" {
  name                = azurerm_healthcare_service.test.name
  resource_group_name = azurerm_healthcare_service.test.resource_group_name
  location            = azurerm_resource_group.test.location
}
`, HealthCareServiceResource{}.basic(data))
}

func (HealthCareServiceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_service" "test" {
  name                = azurerm_healthcare_service.test.name
  resource_group_name = azurerm_healthcare_service.test.resource_group_name
  location            = azurerm_resource_group.test.location
}
`, HealthCareServiceResource{}.complete(data))
}
