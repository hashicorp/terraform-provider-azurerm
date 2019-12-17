package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApiVersionSet_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementApiVersionSet_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_api_version_set"),
			},
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_header(t *testing.T) {
	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_header(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_query(t *testing.T) {
	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_query(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementApiVersionSet_update(t *testing.T) {
	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiVersionSet_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "TestDescription1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestApiVersionSet1%d", ri)),
				),
			},
			{
				Config: testAccAzureRMApiManagementApiVersionSet_update(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "TestDescription2"),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestApiVersionSet2%d", ri)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMApiManagementApiVersionSetDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiVersionSetClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_version_set" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)

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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ApiVersionSetClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Api Management Api Version Set %q (Resource Group %q / Api Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementApiVersionSetClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiVersionSet_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Segment"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementApiVersionSet_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiVersionSet_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "import" {
  name                = "${azurerm_api_management_api_version_set.test.name}"
  resource_group_name = "${azurerm_api_management_api_version_set.test.resource_group_name}"
  api_management_name = "${azurerm_api_management_api_version_set.test.api_management_name}"
  description         = "${azurerm_api_management_api_version_set.test.description}"
  display_name        = "${azurerm_api_management_api_version_set.test.display_name}"
  versioning_scheme   = "${azurerm_api_management_api_version_set.test.versioning_scheme}"
}
`, template)
}

func testAccAzureRMApiManagementApiVersionSet_header(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Header"
  version_header_name = "Header1"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementApiVersionSet_query(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Query"
  version_query_name  = "Query1"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementApiVersionSet_update(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiVersionSet_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  description         = "TestDescription2"
  display_name        = "TestApiVersionSet2%d"
  versioning_scheme   = "Segment"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementApiVersionSet_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}
`, rInt, location, rInt)
}
