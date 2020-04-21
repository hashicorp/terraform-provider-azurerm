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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRegistrationDefinition_basic(t *testing.T) {
	// Multiple tenants are needed to test this resource.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal in second tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_basic(uuid.New().String(), secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMRegistrationDefinition_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_basic(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
				),
			},
			{
				Config:      testAccAzureRMRegistrationDefinition_requiresImport(id, secondTenantID, principalID, data),
				ExpectError: acceptance.RequiresImportError("azurerm_registration_definition"),
			},
		},
	})
}

func TestAccAzureRMRegistrationDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_complete(uuid.New().String(), secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test Registration Definition"),
				),
			},
			data.ImportStep("registration_definition_id"),
		},
	})
}

func TestAccAzureRMRegistrationDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_basic(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
				),
			},
			{
				Config: testAccAzureRMRegistrationDefinition_complete(id, secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test Registration Definition"),
				),
			},
		},
	})
}

func TestAccAzureRMRegistrationDefinition_emptyID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_emptyId(secondTenantID, principalID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "registration_definition_id"),
				),
			},
		},
	})
}

func testCheckAzureRMRegistrationDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ManagedServices.RegistrationDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		scope := rs.Primary.Attributes["scope"]
		registrationDefinitionID := rs.Primary.Attributes["registration_definition_id"]

		resp, err := client.Get(ctx, scope, registrationDefinitionID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Registration Definition %q (Scope: %q) does not exist", registrationDefinitionID, scope)
			}
			return fmt.Errorf("Bad: Get on registrationDefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRegistrationDefinitionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ManagedServices.RegistrationDefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_registration_definition" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		registrationDefinitionID := rs.Primary.Attributes["registration_definition_id"]

		resp, err := client.Get(ctx, scope, registrationDefinitionID)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMRegistrationDefinition_basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id   = "%s"
  registration_definition_name = "acctestrd-%d"
  scope                        = data.azurerm_subscription.primary.id
  managed_by_tenant_id         = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func testAccAzureRMRegistrationDefinition_requiresImport(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_registration_definition" "import" {
  registration_definition_name = azurerm_registration_definition.test.registration_definition_name
  registration_definition_id   = azurerm_registration_definition.test.registration_definition_id
  scope                        = azurerm_registration_definition.test.scope
  managed_by_tenant_id         = azurerm_registration_definition.test.managed_by_tenant_id
  authorization {
    principal_id       = azurerm_registration_definition.test.managed_by_tenant_id
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
`, testAccAzureRMRegistrationDefinition_basic(id, secondTenantID, principalID, data))
}

func testAccAzureRMRegistrationDefinition_complete(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id   = "%s"
  registration_definition_name = "acctestrd-%d"
  description                  = "Acceptance Test Registration Definition"
  scope                        = data.azurerm_subscription.primary.id
  managed_by_tenant_id         = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
`, id, data.RandomInteger, secondTenantID, principalID)
}

func testAccAzureRMRegistrationDefinition_emptyId(secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "test" {
  registration_definition_name = "acctestrd-%d"
  description                  = "Acceptance Test Registration Definition"
  scope                        = data.azurerm_subscription.primary.id
  managed_by_tenant_id         = "%s"

  authorization {
    principal_id       = "%s"
    role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}
`, data.RandomInteger, secondTenantID, principalID)
}
