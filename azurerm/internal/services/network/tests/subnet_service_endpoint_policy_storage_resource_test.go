package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSubnetServiceEndpointStoragePolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointStoragePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointStoragePolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetServiceEndpointStoragePolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointStoragePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointStoragePolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetServiceEndpointStoragePolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointStoragePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointStoragePolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetServiceEndpointStoragePolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetServiceEndpointStoragePolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetServiceEndpointStoragePolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_storage_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointStoragePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointStoragePolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSubnetServiceEndpointStoragePolicy_requiresImport),
		},
	})
}

func testCheckAzureRMSubnetServiceEndpointStoragePolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Subnet Service Endpoint Storage Policy not found: %s", resourceName)
		}

		id, err := parse.SubnetServiceEndpointStoragePolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceEndpointPolicyName, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Subnet Service Endpoint Storage Policy %q (Resource Group %q) does not exist", id.ServiceEndpointPolicyName, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Subnet Service Endpoint Storage Policy: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetServiceEndpointStoragePolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subnet_service_endpoint_storage_policy" {
			continue
		}

		id, err := parse.SubnetServiceEndpointStoragePolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceEndpointPolicyName, "")
		if err == nil {
			return fmt.Errorf("Subnet Service Endpoint Storage Policy still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Subnet Service Endpoint Storage Policy: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMSubnetServiceEndpointStoragePolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMSubnetServiceEndpointStoragePolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMSubnetServiceEndpointStoragePolicy_complete(data acceptance.TestData) string {
	template := testAccAzureRMSubnetServiceEndpointStoragePolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestasasepd%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  definition {
    name        = "def1"
    description = "test definition1"
    service_resources = [
      "/subscriptions/%s",
      azurerm_resource_group.test.id,
      azurerm_storage_account.test.id
    ]
  }
  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, template, data.RandomString, data.RandomInteger, data.Client().SubscriptionID)
}

func testAccAzureRMSubnetServiceEndpointStoragePolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSubnetServiceEndpointStoragePolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_storage_policy" "import" {
  name                = azurerm_subnet_service_endpoint_storage_policy.test.name
  resource_group_name = azurerm_subnet_service_endpoint_storage_policy.test.resource_group_name
  location            = azurerm_subnet_service_endpoint_storage_policy.test.location
}
`, template)
}

func testAccAzureRMSubnetServiceEndpointStoragePolicy_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
