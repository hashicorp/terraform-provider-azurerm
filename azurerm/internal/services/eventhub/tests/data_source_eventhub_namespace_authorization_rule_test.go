package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMEventHubNamespaceAuthorizationRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.eventhub_namespace_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHubNamespaceAuthorizationRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccDataSourceEventHubNamespaceAuthorizationRule_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = "Standard"
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctest-EHN-AR%[1]d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.resource_group_name}"
}
`, data.RandomInteger, data.Locations.Primary)
}
