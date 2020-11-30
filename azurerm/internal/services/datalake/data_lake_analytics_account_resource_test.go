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

func TestAccAzureRMDataLakeAnalyticsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Consumption"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataLakeAnalyticsAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataLakeAnalyticsAccount_requiresImport),
		},
	})
}

func TestAccAzureRMDataLakeAnalyticsAccount_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_tier(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "Commitment_100AUHours"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataLakeAnalyticsAccount_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDataLakeAnalyticsAccount_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataLakeAnalyticsAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.AnalyticsAccountsClient
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
			return fmt.Errorf("Bad: Get on dataLakeAnalyticsAccountClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Analytics Account %q (resource group: %q) does not exist", accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeAnalyticsAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_analytics_account" {
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

		return fmt.Errorf("Data Lake Analytics Account still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeAnalyticsAccount_basic(data acceptance.TestData) string {
	template := testAccAzureRMDataLakeStore_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name
}
`, template, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeAnalyticsAccount_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDataLakeAnalyticsAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "import" {
  name                       = azurerm_data_lake_analytics_account.test.name
  resource_group_name        = azurerm_data_lake_analytics_account.test.resource_group_name
  location                   = azurerm_data_lake_analytics_account.test.location
  default_store_account_name = azurerm_data_lake_analytics_account.test.default_store_account_name
}
`, template)
}

func testAccAzureRMDataLakeAnalyticsAccount_tier(data acceptance.TestData) string {
	template := testAccAzureRMDataLakeStore_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tier = "Commitment_100AUHours"

  default_store_account_name = azurerm_data_lake_store.test.name
}
`, template, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeAnalyticsAccount_withTags(data acceptance.TestData) string {
	template := testAccAzureRMDataLakeStore_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, template, strconv.Itoa(data.RandomInteger)[2:17])
}

func testAccAzureRMDataLakeAnalyticsAccount_withTagsUpdate(data acceptance.TestData) string {
	template := testAccAzureRMDataLakeStore_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name

  tags = {
    environment = "staging"
  }
}
`, template, strconv.Itoa(data.RandomInteger)[2:17])
}
