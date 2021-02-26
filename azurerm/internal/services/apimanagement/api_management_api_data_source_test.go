package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementApiDataSourceResource struct {
}

func TestAccDataSourceAzureRMApiManagementApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api", "test")
	r := ApiManagementApiDataSourceResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue("api1"),
				check.That(data.ResourceName).Key("path").HasValue("api1"),
				check.That(data.ResourceName).Key("protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("protocols.0").HasValue("https"),
				check.That(data.ResourceName).Key("soap_pass_through").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("true"),
				check.That(data.ResourceName).Key("is_current").HasValue("true"),
				check.That(data.ResourceName).Key("is_online").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceAzureRMApiManagementApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api", "test")
	r := ApiManagementApiDataSourceResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue("Butter Parser"),
				check.That(data.ResourceName).Key("path").HasValue("butter-parser"),
				check.That(data.ResourceName).Key("protocols.#").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("What is my purpose? You parse butter."),
				check.That(data.ResourceName).Key("service_url").HasValue("https://example.com/foo/bar"),
				check.That(data.ResourceName).Key("soap_pass_through").HasValue("false"),
				check.That(data.ResourceName).Key("subscription_key_parameter_names.0.header").HasValue("X-Butter-Robot-API-Key"),
				check.That(data.ResourceName).Key("subscription_key_parameter_names.0.query").HasValue("location"),
				check.That(data.ResourceName).Key("is_current").HasValue("true"),
				check.That(data.ResourceName).Key("is_online").HasValue("false"),
			),
		},
	})
}

func (r ApiManagementApiDataSourceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api" "test" {
  name                = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  revision            = azurerm_api_management_api.test.revision
}
`, ApiManagementApiResource{}.basic(data))
}

func (r ApiManagementApiDataSourceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api" "test" {
  name                = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  revision            = azurerm_api_management_api.test.revision
}
`, ApiManagementApiResource{}.complete(data))
}
