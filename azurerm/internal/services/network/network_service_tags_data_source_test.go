package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NetworkServiceTagsDataSource struct {
}

func TestAccDataSourceAzureRMServiceTags_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_region(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")
	r := NetworkServiceTagsDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.region(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
			),
		},
	})
}

func (NetworkServiceTagsDataSource) basic() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "westcentralus"
  service  = "AzureKeyVault"
}`
}

func (NetworkServiceTagsDataSource) region() string {
	return `data "azurerm_network_service_tags" "test" {
  location        = "westcentralus"
  service         = "AzureKeyVault"
  location_filter = "australiacentral"
}`
}
