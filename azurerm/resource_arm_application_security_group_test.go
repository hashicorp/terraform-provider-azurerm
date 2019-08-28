package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApplicationSecurityGroup_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_application_security_group.test"
	config := testAccAzureRMApplicationSecurityGroup_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationSecurityGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	ri := tf.AccRandTimeInt()
	location := testLocation()
	resourceName := "azurerm_application_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationSecurityGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApplicationSecurityGroup_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_app_service_custom_hostname_binding"),
			},
		},
	})
}

func TestAccAzureRMApplicationSecurityGroup_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_application_security_group.test"
	config := testAccAzureRMApplicationSecurityGroup_complete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationSecurityGroup_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()
	resourceName := "azurerm_application_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationSecurityGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMApplicationSecurityGroup_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationSecurityGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_security_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.ApplicationSecurityGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Application Security Group still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMApplicationSecurityGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Application Security Group: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).network.ApplicationSecurityGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Application Security Group %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on applicationSecurityGroupsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApplicationSecurityGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMApplicationSecurityGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApplicationSecurityGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_application_security_group" "import" {
  name                = "${azurerm_application_security_group.test.name}"
  location            = "${azurerm_application_security_group.test.location}"
  resource_group_name = "${azurerm_application_security_group.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMApplicationSecurityGroup_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    Hello = "World"
  }
}
`, rInt, location, rInt)
}
