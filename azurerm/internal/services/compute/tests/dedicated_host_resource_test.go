package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHost_autoReplaceOnFailure(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMDedicatedHost_autoReplaceOnFailure(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMDedicatedHost_autoReplaceOnFailure(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMDedicatedHost_autoReplaceOnFailure(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHost_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHost_licenceType(data, "None"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDedicatedHost_licenceType(data, "Windows_Server_Hybrid"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDedicatedHost_licenceType(data, "Windows_Server_Perpetual"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDedicatedHost_licenceType(data, "None"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHost_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHost_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHost_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDedicatedHost_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDedicatedHost_requiresImport),
		},
	})
}

func testCheckAzureRMDedicatedHostExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Dedicated Host not found: %s", resourceName)
		}

		id, err := parse.DedicatedHostID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DedicatedHostsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Dedicated Host %q (Host Group Name %q / Resource Group %q) does not exist", id.HostName, id.HostGroupName, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on Compute.DedicatedHostsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDedicatedHostDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.DedicatedHostsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dedicated_host" {
			continue
		}

		id, err := parse.DedicatedHostID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.DedicatedHostsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDedicatedHost_basic(data acceptance.TestData) string {
	template := testAccAzureRMDedicatedHost_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}
`, template, data.RandomInteger)
}

func testAccAzureRMDedicatedHost_autoReplaceOnFailure(data acceptance.TestData, replace bool) string {
	template := testAccAzureRMDedicatedHost_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  auto_replace_on_failure = %t
}
`, template, data.RandomInteger, replace)
}

func testAccAzureRMDedicatedHost_licenceType(data acceptance.TestData, licenseType string) string {
	template := testAccAzureRMDedicatedHost_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  license_type            = %q
}
`, template, data.RandomInteger, licenseType)
}

func testAccAzureRMDedicatedHost_complete(data acceptance.TestData) string {
	template := testAccAzureRMDedicatedHost_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  license_type            = "Windows_Server_Hybrid"
  auto_replace_on_failure = false
}
`, template, data.RandomInteger)
}

func testAccAzureRMDedicatedHost_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDedicatedHost_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_dedicated_host" "import" {
  name                    = azurerm_dedicated_host.test.name
  location                = azurerm_dedicated_host.test.location
  dedicated_host_group_id = azurerm_dedicated_host.test.dedicated_host_group_id
  sku_name                = azurerm_dedicated_host.test.sku_name
  platform_fault_domain   = azurerm_dedicated_host.test.platform_fault_domain
}
`, template)
}

func testAccAzureRMDedicatedHost_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctest-DHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
