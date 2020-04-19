package tests

import (
	"fmt"
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
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_basic(uuid.New().String(), data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
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
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_basic(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMRegistrationDefinition_requiresImport(id, data),
				ExpectError: acceptance.RequiresImportError("azurerm_registration_definition"),
			},
		},
	})
}

func TestAccAzureRMRegistrationDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_complete(uuid.New().String(), data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep("registration_definition_id"),
		},
	})
}

func TestAccAzureRMRegistrationDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_basic(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestMatchResourceAttr(data.ResourceName, "managed_by_tenant_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.role_definition_id", validate.UUIDRegExp),
				),
			},
			{
				Config: testAccAzureRMRegistrationDefinition_updated(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Acceptance Test Registration Definition"),
					resource.TestMatchResourceAttr(data.ResourceName, "managed_by_tenant_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "authorization.role_definition_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMRegistrationDefinition_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_registration_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRegistrationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRegistrationDefinition_emptyId(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testCheckAzureRMRegistrationDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
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

func testAccAzureRMRegistrationDefinition_basic(id string, data acceptance.TestData) string {
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
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}
`, data.RandomInteger, id, data.RandomInteger)
}

func testAccAzureRMRegistrationDefinition_requiresImport(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_registration_definition" "import" {
  registration_definition_id = azurerm_registration_definition.test.registration_definition_id
  name                       = azurerm_registration_definition.test.name
  scope                      = azurerm_registration_definition.test.scope
  managed_by_tenant_id       = azurerm_registration_definition.test.managed_by_tenant_id

  authorization {
	principal_id        = azurerm_registration_definition.test.authorization.principal_id
	role_definition_id  = azurerm_registration_definition.test.authorization.role_definition_id
  }
}
`, testAccAzureRMRegistrationDefinition_basic(id, data))
}

func testAccAzureRMRegistrationDefinition_complete(id string, data acceptance.TestData) string {
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
  description				 = "Acceptance Test Registration Definition"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}
`, data.RandomInteger, id, data.RandomInteger)
}

func testAccAzureRMRegistrationDefinition_updated(id string, data acceptance.TestData) string {
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
  description				 = "Acceptance Test Registration Definition"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = azuread_service_principal.test.id
	role_definition_id  = data.azurerm_role_definition.builtin.name
  }
}
`, data.RandomInteger, id, data.RandomInteger)
}

func testAccAzureRMRegistrationDefinition_emptyId(data acceptance.TestData) string {
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
	managed_by_tenant_id       	= "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"
		
	authorization {
		principal_id        = azuread_service_principal.test.id
		role_definition_id  = data.azurerm_role_definition.builtin.name
	}
}
`, data.RandomInteger, data.RandomInteger)
}
