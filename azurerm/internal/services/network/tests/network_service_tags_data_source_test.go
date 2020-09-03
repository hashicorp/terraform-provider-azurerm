package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMServiceTags_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefixes.#"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMServiceTags_region(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_service_tags", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceTags_region(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefixes.#"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMServiceTags_basic() string {
	return `data "azurerm_network_service_tags" "test" {
  location = "westcentralus"
  service  = "AzureKeyVault"
}`
}

func testAccDataSourceAzureRMServiceTags_region() string {
	return `data "azurerm_network_service_tags" "test" {
  location        = "westcentralus"
  service         = "AzureKeyVault"
  location_filter = "australiacentral"
}`
}
