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

func TestAccAzureRMServiceEndpointPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceEndpointPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceEndpointPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceEndpointPolicy_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceEndpointPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceEndpointPolicy_requiresImport),
		},
	})
}

func testCheckAzureRMServiceEndpointPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Service Endpoint Policy not found: %s", resourceName)
		}

		id, err := parse.ServiceEndpointPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Service Endpoint Policy %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Network.ServiceEndpointPolicy: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMServiceEndpointPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_endpoint_policy" {
			continue
		}

		id, err := parse.ServiceEndpointPolicyID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err == nil {
			return fmt.Errorf("Network.ServiceEndpointPolicy still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Network.ServiceEndpointPolicy: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMServiceEndpointPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMServiceEndpointPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_service_endpoint_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceEndpointPolicy_complete(data acceptance.TestData) string {
	template := testAccAzureRMServiceEndpointPolicy_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_service_endpoint_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMServiceEndpointPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMServiceEndpointPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_service_endpoint_policy" "import" {
  name                = azurerm_service_endpoint_policy.test.name
  resource_group_name = azurerm_service_endpoint_policy.test.resource_group_name
  location            = azurerm_service_endpoint_policy.test.location
}
`, template)
}

func testAccAzureRMServiceEndpointPolicy_template(data acceptance.TestData) string {
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
