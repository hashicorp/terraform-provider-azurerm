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

func TestAccAzureRMSubnetServiceEndpointPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetServiceEndpointPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetServiceEndpointPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetServiceEndpointPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSubnetServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSubnetServiceEndpointPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSubnetServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSubnetServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSubnetServiceEndpointPolicy_requiresImport),
		},
	})
}

func testCheckAzureRMSubnetServiceEndpointPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Subnet Service Endpoint Policy not found: %s", resourceName)
		}

		id, err := parse.SubnetServiceEndpointPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Subnet Service Endpoint Policy %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Subnet Service Endpoint Policy: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetServiceEndpointPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subnet_service_endpoint_policy" {
			continue
		}

		id, err := parse.SubnetServiceEndpointPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err == nil {
			return fmt.Errorf("Subnet Service Endpoint Policy still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Subnet Service Endpoint Policy: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMSubnetServiceEndpointPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMSubnetServiceEndpointPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMSubnetServiceEndpointPolicy_complete(data acceptance.TestData) string {
	template := testAccAzureRMSubnetServiceEndpointPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestasasepd%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_subnet_service_endpoint_policy" "test" {
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

func testAccAzureRMSubnetServiceEndpointPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSubnetServiceEndpointPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet_service_endpoint_policy" "import" {
  name                = azurerm_subnet_service_endpoint_policy.test.name
  resource_group_name = azurerm_subnet_service_endpoint_policy.test.resource_group_name
  location            = azurerm_subnet_service_endpoint_policy.test.location
}
`, template)
}

func testAccAzureRMSubnetServiceEndpointPolicy_template(data acceptance.TestData) string {
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
