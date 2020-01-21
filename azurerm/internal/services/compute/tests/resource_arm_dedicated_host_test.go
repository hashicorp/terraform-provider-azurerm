package tests

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

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
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", string(compute.DedicatedHostLicenseTypesNone)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_replace_on_failure", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "DSv3-Type1"),
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
					resource.TestCheckResourceAttr(data.ResourceName, "platform_fault_domain", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "DSv3-Type1"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", string(compute.DedicatedHostLicenseTypesWindowsServerHybrid)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_replace_on_failure", "false"),
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
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", string(compute.DedicatedHostLicenseTypesNone)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_replace_on_failure", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDedicatedHost_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", string(compute.DedicatedHostLicenseTypesWindowsServerHybrid)),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_replace_on_failure", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDedicatedHost_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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

		if resp, err := client.Get(ctx, id.ResourceGroup, id.HostGroup, id.Name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Dedicated Host %q (Host Group Name %q / Resource Group %q) does not exist", id.Name, id.HostGroup, id.ResourceGroup)
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

		if resp, err := client.Get(ctx, id.ResourceGroup, id.HostGroup, id.Name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.DedicatedHostsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDedicatedHost_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%s"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}


resource "azurerm_dedicated_host" "test" {
  name                  = "acctestDH-compute-%s"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  host_group_name       = azurerm_dedicated_host_group.test.name
  sku_name              = "DSv3-Type1"
  platform_fault_domain = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMDedicatedHost_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%s"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}


resource "azurerm_dedicated_host" "test" {
  name                    = "acctestDH-compute-%s"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  host_group_name         = azurerm_dedicated_host_group.test.name
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  license_type            = "Windows_Server_Hybrid"
  auto_replace_on_failure = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMDedicatedHost_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_dedicated_host" "import" {
  name                  = azurerm_dedicated_host.test.name
  resource_group_name   = azurerm_dedicated_host.test.resource_group_name
  location              = azurerm_dedicated_host.test.location
  host_group_name       = azurerm_dedicated_host.test.host_group_name
  sku_name              = azurerm_dedicated_host.test.sku_name
  platform_fault_domain = azurerm_dedicated_host.test.platform_fault_domain
}
`, testAccAzureRMDedicatedHost_basic(data))
}
