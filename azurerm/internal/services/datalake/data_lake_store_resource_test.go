package datalake_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMDataLakeStore_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Consumption"),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_type", "ServiceManaged"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataLakeStore_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDataLakeStore_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_data_lake_store"),
			},
		},
	})
}

func TestAccAzureRMDataLakeStore_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_tier(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Commitment_1TB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataLakeStore_encryptionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_encryptionDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_state", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "encryption_type", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataLakeStore_firewallUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_firewall(data, "Enabled", "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_allow_azure_ips", "Enabled"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_firewall(data, "Enabled", "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_state", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_allow_azure_ips", "Disabled"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_firewall(data, "Disabled", "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_state", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_allow_azure_ips", "Enabled"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_firewall(data, "Disabled", "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_state", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "firewall_allow_azure_ips", "Disabled"),
				),
			},
		},
	})
}

func TestAccAzureRMDataLakeStore_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStore_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStore_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataLakeStoreExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.StoreAccountsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		accountName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for data lake store: %s", accountName)
		}

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeStoreAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Store %q (resource group: %q) does not exist", accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeStoreDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.StoreAccountsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_store" {
			continue
		}

		accountName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, accountName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Data Lake Store still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeStore_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeStore_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDataLakeStore_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_store" "import" {
  name                = azurerm_data_lake_store.test.name
  resource_group_name = azurerm_data_lake_store.test.resource_group_name
  location            = azurerm_data_lake_store.test.location
}
`, template)
}

func testAccAzureRMDataLakeStore_tier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tier                = "Commitment_1TB"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeStore_encryptionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  encryption_state    = "Disabled"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeStore_firewall(data acceptance.TestData, firewallState string, firewallAllowAzureIPs string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                     = "acctest%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  firewall_state           = "%s"
  firewall_allow_azure_ips = "%s"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17], firewallState, firewallAllowAzureIPs)
}

func testAccAzureRMDataLakeStore_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeStore_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}
