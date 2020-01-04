package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCognitiveAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCognitiveAccount_speechServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_speechServices(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "SpeechServices"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCognitiveAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMCognitiveAccount_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_cognitive_account"),
			},
		},
	})
}

func TestAccAzureRMCognitiveAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCognitiveAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppCognitiveAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCognitiveAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
				),
			},
			{
				Config: testAccAzureRMCognitiveAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCognitiveAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Face"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
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

func testAccAzureRMCognitiveAccount_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCognitiveAccount_speechServices(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCognitiveAccount_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMCognitiveAccount_basic(data)
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

func testAccAzureRMCognitiveAccount_complete(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
