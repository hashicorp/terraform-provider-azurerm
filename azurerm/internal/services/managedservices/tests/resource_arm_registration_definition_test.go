package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
			data.ImportStep("role_definition_id", "scope"),
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
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "0"),
				),
			},
			{
				Config: testAccAzureRMRegistrationDefinition_updated(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRegistrationDefinitionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/read"),
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
		roleDefinitionId := rs.Primary.Attributes["role_definition_id"]

		resp, err := client.Get(ctx, scope, roleDefinitionId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Role Definition %q (Scope: %q) does not exist", roleDefinitionId, scope)
			}
			return fmt.Errorf("Bad: Get on roleDefinitionsClient: %+v", err)
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
		roleDefinitionId := rs.Primary.Attributes["role_definition_id"]

		resp, err := client.Get(ctx, scope, roleDefinitionId)

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

resource "azurerm_registration_definition" "test" {
  registration_definition_id = "%s"
  name                       = "acctestrd-%d"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = "0ef18a36-5705-4d3f-99e8-a2a73de87d10"
	role_definition_id  = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
  }
}
`, id, data.RandomInteger)
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

resource "azurerm_registration_definition" "test" {
  registration_definition_id = "%s"
  name                       = "acctestrd-%d"
  description				 = "Acceptance Test Registration Definition"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = "0ef18a36-5705-4d3f-99e8-a2a73de87d10"
	role_definition_id  = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
  }
}
`, id, data.RandomInteger)
}

func testAccAzureRMRegistrationDefinition_updated(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "test" {
  registration_definition_id = "%s"
  name                       = "acctestrd-%d"
  description				 = "Acceptance Test Registration Definition"
  scope                      = data.azurerm_subscription.primary.id
  managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"

  authorization {
	principal_id        = "0ef18a36-5705-4d3f-99e8-a2a73de87d10"
	role_definition_id  = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
  }
}
`, id, data.RandomInteger)
}

func testAccAzureRMRegistrationDefinition_emptyId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "test" {
	name                       = "acctestrd-%d"
	description				 = "Acceptance Test Registration Definition"
	scope                      = data.azurerm_subscription.primary.id
	managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"
  
	authorization {
	  principal_id        = "0ef18a36-5705-4d3f-99e8-a2a73de87d10"
	  role_definition_id  = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
	}
}
`, data.RandomInteger)
}

func testAzureRMRegistrationDefinition_updateEmptyId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_registration_definition" "test" {
	name                       = "acctestrd-%d"
	description				 = "Acceptance Test Registration Definition"
	scope                      = data.azurerm_subscription.primary.id
	managed_by_tenant_id       = "70b0b1ee-ab4d-4c83-816a-2957afa9ea0b"
  
	authorization {
	  principal_id        = "0ef18a36-5705-4d3f-99e8-a2a73de87d10"
	  role_definition_id  = "43d0d8ad-25c7-4714-9337-8ba259a9fe05"
	}
}
`, data.RandomInteger)
}
