package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RoleDefinitionResource struct{}

func TestAccRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(uuid.New().String(), data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}
	id := uuid.New().String()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(id, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(id, data)
		}),
	})
}

func TestAccRoleDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(uuid.New().String(), data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("role_definition_id", "scope"),
	})
}

func TestAccRoleDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}
	id := uuid.New().String()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(id, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdated(id, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleDefinition_updateEmptyId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateEmptyId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleDefinition_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRoleDefinition_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_definition", "test")
	r := RoleDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.TestAccRoleDefinition_managementGroup(uuid.New().String(), data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("scope"),
	})
}

func (r RoleDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	scope := state.Attributes["scope"]
	roleDefinitionId := state.Attributes["role_definition_id"]

	resp, err := client.Authorization.RoleDefinitionsClient.Get(ctx, scope, roleDefinitionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Role Definition %q (Scope %q): %+v", roleDefinitionId, scope, err)
	}
	return utils.Bool(true), nil
}

func (r RoleDefinitionResource) basic(id string, data acceptance.TestData) string {
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

func (r RoleDefinitionResource) requiresImport(id string, data acceptance.TestData) string {
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
`, r.basic(id, data))
}

func (r RoleDefinitionResource) complete(id string, data acceptance.TestData) string {
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

func (r RoleDefinitionResource) basicUpdated(id string, data acceptance.TestData) string {
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

func (r RoleDefinitionResource) emptyId(data acceptance.TestData) string {
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

func (r RoleDefinitionResource) updateEmptyId(data acceptance.TestData) string {
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

func (r RoleDefinitionResource) TestAccRoleDefinition_managementGroup(id string, data acceptance.TestData) string {
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
