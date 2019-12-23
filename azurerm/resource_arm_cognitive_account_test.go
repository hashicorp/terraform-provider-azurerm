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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCognitiveAccount_basic(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCognitiveAccount_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
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

func TestAccAzureRMCognitiveAccount_speechServices(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCognitiveAccount_speechServices(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "SpeechServices"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
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

func TestAccAzureRMCognitiveAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_cognitive_account.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMCognitiveAccount_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_cognitive_account"),
			},
		},
	})
}

func TestAccAzureRMCognitiveAccount_complete(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCognitiveAccount_complete(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
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

func TestAccAzureRMCognitiveAccount_update(t *testing.T) {
	resourceName := "azurerm_cognitive_account.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
			{
				Config: testAccAzureRMCognitiveAccount_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
				),
			},
		},
	})
}

func testCheckAzureRMAppCognitiveAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cognitive.AccountsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cognitive_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetProperties(ctx, resourceGroup, name)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Cognitive Services Account still exists:\n%#v", resp)
			}

			return nil
		}
	}

	return nil
}

func testCheckAzureRMCognitiveAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cognitive.AccountsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.GetProperties(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Cognitive Services Account %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on cognitiveAccountsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMCognitiveAccount_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Face"

  sku {
    name = "S0"
    tier = "Standard"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCognitiveAccount_speechServices(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "SpeechServices"

  sku {
    name = "S0"
    tier = "Standard"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCognitiveAccount_requiresImport(rInt int, location string) string {
	template := testAccAzureRMCognitiveAccount_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_account" "import" {
  name                = "${azurerm_cognitive_account.test.name}"
  location            = "${azurerm_cognitive_account.test.location}"
  resource_group_name = "${azurerm_cognitive_account.test.resource_group_name}"
  kind                = "${azurerm_cognitive_account.test.kind}"

  sku {
    name = "S0"
    tier = "Standard"
  }
}
`, template)
}

func testAccAzureRMCognitiveAccount_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctestcogacc-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Face"

  sku {
    name = "S0"
    tier = "Standard"
  }

  tags = {
    Acceptance = "Test"
  }
}
`, rInt, location, rInt)
}
