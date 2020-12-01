package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementApiVersionSet_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_header(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_header(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_query(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_query(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "TestDescription1"),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("TestApiVersionSet1%d", data.RandomInteger)),
				),
			},
			{
				Config: testAccAzureRMApiManagementApiVersionSet_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "TestDescription2"),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("TestApiVersionSet2%d", data.RandomInteger)),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApiManagementApiVersionSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiVersionSetClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_version_set" {
			continue
		}

		id, err := parse.ApiVersionSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMApiManagementApiVersionSetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiVersionSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ApiVersionSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Api Management Api Version Set %q (Resource Group %q / Api Management Service %q) does not exist", id.Name, id.ResourceGroup, id.ServiceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementApiVersionSetClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiVersionSet_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Segment"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApiVersionSet_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiVersionSet_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "import" {
  name                = azurerm_api_management_api_version_set.test.name
  resource_group_name = azurerm_api_management_api_version_set.test.resource_group_name
  api_management_name = azurerm_api_management_api_version_set.test.api_management_name
  description         = azurerm_api_management_api_version_set.test.description
  display_name        = azurerm_api_management_api_version_set.test.display_name
  versioning_scheme   = azurerm_api_management_api_version_set.test.versioning_scheme
}
`, template)
}

func testAccAzureRMApiManagementApiVersionSet_header(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Header"
  version_header_name = "Header1"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApiVersionSet_query(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Query"
  version_query_name  = "Query1"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApiVersionSet_update(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription2"
  display_name        = "TestApiVersionSet2%d"
  versioning_scheme   = "Segment"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementApiVersionSet_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
