package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_basic(uuid.New().String(), data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRoleDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_basic(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMRoleDefinition_requiresImport(id, data),
				ExpectError: acceptance.RequiresImportError("azurerm_role_definition"),
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_complete(uuid.New().String(), data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep("role_definition_id", "scope"),
		},
	})
}

func TestAccAzureRMRoleDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_basic(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMRoleDefinition_updated(id, data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRoleDefinition_updateEmptyId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_emptyId(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMRoleDefinition_updateEmptyId(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRoleDefinition_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_emptyId(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRoleDefinition_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_managementGroup(uuid.New().String(), data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep("scope"),
		},
	})
}

func testCheckAzureRMRoleDefinitionExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMRoleDefinitionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_role_definition" {
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

func testAccAzureRMRoleDefinition_basic(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["*"]
    not_actions = []
  }
}
`, id, data.RandomInteger)
}

func testAccAzureRMRoleDefinition_requiresImport(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_definition" "import" {
  role_definition_id = azurerm_role_definition.test.role_definition_id
  name               = azurerm_role_definition.test.name
  scope              = azurerm_role_definition.test.scope

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}
`, testAccAzureRMRoleDefinition_basic(id, data))
}

func testAccAzureRMRoleDefinition_complete(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = data.azurerm_subscription.primary.id
  description        = "Acceptance Test Role Definition"

  permissions {
    actions          = ["*"]
    data_actions     = ["Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"]
    not_actions      = ["Microsoft.Authorization/*/read"]
    not_data_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}
`, id, data.RandomInteger)
}

func testAccAzureRMRoleDefinition_updated(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d-updated"
  scope              = data.azurerm_subscription.primary.id
  description        = "Acceptance Test Role Definition"

  permissions {
    actions     = ["*"]
    not_actions = ["Microsoft.Authorization/*/read"]
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}
`, id, data.RandomInteger)
}

func testAccAzureRMRoleDefinition_emptyId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "test" {
  name  = "acctestrd-%d"
  scope = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}
`, data.RandomInteger)
}

func testAccAzureRMRoleDefinition_updateEmptyId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "test" {
  name  = "acctestrd-%d"
  scope = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["*"]
    not_actions = ["Microsoft.Authorization/*/read"]
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}
`, data.RandomInteger)
}

func testAccAzureRMRoleDefinition_managementGroup(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_management_group" "test" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = azurerm_management_group.test.id

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    azurerm_management_group.test.id,
    data.azurerm_subscription.primary.id,
  ]
}
`, id, data.RandomInteger)
}
