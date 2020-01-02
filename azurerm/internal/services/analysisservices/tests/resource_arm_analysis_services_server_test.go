package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAnalysisServicesServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.label", "test"),
				),
			},
			{
				Config: testAccAzureRMAnalysisServicesServer_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "prod"),
				),
			},
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_querypoolConnectionMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(data, "All"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "querypool_connection_mode", "All"),
				),
			},
			{
				Config: testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(data, "ReadOnly"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "querypool_connection_mode", "ReadOnly"),
				),
			},
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_firewallSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	config1 := testAccAzureRMAnalysisServicesServer_firewallSettings1(data, true)

	config2 := testAccAzureRMAnalysisServicesServer_firewallSettings2(data, false)

	config3 := testAccAzureRMAnalysisServicesServer_firewallSettings3(data, true)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_power_bi_service", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.#", "0"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_power_bi_service", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.0.name", "test1"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.0.range_start", "92.123.234.11"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.0.range_end", "92.123.234.12"),
				),
			},
			{
				Config: config3,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_power_bi_service", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.0.name", "test1"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.0.range_start", "92.123.234.11"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.0.range_end", "92.123.234.13"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.1.name", "test2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.1.range_start", "226.202.187.57"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipv4_firewall_rule.1.range_end", "226.208.192.47"),
				),
			},
		},
	})
}

// ARM_ACC_EMAIL1 and ARM_ACC_EMAIL2 must be set and existing emails in the tenant's AD to work properly
func TestAccAzureRMAnalysisServicesServer_adminUsers(t *testing.T) {
	const ArmAccAdminEmail1 = "ARM_ACCTEST_ADMIN_EMAIL1"
	const ArmAccAdminEmail2 = "ARM_ACCTEST_ADMIN_EMAIL2"
	if os.Getenv(ArmAccAdminEmail1) == "" || os.Getenv(ArmAccAdminEmail2) == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env '%s' and '%s' set", ArmAccAdminEmail1, ArmAccAdminEmail2))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	email1 := os.Getenv(ArmAccAdminEmail1)
	email2 := os.Getenv(ArmAccAdminEmail2)
	preAdminUsers := []string{email1}
	postAdminUsers := []string{email1, email2}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_adminUsers(data, preAdminUsers),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAnalysisServicesServer_adminUsers(data, postAdminUsers),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_serverFullName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_serverFullName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_full_name"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_backupBlobContainerUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_backupBlobContainerUri(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "backup_blob_container_uri"),
				),
			},
		},
	})
}

func testAccAzureRMAnalysisServicesServer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "B1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAnalysisServicesServer_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "B1"

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAnalysisServicesServer_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "B1"

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(data acceptance.TestData, connectionMode string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                      = "acctestass%d"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  sku                       = "B1"
  querypool_connection_mode = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, connectionMode)
}

func testAccAzureRMAnalysisServicesServer_firewallSettings1(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                    = "acctestass%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  sku                     = "B1"
  enable_power_bi_service = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func testAccAzureRMAnalysisServicesServer_firewallSettings2(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                    = "acctestass%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  sku                     = "B1"
  enable_power_bi_service = %t

  ipv4_firewall_rule {
    name        = "test1"
    range_start = "92.123.234.11"
    range_end   = "92.123.234.12"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func testAccAzureRMAnalysisServicesServer_firewallSettings3(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                    = "acctestass%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  sku                     = "B1"
  enable_power_bi_service = %t

  ipv4_firewall_rule {
    name        = "test1"
    range_start = "92.123.234.11"
    range_end   = "92.123.234.13"
  }

  ipv4_firewall_rule {
    name        = "test2"
    range_start = "226.202.187.57"
    range_end   = "226.208.192.47"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func testAccAzureRMAnalysisServicesServer_adminUsers(data acceptance.TestData, adminUsers []string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "B1"
  admin_users         = ["%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, strings.Join(adminUsers, "\", \""))
}

func testAccAzureRMAnalysisServicesServer_serverFullName(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "B1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMAnalysisServicesServer_backupBlobContainerUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestass%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "assbackup"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  container_name    = "${azurerm_storage_container.test.name}"
  https_only        = true

  start  = "2018-06-01"
  expiry = "2048-06-01"

  permissions {
    read   = true
    add    = true
    create = true
    write  = true
    delete = true
    list   = true
  }
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "B1"

  backup_blob_container_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}${data.azurerm_storage_account_blob_container_sas.test.sas}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testCheckAzureRMAnalysisServicesServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AnalysisServices.ServerClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_analysis_services_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.GetDetails(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAnalysisServicesServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		analysisServicesServerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Analysis Services Server: %s", analysisServicesServerName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).AnalysisServices.ServerClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.GetDetails(ctx, resourceGroup, analysisServicesServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Analysis Services Server %q (resource group: %q) does not exist", analysisServicesServerName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on analysisServicesServerClient: %+v", err)
		}

		return nil
	}
}
