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

func TestAccAzureRMRegistrationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_assignment", "test")
	secondTenantID := os.Getenv("ARM_SECOND_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationAssignment_basic(uuid.New().String(), secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationAssignmentExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMRegistrationAssignment_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_registration_assignment", "test")
	secondTenantID := os.Getenv("ARM_SECOND_TENANT_ID")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationAssignment_basic(id, secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationAssignmentExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMRegistrationAssignment_requiresImport(id, secondTenantID, data),
				ExpectError: acceptance.RequiresImportError("azurerm_registration_assignment"),
			},
		},
	})
}

func TestAccAzureRMRegistrationAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_assignment", "test")
	secondTenantID := os.Getenv("ARM_SECOND_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationAssignment_complete(uuid.New().String(), secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep("registration_assignment_id"),
		},
	})
}

func TestAccAzureRMRegistrationAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_assignment", "test")
	secondTenantID := os.Getenv("ARM_SECOND_TENANT_ID")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationAssignment_basic(id, secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationAssignmentExists(data.ResourceName),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
				),
			},
			{
				Config: testAccAzureRMRegistrationAssignment_updated(id, secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestMatchResourceAttr(data.ResourceName, "registration_definition_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMRegistrationAssignment_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_assignment", "test")
	secondTenantID := os.Getenv("ARM_SECOND_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationAssignment_emptyId(secondTenantID, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testCheckAzureRMRegistrationAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		scope := rs.Primary.Attributes["scope"]
		RegistrationAssignmentID := rs.Primary.Attributes["registration_assignment_id"]

		resp, err := client.Get(ctx, scope, RegistrationAssignmentID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Registration Definition %q (Scope: %q) does not exist", RegistrationAssignmentID, scope)
			}
			return fmt.Errorf("Bad: Get on RegistrationAssignmentsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRegistrationAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_registration_assignment" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		RegistrationAssignmentID := rs.Primary.Attributes["registration_assignment_id"]

		resp, err := client.Get(ctx, scope, RegistrationAssignmentID)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMRegistrationAssignment_basic(id string, secondTenantID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azuread_application" "test" {
	name = "acctestspa-%d"
}
		
resource "azuread_service_principal" "test" {
	application_id = azuread_application.test.application_id
}
	  
data "azurerm_role_definition" "builtin" {
	name = "Contributor"
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id = "%s"
  name                       = "acctestrd-%d"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "%s"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}

resource "azurerm_registration_assignment" "test" {
   registration_assignment_id = "%s"
   scope = data.azurerm_subscription.primary.id
   registration_definition_id = azurerm_registration_definition.test.id
}

`, data.RandomInteger, id, data.RandomInteger, secondTenantID, id)
}

func testAccAzureRMRegistrationAssignment_requiresImport(id string, secondTenantID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_registration_assignment" "import" {
  registration_assignment_id = azurerm_registration_assignment.test.registration_assignment_id
  registration_definition_id = azurerm_registration_assignment.test.registration_definition_id
  scope                      = azurerm_registration_assignment.test.scope
}
`, testAccAzureRMRegistrationAssignment_basic(id, secondTenantID, data))
}

func testAccAzureRMRegistrationAssignment_complete(id string, secondTenantID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azuread_application" "test" {
	name = "acctestspa-%d"
}
		
resource "azuread_service_principal" "test" {
	application_id = azuread_application.test.application_id
}
	  
data "azurerm_role_definition" "builtin" {
	name = "Contributor"
}

resource "azurerm_registration_definition" "test" {
  name                       = "acctestrd-%d"
  description				 = "Acceptance Test Registration Definition"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "%s"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}

resource "azurerm_registration_assignment" "test" {
	registration_assignment_id = "%s"
	scope = data.azurerm_subscription.primary.id
	registration_definition_id = azurerm_registration_definition.test.id
 }

`, data.RandomInteger, data.RandomInteger, secondTenantID, id)
}

func testAccAzureRMRegistrationAssignment_updated(id string, secondTenantID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}
	 
resource "azuread_application" "test" {
	name = "acctestspa-%d"
}
		
resource "azuread_service_principal" "test" {
	application_id = azuread_application.test.application_id
}
	  
data "azurerm_role_definition" "builtin" {
	name = "Contributor"
}

resource "azurerm_registration_definition" "test" {
  name                       = "acctestrd-%d"
  description				 = "Acceptance Test Registration Definition"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "%s"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}

resource "azurerm_registration_assignment" "test" {
	registration_assignment_id = "%s"
	scope = data.azurerm_subscription.primary.id
	registration_definition_id = azurerm_registration_definition.test.id
 }
`, data.RandomInteger, data.RandomInteger, secondTenantID, id)
}

func testAccAzureRMRegistrationAssignment_emptyId(secondTenantID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
	features {}
}
	  
data "azurerm_subscription" "primary" {
}
	 
resource "azuread_application" "test" {
	name = "acctestspa-%d"
}
		
resource "azuread_service_principal" "test" {
	application_id = azuread_application.test.application_id
}
	  
data "azurerm_role_definition" "builtin" {
	name = "Contributor"
}
	  
resource "azurerm_registration_definition" "test" {
	name                       	= "acctestrd-%d"
	description				 	= "Acceptance Test Registration Definition"
	scope                      	= data.azurerm_subscription.primary.id
	managed_by_tenant_id       	= "%s"
		
	authorization {
		principal_id        = azuread_service_principal.test.id
		role_definition_id  = data.azurerm_role_definition.builtin.name
	}
}

resource "azurerm_registration_assignment" "test" {
	scope = data.azurerm_subscription.primary.id
	registration_definition_id = azurerm_registration_definition.test.id
 }

`, data.RandomInteger, data.RandomInteger, secondTenantID)
}
