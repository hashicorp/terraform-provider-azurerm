package lighthouse_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLighthouseAssignment_basic(t *testing.T) {
	// Multiple tenants are needed to test this resource.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal from second Tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseAssignment_basic(uuid.New().String(), secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func TestAccAzureRMLighthouseAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_assignment", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseAssignment_basic(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
			{
				Config:      testAccAzureRMLighthouseAssignment_requiresImport(id, secondTenantID, principalID, data),
				ExpectError: acceptance.RequiresImportError("azurerm_lighthouse_assignment"),
			},
		},
	})
}

func TestAccAzureRMLighthouseAssignment_emptyID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lighthouse_assignment", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLighthouseAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLighthouseAssignment_emptyId(secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLighthouseAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testCheckAzureRMLighthouseAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Lighthouse.AssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		scope := rs.Primary.Attributes["scope"]
		lighthouseAssignmentName := rs.Primary.Attributes["name"]
		expandLighthouseDefinition := true

		resp, err := client.Get(ctx, scope, lighthouseAssignmentName, &expandLighthouseDefinition)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Lighthouse Assignment %q (Scope: %q) does not exist", lighthouseAssignmentName, scope)
			}
			return fmt.Errorf("Bad: Get on LighthouseAssignmentsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMLighthouseAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Lighthouse.AssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_lighthouse_assignment" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		lighthouseAssignmentName := rs.Primary.Attributes["name"]
		expandLighthouseDefinition := true

		resp, err := client.Get(ctx, scope, lighthouseAssignmentName, &expandLighthouseDefinition)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMLighthouseAssignment_basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
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

resource "azurerm_lighthouse_assignment" "test" {
  name                     = "%s"
  scope                    = data.azurerm_subscription.primary.id
  lighthouse_definition_id = azurerm_lighthouse_definition.test.id
}

`, data.RandomInteger, secondTenantID, principalID, id)
}

func testAccAzureRMLighthouseAssignment_requiresImport(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lighthouse_assignment" "import" {
  name                     = azurerm_lighthouse_assignment.test.name
  lighthouse_definition_id = azurerm_lighthouse_assignment.test.lighthouse_definition_id
  scope                    = azurerm_lighthouse_assignment.test.scope
}
`, testAccAzureRMLighthouseAssignment_basic(id, secondTenantID, principalID, data))
}

func testAccAzureRMLighthouseAssignment_emptyId(secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
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

resource "azurerm_lighthouse_assignment" "test" {
  scope                    = data.azurerm_subscription.primary.id
  lighthouse_definition_id = azurerm_lighthouse_definition.test.id
}
`, data.RandomInteger, secondTenantID, principalID)
}
