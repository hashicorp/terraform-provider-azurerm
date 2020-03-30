package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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
				),
			},
			data.ImportStep(),
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

func testCheckAzureRMManagedApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Managed Application not found: %s", resourceName)
		}

		id, err := parse.ManagedApplicationID(rs.Primary.ID)
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

		id, err := parse.ManagedApplicationID(rs.Primary.ID)
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
  name                      = "acctestManagedApp%d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  kind                      = "ServiceCatalog"
  managed_resource_group_id = "/subscriptions/${data.azurerm_client_config.test.subscription_id}/resourceGroups/infraGroup%d"
  application_definition_id = azurerm_managed_application_definition.test.id

  parameters = <<PARAMETERS
    {
        "location": {
            "value": "${azurerm_resource_group.test.location}"
        },
        "storageAccountNamePrefix": {
            "value": "store%s"
        },
        "storageAccountType": {
            "value": "Standard_LRS"
        }
    }
    PARAMETERS
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func testAccAzureRMManagedApplication_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application" "import" {
  name                = azurerm_managed_application.test.name
  location            = azurerm_managed_application.test.location
  resource_group_name = azurerm_managed_application.test.resource_group_name
}
`, testAccAzureRMManagedApplication_basic(data))
}

func testAccAzureRMManagedApplication_complete(data acceptance.TestData) string {
	template := testAccAzureRMManagedApplication_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application" "test" {
  name                      = "acctestCompleteManagedApp%d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  kind                      = "MarketPlace"
  managed_resource_group_id = "/subscriptions/${data.azurerm_client_config.test.subscription_id}/resourceGroups/completeInfraGroup%d"

  plan {
    name      = "meraki-vmx100"
    product   = "meraki-vmx"
    publisher = "cisco"
    version   = "1.0.44"
  }

  parameters = <<PARAMETERS
    {
        "baseUrl": {
            "value": ""
        },
        "location": {
            "value": "${azurerm_resource_group.test.location}"
        },
        "merakiAuthToken": {
            "value": "f451adfb-d00b-4612-8799-b29294217d4a"
        },
        "subnetAddressPrefix": {
            "value": "10.0.0.0/24"
        },
        "subnetName": {
            "value": "acctestSubnet"
        },
        "virtualMachineSize": {
            "value": "Standard_DS12_v2"
        },
        "virtualNetworkAddressPrefix": {
            "value": "10.0.0.0/16"
        },
        "virtualNetworkName": {
            "value": "acctestVnet"
        },
        "virtualNetworkNewOrExisting": {
            "value": "new"
        },
        "virtualNetworkResourceGroup": {
            "value": "acctestVnetRg"
        },
        "vmName": {
            "value": "acctestVM"
        }
    }
    PARAMETERS

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
`, data.RandomInteger, "westus2", data.RandomInteger)
}
