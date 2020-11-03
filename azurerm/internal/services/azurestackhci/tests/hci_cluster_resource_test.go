package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/azurestackhci/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMHCICluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHCIClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHCICluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHCIClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHCICluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHCIClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHCICluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHCIClusterExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMHCICluster_requiresImport),
		},
	})
}

func TestAccAzureRMHCICluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHCIClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHCICluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHCIClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMHCICluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hci_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHCIClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHCICluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHCIClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMHCICluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHCIClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMHCIClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).AzureStackHCI.ClusterClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("AzureStackHCI Cluster not found: %s", resourceName)
		}

		id, err := parse.HCIClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: AzureStackHCI Cluster %q does not exist", id.Name)
			}

			return fmt.Errorf("bad: Get on AzureStackHCI.ClusterClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMHCIClusterDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AzureStackHCI.ClusterClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hci_cluster" {
			continue
		}

		id, err := parse.HCIClusterID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on AzureStackHCI.ClusterClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMHCICluster_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-AzureStackHCI-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMHCICluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMHCICluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "test" {
  name                = "acctest-hci-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  aad_client_id       = data.azurerm_client_config.current.client_id
  aad_tenant_id       = data.azurerm_client_config.current.tenant_id
}
`, template, data.RandomInteger)
}

func testAccAzureRMHCICluster_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMHCICluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "import" {
  name                = azurerm_hci_cluster.test.name
  resource_group_name = azurerm_hci_cluster.test.resource_group_name
  location            = azurerm_hci_cluster.test.location
  aad_client_id       = azurerm_hci_cluster.test.aad_client_id
  aad_tenant_id       = azurerm_hci_cluster.test.aad_tenant_id
}
`, config)
}

func testAccAzureRMHCICluster_complete(data acceptance.TestData) string {
	template := testAccAzureRMHCICluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hci_cluster" "test" {
  name                = "acctest-hci-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  aad_client_id       = data.azurerm_client_config.current.client_id
  aad_tenant_id       = data.azurerm_client_config.current.tenant_id

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
