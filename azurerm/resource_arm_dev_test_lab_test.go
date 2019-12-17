package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDevTestLab_basic(t *testing.T) {
	resourceName := "azurerm_dev_test_lab.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestLab_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAzureRMDevTestLab_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dev_test_lab.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestLab_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDevTestLab_requiresImport(rInt, location),
				ExpectError: acceptance.RequiresImportError("azurerm_dev_test_lab"),
			},
		},
	})
}

func TestAccAzureRMDevTestLab_complete(t *testing.T) {
	resourceName := "azurerm_dev_test_lab.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestLab_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "World"),
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

func testCheckAzureRMDevTestLabExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		labName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DevTest Lab: %s", labName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, labName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestLabsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Lab %q (Resource Group: %q) does not exist", labName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDevTestLabDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.LabsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_lab" {
			continue
		}

		labName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, "")

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Lab still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDevTestLab_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMDevTestLab_requiresImport(rInt int, location string) string {
	template := testAccAzureRMDevTestLab_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_lab" "import" {
  name                = "${azurerm_dev_test_lab.test.name}"
  location            = "${azurerm_dev_test_lab.test.location}"
  resource_group_name = "${azurerm_dev_test_lab.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMDevTestLab_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  storage_type        = "Standard"

  tags = {
    Hello = "World"
  }
}
`, rInt, location, rInt)
}
