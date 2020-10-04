package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLighthouseDefinition_basic(t *testing.T) {
	// Multiple tenants are needed to test this resource.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal from second Tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseDefinition_basic(uuid.New().String(), secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMLighthouseDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseDefinition_basic(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				),
			},
			{
				Config:      testAccAzureRMLighthouseDefinition_requiresImport(id, secondTenantID, principalID, data),
				ExpectError: acceptance.RequiresImportError("azurerm_lighthouse_definition"),
			},
		},
	})
}

func TestAccAzureRMLighthouseDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseDefinition_complete(uuid.New().String(), secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test Lighthouse Definition"),
				),
			},
			data.ImportStep("lighthouse_definition_id"),
		},
	})
}

func TestAccAzureRMLighthouseDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseDefinition_basic(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
				),
			},
			{
				Config: testAccAzureRMLighthouseDefinition_complete(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "lighthouse_definition_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test Lighthouse Definition"),
				),
			},
		},
	})
}

func TestAccAzureRMLighthouseDefinition_emptyID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseDefinition_emptyId(secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "lighthouse_definition_id"),
				),
			},
		},
	})
}

func testCheckAzureRMLighthouseDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Lighthouse.DefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		scope := rs.Primary.Attributes["scope"]
		lighthouseDefinitionID := rs.Primary.Attributes["lighthouse_definition_id"]

		resp, err := client.Get(ctx, scope, lighthouseDefinitionID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Lighthouse Definition %q (Scope: %q) does not exist", lighthouseDefinitionID, scope)
			}
			return fmt.Errorf("Bad: Get on lighthouseDefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLighthouseDefinitionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Lighthouse.DefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_lighthouse_definition" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		lighthouseDefinitionID := rs.Primary.Attributes["lighthouse_definition_id"]

		resp, err := client.Get(ctx, scope, lighthouseDefinitionID)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMLighthouseDefinition_basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  managing_tenant_id       = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func testAccAzureRMLighthouseDefinition_requiresImport(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lighthouse_definition" "import" {
  name                     = azurerm_lighthouse_definition.test.name
  lighthouse_definition_id = azurerm_lighthouse_definition.test.lighthouse_definition_id
  managing_tenant_id       = azurerm_lighthouse_definition.test.managing_tenant_id
  authorization {
    principal_id       = azurerm_lighthouse_definition.test.managing_tenant_id
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
`, testAccAzureRMLighthouseDefinition_basic(id, secondTenantID, principalID, data))
}

func testAccAzureRMLighthouseDefinition_complete(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "test" {
  lighthouse_definition_id = "%s"
  name                     = "acctest-LD-%d"
  description              = "Acceptance Test Lighthouse Definition"
  managing_tenant_id       = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func testAccAzureRMLighthouseDefinition_emptyId(secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "test" {
  name               = "acctest-LD-%d"
  description        = "Acceptance Test Lighthouse Definition"
  managing_tenant_id = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}
`, data.RandomInteger, secondTenantID, principalID)
}
