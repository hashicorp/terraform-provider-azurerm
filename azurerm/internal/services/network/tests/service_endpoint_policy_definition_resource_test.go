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

func TestAccAzureRMServiceEndpointPolicyDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceEndpointPolicyDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicyDefinition_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceEndpointPolicyDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceEndpointPolicyDefinition_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMServiceEndpointPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMServiceEndpointPolicyDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_endpoint_policy_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceEndpointPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceEndpointPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceEndpointPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMServiceEndpointPolicyDefinition_requiresImport),
		},
	})
}

func testCheckAzureRMServiceEndpointPolicyDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPolicyDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Service Endpoint Policy Definition not found: %s", resourceName)
		}

		id, err := parse.ServiceEndpointPolicyDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Policy, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Service Endpoint Policy Definition %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Network.ServiceEndpointPolicyDefinition: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMServiceEndpointPolicyDefinitionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.ServiceEndpointPolicyDefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_endpoint_policy_definition" {
			continue
		}

		id, err := parse.ServiceEndpointPolicyDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Policy, id.Name)
		if err == nil {
			return fmt.Errorf("Network.ServiceEndpointPolicyDefinition still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Network.ServiceEndpointPolicyDefinition: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMServiceEndpointPolicyDefinition_basic(data acceptance.TestData) string {
	template := testAccAzureRMServiceEndpointPolicyDefinition_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_service_endpoint_policy_definition" "test" {
  name                  = "acctestSEPD-%d"
  policy_id             = azurerm_service_endpoint_policy.test.id
  service_endpoint_name = "Microsoft.Storage"
  service_resources = [
    "/subscriptions/%s",
  ]
}
`, template, data.RandomInteger, data.Client().SubscriptionID)
}

func testAccAzureRMServiceEndpointPolicyDefinition_complete(data acceptance.TestData) string {
	template := testAccAzureRMServiceEndpointPolicyDefinition_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestasasepd%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_service_endpoint_policy_definition" "test" {
  name                  = "acctestSEPD-%d"
  policy_id             = azurerm_service_endpoint_policy.test.id
  description           = "test policy def"
  service_endpoint_name = "Microsoft.Storage"
  service_resources = [
    "/subscriptions/%s",
    azurerm_resource_group.test.id,
    azurerm_storage_account.test.id
  ]
}
`, template, data.RandomString, data.RandomInteger, data.Client().SubscriptionID)
}

func testAccAzureRMServiceEndpointPolicyDefinition_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMServiceEndpointPolicyDefinition_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_service_endpoint_policy_definition" "import" {
  name                  = azurerm_service_endpoint_policy_definition.test.name
  policy_id             = azurerm_service_endpoint_policy_definition.test.policy_id
  service_endpoint_name = azurerm_service_endpoint_policy_definition.test.service_endpoint_name
  service_resources     = azurerm_service_endpoint_policy_definition.test.service_resources
}
`, template)
}

func testAccAzureRMServiceEndpointPolicyDefinition_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_endpoint_policy" "test" {
  name                = "acctestSEP-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
