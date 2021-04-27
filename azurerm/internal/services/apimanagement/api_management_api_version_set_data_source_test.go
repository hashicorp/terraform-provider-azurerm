package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementApiVersionSetDataSource struct {
}

func TestAccDataSourceApiManagementApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("api_management_name").Exists(),
			),
		},
	})
}

func (ApiManagementApiVersionSetDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api_version_set" "test" {
  name                = azurerm_api_management_api_version_set.test.name
  resource_group_name = azurerm_api_management_api_version_set.test.resource_group_name
  api_management_name = azurerm_api_management_api_version_set.test.api_management_name
}
`, ApiManagementApiVersionSetResource{}.basic(data))
}
