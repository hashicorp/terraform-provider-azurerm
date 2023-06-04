package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RoleAssignmentMarketplaceResource struct{}

func TestAccRoleAssignmentMarketplace_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	r := RoleAssignmentMarketplaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyNameConfig(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignmentMarketplace_roleName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.roleNameConfig(id),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("role_definition_name").HasValue("Marketplace Admin"),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignmentMarketplace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.roleNameConfig(id),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("role_definition_name").HasValue("Marketplace Admin"),
			),
		},
		{
			Config:      r.requiresImportConfig(id),
			ExpectError: acceptance.RequiresImportError("azurerm_role_assignment_marketplace"),
		},
	})
}

func TestAccRoleAssignmentMarketplace_builtin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.builtinConfig(id),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignmentMarketplace_ServicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipal(ri, id),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "principal_type", "ServicePrincipal"),
			),
		},
	})
}

func TestAccRoleAssignmentMarketplace_ServicePrincipalWithType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipalWithType(ri, id),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRoleAssignmentMarketplace_ServicePrincipalGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment_marketplace", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.group(ri, id),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r RoleAssignmentMarketplaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.RoleAssignmentMarketplaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Authorization.RoleAssignmentsClient.GetByID(ctx, state.ID, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Role Assignment for role %q: %+v", id.Name, err)
	}
	return utils.Bool(true), nil
}

func (RoleAssignmentMarketplaceResource) emptyNameConfig() string {
	return `
data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment_marketplace" "test" {
  role_definition_id = "${data.azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.object_id}"
  description        = "Test Role Assignment"
}
`
}

func (RoleAssignmentMarketplaceResource) roleNameConfig(id string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment_marketplace" "test" {
  name                 = "%s"
  role_definition_name = "Marketplace Admin"
  principal_id         = data.azurerm_client_config.test.object_id
}
`, id)
}

func (RoleAssignmentMarketplaceResource) requiresImportConfig(id string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment_marketplace" "import" {
  name                 = azurerm_role_assignment_marketplace.test.name
  role_definition_name = azurerm_role_assignment_marketplace.test.role_definition_name
  principal_id         = azurerm_role_assignment_marketplace.test.principal_id
}
`, RoleAssignmentMarketplaceResource{}.roleNameConfig(id))
}

func (RoleAssignmentMarketplaceResource) builtinConfig(id string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {
}

data "azurerm_role_definition" "test" {
  name = "Site Recovery Reader"
}

resource "azurerm_role_assignment_marketplace" "test" {
  name               = "%s"
  role_definition_id = "${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id
}
`, id)
}

func (RoleAssignmentMarketplaceResource) servicePrincipal(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_role_assignment_marketplace" "test" {
  name                 = "%s"
  role_definition_name = "Reader"
  principal_id         = azuread_service_principal.test.id
}
`, rInt, roleAssignmentID)
}

func (RoleAssignmentMarketplaceResource) servicePrincipalWithType(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_role_assignment_marketplace" "test" {
  name                             = "%s"
  role_definition_name             = "Reader"
  principal_id                     = azuread_service_principal.test.id
  skip_service_principal_aad_check = true
}
`, rInt, roleAssignmentID)
}

func (RoleAssignmentMarketplaceResource) group(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_group" "test" {
  display_name     = "acctestspa-%d"
  security_enabled = true
}

resource "azurerm_role_assignment_marketplace" "test" {
  name                 = "%s"
  role_definition_name = "Reader"
  principal_id         = azuread_group.test.id
}
`, rInt, roleAssignmentID)
}
