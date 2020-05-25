package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataShareDataSetKustoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_kusto_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetKustoCluster_basic(data),
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

func TestAccAzureRMDataShareDataSetKustoCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_kusto_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataShareDataSetKustoCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataShareDataSetKustoCluster_requiresImport),
		},
	})
}

func testAccAzureRMDataShareDataSetKustoCluster_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-datashare-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-dsa-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind = "InPlace"
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

data "azuread_service_principal" "test" {
  display_name = azurerm_data_share_account.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_kusto_cluster.test.id
  role_definition_name = "Contributor"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func testAccAzureRMDataShareDataSetKustoCluster_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetKustoCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_kusto_cluster" "test" {
  name             = "acctest-dsbds-%d"
  share_id         = azurerm_data_share.test.id
  kusto_cluster_id = azurerm_kusto_cluster.test.id
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, config, data.RandomInteger)
}

func testAccAzureRMDataShareDataSetKustoCluster_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetKustoCluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_kusto_cluster" "import" {
  name             = azurerm_data_share_dataset_kusto_cluster.test.name
  share_id         = azurerm_data_share.test.id
  kusto_cluster_id = azurerm_kusto_cluster.test.id
}
`, config)
}
