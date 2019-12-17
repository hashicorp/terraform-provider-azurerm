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

func TestAccAzureRMAPIManagementGroup_basic(t *testing.T) {
	resourceName := "azurerm_api_management_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAPIManagementGroup_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(resourceName, "type", "custom"),
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

func TestAccAzureRMAPIManagementGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(resourceName, "type", "custom"),
				),
			},
			{
				Config:      testAccAzureRMAPIManagementGroup_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_group"),
			},
		},
	})
}

func TestAccAzureRMAPIManagementGroup_complete(t *testing.T) {
	resourceName := "azurerm_api_management_group.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAPIManagementGroup_complete(ri, acceptance.Location(), "Test Group", "A test description.")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Group"),
					resource.TestCheckResourceAttr(resourceName, "description", "A test description."),
					resource.TestCheckResourceAttr(resourceName, "type", "external"),
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

func TestAccAzureRMAPIManagementGroup_descriptionDisplayNameUpdate(t *testing.T) {
	resourceName := "azurerm_api_management_group.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMAPIManagementGroup_complete(ri, acceptance.Location(), "Original Group", "The original description.")
	postConfig := testAccAzureRMAPIManagementGroup_complete(ri, acceptance.Location(), "Modified Group", "A modified description.")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Original Group"),
					resource.TestCheckResourceAttr(resourceName, "description", "The original description."),
					resource.TestCheckResourceAttr(resourceName, "type", "external"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Modified Group"),
					resource.TestCheckResourceAttr(resourceName, "description", "A modified description."),
					resource.TestCheckResourceAttr(resourceName, "type", "external"),
				),
			},
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Original Group"),
					resource.TestCheckResourceAttr(resourceName, "description", "The original description."),
					resource.TestCheckResourceAttr(resourceName, "type", "external"),
				),
			},
		},
	})
}

func testCheckAzureRMAPIManagementGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.GroupClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_group" {
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

func testCheckAzureRMAPIManagementGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.GroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Group %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.GroupClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAPIManagementGroup_basic(rInt int, location string) string {
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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "Test Group"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAPIManagementGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAPIManagementGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_group" "import" {
  name                = "${azurerm_api_management_group.test.name}"
  resource_group_name = "${azurerm_api_management_group.test.resource_group_name}"
  api_management_name = "${azurerm_api_management_group.test.api_management_name}"
  display_name        = "${azurerm_api_management_group.test.display_name}"
}
`, template)
}

func testAccAzureRMAPIManagementGroup_complete(rInt int, location string, displayName, description string) string {
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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "%s"
  description         = "%s"
  type                = "external"
}
`, rInt, location, rInt, rInt, displayName, description)
}
