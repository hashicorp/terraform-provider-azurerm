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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RoleAssignmentResource struct{}

func TestAccRoleAssignment_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyNameConfig(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignment_roleName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleNameConfig(id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("role_definition_name").HasValue("Log Analytics Reader"),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleNameConfig(id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("role_definition_name").HasValue("Log Analytics Reader"),
			),
		},
		{
			Config:      r.requiresImportConfig(id),
			ExpectError: acceptance.RequiresImportError("azurerm_role_assignment"),
		},
	})
}

func TestAccRoleAssignment_dataActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dataActionsConfig(id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttrSet(data.ResourceName, "role_definition_id"),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignment_builtin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.builtinConfig(id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignment_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	rInt := acceptance.RandTimeInt()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.customConfig(roleDefinitionId, roleAssignmentId, rInt),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignment_ServicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.servicePrincipal(ri, id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(data.ResourceName, "principal_type", "ServicePrincipal"),
			),
		},
	})
}

func TestAccRoleAssignment_ServicePrincipalWithType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.servicePrincipalWithType(ri, id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRoleAssignment_ServicePrincipalGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.group(ri, id),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

// TODO - "real" management group with appropriate required for testing
func TestAccRoleAssignment_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	groupId := uuid.New().String()

	r := RoleAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.managementGroupConfig(groupId),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r RoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.RoleAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Authorization.RoleAssignmentsClient.GetByID(ctx, state.ID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Role Assignment for role %q: %+v", id.Name, err)
	}
	return utils.Bool(true), nil
}

func (RoleAssignmentResource) emptyNameConfig() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "test" {
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.object_id}"
}
`
}

func (RoleAssignmentResource) roleNameConfig(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Log Analytics Reader"
  principal_id         = data.azurerm_client_config.test.object_id
}
`, id)
}

func (RoleAssignmentResource) requiresImportConfig(id string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "import" {
  name                 = azurerm_role_assignment.test.name
  scope                = azurerm_role_assignment.test.scope
  role_definition_name = azurerm_role_assignment.test.role_definition_name
  principal_id         = azurerm_role_assignment.test.principal_id
}
`, RoleAssignmentResource{}.roleNameConfig(id))
}

func (RoleAssignmentResource) dataActionsConfig(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Virtual Machine User Login"
  principal_id         = data.azurerm_client_config.test.object_id
}
`, id)
}

func (RoleAssignmentResource) builtinConfig(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

data "azurerm_role_definition" "test" {
  name = "Site Recovery Reader"
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id
}
`, id)
}

func (RoleAssignmentResource) customConfig(roleDefinitionId string, roleAssignmentId string, rInt int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = data.azurerm_subscription.primary.id
  description        = "Created by the Role Assignment Acceptance Test"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.test.role_definition_resource_id
  principal_id       = data.azurerm_client_config.test.object_id
}
`, roleDefinitionId, rInt, roleAssignmentId)
}

func (RoleAssignmentResource) servicePrincipal(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_subscription" "current" {
}

resource "azuread_application" "test" {
  name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Reader"
  principal_id         = azuread_service_principal.test.id
}
`, rInt, roleAssignmentID)
}

func (RoleAssignmentResource) servicePrincipalWithType(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_subscription" "current" {
}

resource "azuread_application" "test" {
  name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_role_assignment" "test" {
  name                             = "%s"
  scope                            = data.azurerm_subscription.current.id
  role_definition_name             = "Reader"
  principal_id                     = azuread_service_principal.test.id
  skip_service_principal_aad_check = true
}
`, rInt, roleAssignmentID)
}

func (RoleAssignmentResource) group(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_subscription" "current" {
}

resource "azuread_group" "test" {
  name = "acctestspa-%d"
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Reader"
  principal_id         = azuread_group.test.id
}
`, rInt, roleAssignmentID)
}

func (RoleAssignmentResource) managementGroupConfig(groupId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_management_group" "test" {
  group_id = "%s"
}

resource "azurerm_role_assignment" "test" {
  scope              = azurerm_management_group.test.id
  role_definition_id = data.azurerm_role_definition.test.id
  principal_id       = data.azurerm_client_config.test.object_id
}
`, groupId)
}
