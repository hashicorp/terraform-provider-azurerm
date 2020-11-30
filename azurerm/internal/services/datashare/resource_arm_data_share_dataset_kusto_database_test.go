package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataShareDataSetKustoDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_database", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_kusto_database"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetKustoDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kusto_cluster_location"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataShareDataSetKustoDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_database", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_kusto_database"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetKustoDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataShareDataSetKustoDatabase_requiresImport),
		},
	})
}

func testAccAzureRMDataShareDataSetKustoDatabase_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datashare-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_DS_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "InPlace"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestKD-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_kusto_cluster.test.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_data_share_account.test.identity.0.principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func testAccAzureRMDataShareDataSetKustoDatabase_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetKustoDatabase_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_kusto_database" "test" {
  name              = "acctest-DSKD-%d"
  share_id          = azurerm_data_share.test.id
  kusto_database_id = azurerm_kusto_database.test.id
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccAzureRMDataShareDataSetKustoDatabase_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetKustoDatabase_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_kusto_database" "import" {
  name              = azurerm_data_share_dataset_kusto_database.test.name
  share_id          = azurerm_data_share.test.id
  kusto_database_id = azurerm_kusto_database.test.id
}
`, config)
}
