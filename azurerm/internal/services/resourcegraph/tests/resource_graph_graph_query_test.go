package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resourcegraph/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMResourceGraphGraphQuery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_graph_query", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGraphGraphQueryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGraphGraphQuery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGraphGraphQueryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceGraphGraphQuery_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_graph_query", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGraphGraphQueryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGraphGraphQuery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGraphGraphQueryExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMResourceGraphGraphQuery_requiresImport),
		},
	})
}

func TestAccAzureRMResourceGraphGraphQuery_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_graph_query", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGraphGraphQueryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGraphGraphQuery_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGraphGraphQueryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMResourceGraphGraphQuery_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_graph_graph_query", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMResourceGraphGraphQueryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMResourceGraphGraphQuery_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGraphGraphQueryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMResourceGraphGraphQuery_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGraphGraphQueryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Docker VMs in PROD"),
					resource.TestCheckResourceAttr(data.ResourceName, "query", "where isnotnull(tags['Prod']) and properties.extensions[0].Name == 'docker'"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
		},
	})
}

func testCheckAzureRMResourceGraphGraphQueryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ResourceGraph.GraphQueryClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resourceGraph GraphQuery not found: %s", resourceName)
		}
		id, err := parse.ResourceGraphGraphQueryID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.ResourceName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: ResourceGraph GraphQuery %q does not exist", resourceName)
			}
			return fmt.Errorf("bad: Get on ResourceGraph.GraphQueryClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMResourceGraphGraphQueryDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ResourceGraph.GraphQueryClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_resource_graph_graph_query" {
			continue
		}
		id, err := parse.ResourceGraphGraphQueryID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.ResourceName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on ResourceGraph.GraphQueryClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMResourceGraphGraphQuery_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMResourceGraphGraphQuery_basic(data acceptance.TestData) string {
	template := testAccAzureRMResourceGraphGraphQuery_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_graph_graph_query" "test" {
  resource_group_name = azurerm_resource_group.test.name
  resource_name = "acctest-rggq-%d"
  query = "where isnotnull(tags['Prod']) and properties.extensions[0].Name == 'docker'"
}
`, template, data.RandomInteger)
}

func testAccAzureRMResourceGraphGraphQuery_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMResourceGraphGraphQuery_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_graph_graph_query" "import" {
  resource_group_name = azurerm_resource_graph_graph_query.test.resource_group_name
  query = azurerm_resource_graph_graph_query.test.query
  resource_name = azurerm_resource_graph_graph_query.test.resource_name
}
`, config)
}

func testAccAzureRMResourceGraphGraphQuery_complete(data acceptance.TestData) string {
	template := testAccAzureRMResourceGraphGraphQuery_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_graph_graph_query" "test" {
  resource_group_name = azurerm_resource_group.test.name
  resource_name = "acctest-rggq-%d"
  query = "where isnotnull(tags['Prod']) and properties.extensions[0].Name == 'docker'"
  description = "Docker VMs in PROD"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
