package eventhub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccEventHubNamespaceAuthorizationRuleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEventHubNamespaceDataSourceAuthorizationRule_basic(data, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "listen"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "manage"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "send"),
				),
			},
		},
	})
}

func TestAccEventHubNamespaceAuthorizationRuleDataSource_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventhub_namespace_authorization_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				// `primary_connection_string_alias` and `secondary_connection_string_alias` are still `nil` while `data.azurerm_eventhub_namespace_authorization_rule` is retrieving resource. since `azurerm_eventhub_namespace_disaster_recovery_config` hasn't been created.
				// So these two properties should be checked in the second run.
				// And `depends_on` cannot be applied to `azurerm_eventhub_namespace_authorization_rule`.
				// Because it would throw error message `BreakPairing operation is only allowed on primary namespace with valid secondary namespace.` while destroying `azurerm_eventhub_namespace_disaster_recovery_config` if `depends_on` is applied.
				Config: testAccEventHubNamespaceAuthorizationRule_withAliasConnectionString(data),
			},
			{
				Config: testAccEventHubNamespaceDataSourceAuthorizationRule_withAliasConnectionString(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string_alias"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string_alias"),
				),
			},
		},
	})
}

func testAccEventHubNamespaceDataSourceAuthorizationRule_basic(data acceptance.TestData, listen, send, manage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "acctest-EHN-AR%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}

data "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = azurerm_eventhub_namespace_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, listen, send, manage)
}

func testAccEventHubNamespaceDataSourceAuthorizationRule_withAliasConnectionString(data acceptance.TestData) string {
	template := testAccEventHubNamespaceAuthorizationRule_withAliasConnectionString(data)
	return fmt.Sprintf(`
%s

data "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = azurerm_eventhub_namespace_authorization_rule.test.name
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
