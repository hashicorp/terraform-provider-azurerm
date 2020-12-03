package managedapplications_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplications/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMManagedApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedApplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedApplicationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMManagedApplication_requiresImport),
		},
	})
}

func TestAccAzureRMManagedApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedApplication_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedApplicationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagedApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagedApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagedApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "ServiceCatalog"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMManagedApplication_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MarketPlace"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMManagedApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagedApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "ServiceCatalog"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMManagedApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Managed Application not found: %s", resourceName)
		}

		id, err := parse.ApplicationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).ManagedApplication.ApplicationClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Managed Application %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on ManagedApplication.ApplicationClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMManagedApplicationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ManagedApplication.ApplicationClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_managed_application" {
			continue
		}

		id, err := parse.ApplicationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on ManagedApplication.ApplicationClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMManagedApplication_basic(data acceptance.TestData) string {
	template := testAccAzureRMManagedApplication_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameters = {
    location                 = azurerm_resource_group.test.location
    storageAccountNamePrefix = "store%s"
    storageAccountType       = "Standard_LRS"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMManagedApplication_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application" "import" {
  name                        = azurerm_managed_application.test.name
  location                    = azurerm_managed_application.test.location
  resource_group_name         = azurerm_managed_application.test.resource_group_name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%d"
}
`, testAccAzureRMManagedApplication_basic(data), data.RandomInteger)
}

func testAccAzureRMManagedApplication_complete(data acceptance.TestData) string {
	template := testAccAzureRMManagedApplication_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_agreement" "test" {
  publisher = "cisco"
  offer     = "meraki-vmx"
  plan      = "meraki-vmx100"
}

resource "azurerm_managed_application" "test" {
  name                        = "acctestCompleteManagedApp%d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "MarketPlace"
  managed_resource_group_name = "completeInfraGroup%d"

  plan {
    name      = azurerm_marketplace_agreement.test.plan
    product   = azurerm_marketplace_agreement.test.offer
    publisher = azurerm_marketplace_agreement.test.publisher
    version   = "1.0.44"
  }

  parameters = {
    baseUrl                     = ""
    location                    = azurerm_resource_group.test.location
    merakiAuthToken             = "f451adfb-d00b-4612-8799-b29294217d4a"
    subnetAddressPrefix         = "10.0.0.0/24"
    subnetName                  = "acctestSubnet"
    virtualMachineSize          = "Standard_DS12_v2"
    virtualNetworkAddressPrefix = "10.0.0.0/16"
    virtualNetworkName          = "acctestVnet"
    virtualNetworkNewOrExisting = "new"
    virtualNetworkResourceGroup = "acctestVnetRg"
    vmName                      = "acctestVM"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagedApplication_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Contributor"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mapp-%d"
  location = "%s"
}

resource "azurerm_managed_application_definition" "test" {
  name                = "acctestManagedAppDef%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lock_level          = "ReadOnly"
  package_file_uri    = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name        = "TestManagedAppDefinition"
  description         = "Test Managed App Definition"
  package_enabled     = true

  authorization {
    service_principal_id = data.azurerm_client_config.test.object_id
    role_definition_id   = split("/", data.azurerm_role_definition.test.id)[length(split("/", data.azurerm_role_definition.test.id)) - 1]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
