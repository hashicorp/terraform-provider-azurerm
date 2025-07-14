// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RoleAssignmentMarketplaceResource struct{}

func TestAccRoleAssignmentMarketplace_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	r := RoleAssignmentMarketplaceResource{}
	roleName := "Reader"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.emptyNameConfig(roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignmentMarketplace_roleName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}
	roleName := "Log Analytics Reader"

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.roleNameConfig(id, roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("role_definition_name").HasValue(roleName),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignmentMarketplace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}
	roleName := "Managed Applications Reader"

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.roleNameConfig(id, roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_id").Exists(),
				check.That(data.ResourceName).Key("role_definition_name").HasValue(roleName),
			),
		},
		{
			Config:      r.requiresImportConfig(id, roleName),
			ExpectError: acceptance.RequiresImportError("azurerm_marketplace_role_assignment"),
		},
	})
}

func TestAccRoleAssignmentMarketplace_builtin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}
	roleName := "Monitoring Reader"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.builtinConfig(id, roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("skip_service_principal_aad_check"),
	})
}

func TestAccRoleAssignmentMarketplace_ServicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}
	roleName := "Contributor"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.servicePrincipal(ri, id, roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "principal_type", "ServicePrincipal"),
			),
		},
	})
}

func TestAccRoleAssignmentMarketplace_ServicePrincipalWithType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}
	roleName := "Log Analytics Contributor"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.servicePrincipalWithType(ri, id, roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccRoleAssignmentMarketplace_ServicePrincipalGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_marketplace_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	r := RoleAssignmentMarketplaceResource{}
	roleName := "Monitoring Contributor"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Last error may cause the role already assigned. Need to delete it before a new test.
			PreConfig: r.deleteAssignedRole(t, roleName),
			Config:    r.group(ri, id, roleName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r RoleAssignmentMarketplaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ScopedRoleAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	options := roleassignments.DefaultGetByIdOperationOptions()
	if id.TenantId != "" {
		options.TenantId = &id.TenantId
	}

	resp, err := client.Authorization.ScopedRoleAssignmentsClient.GetById(ctx, commonids.NewScopeID(id.ID()), options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (RoleAssignmentMarketplaceResource) emptyNameConfig(roleName string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "%s"
}

provider "azurerm" {
  features {}
}

resource "azurerm_marketplace_role_assignment" "test" {
  role_definition_id = "${data.azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.object_id}"
  description        = "Test Role Assignment"

  lifecycle {
    ignore_changes = [
      name,
      role_definition_name,
    ]
  }
}
`, roleName)
}

func (RoleAssignmentMarketplaceResource) roleNameConfig(id string, roleName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {
}

resource "azurerm_marketplace_role_assignment" "test" {
  name                 = "%s"
  role_definition_name = "%s"
  principal_id         = data.azurerm_client_config.test.object_id

  lifecycle {
    ignore_changes = [
      role_definition_id,
    ]
  }
}
`, id, roleName)
}

func (RoleAssignmentMarketplaceResource) requiresImportConfig(id string, roleName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_marketplace_role_assignment" "import" {
  name                 = azurerm_marketplace_role_assignment.test.name
  role_definition_name = azurerm_marketplace_role_assignment.test.role_definition_name
  principal_id         = azurerm_marketplace_role_assignment.test.principal_id

  lifecycle {
    ignore_changes = [
      role_definition_id,
    ]
  }
}
`, RoleAssignmentMarketplaceResource{}.roleNameConfig(id, roleName))
}

func (RoleAssignmentMarketplaceResource) builtinConfig(id string, roleName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {
}

data "azurerm_role_definition" "test" {
  name = "%s"
}

resource "azurerm_marketplace_role_assignment" "test" {
  name               = "%s"
  role_definition_id = "${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id

  lifecycle {
    ignore_changes = [
      role_definition_name,
    ]
  }
}
`, roleName, id)
}

func (RoleAssignmentMarketplaceResource) servicePrincipal(rInt int, roleAssignmentID string, roleName string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

provider "azurerm" {
  features {}
}

resource "azurerm_marketplace_role_assignment" "test" {
  name                 = "%s"
  role_definition_name = "%s"
  principal_id         = azuread_service_principal.test.id

  lifecycle {
    ignore_changes = [
      role_definition_id,
    ]
  }
}
`, rInt, roleAssignmentID, roleName)
}

func (RoleAssignmentMarketplaceResource) servicePrincipalWithType(rInt int, roleAssignmentID string, roleName string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

provider "azurerm" {
  features {}
}

resource "azurerm_marketplace_role_assignment" "test" {
  name                             = "%s"
  role_definition_name             = "%s"
  principal_id                     = azuread_service_principal.test.id
  skip_service_principal_aad_check = true

  lifecycle {
    ignore_changes = [
      role_definition_id,
    ]
  }
}
`, rInt, roleAssignmentID, roleName)
}

func (RoleAssignmentMarketplaceResource) group(rInt int, roleAssignmentID string, roleName string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_group" "test" {
  display_name     = "acctestspa-%d"
  security_enabled = true
}

provider "azurerm" {
  features {}
}

resource "azurerm_marketplace_role_assignment" "test" {
  name                 = "%s"
  role_definition_name = "%s"
  principal_id         = azuread_group.test.id

  lifecycle {
    ignore_changes = [
      role_definition_id,
    ]
  }
}
`, rInt, roleAssignmentID, roleName)
}

func (r RoleAssignmentMarketplaceResource) deleteAssignedRole(t *testing.T, roleName string) func() {
	return func() {
		clientManager, err := testclient.Build()
		if err != nil {
			t.Fatalf("building client: %+v", err)
		}

		ctx, cancel := context.WithDeadline(clientManager.StopContext, time.Now().Add(30*time.Minute))
		defer cancel()

		roleDefinitionsClient := clientManager.Authorization.ScopedRoleDefinitionsClient
		roleDefinitions, err := roleDefinitionsClient.List(ctx, commonids.NewScopeID(authorization.MarketplaceScope), roledefinitions.ListOperationOptions{Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", roleName))})
		if err != nil {
			t.Fatalf("loading Role Definition List: %+v", err)
		}

		if roleDefinitions.Model == nil || len(*roleDefinitions.Model) != 1 || (*roleDefinitions.Model)[0].Id == nil {
			t.Fatalf("loading Role Definition List: failed to find role '%s'", roleName)
		}

		roleAssignmentsClient := clientManager.Authorization.ScopedRoleAssignmentsClient
		roleAssignments, err := roleAssignmentsClient.ListForScope(ctx, commonids.NewScopeID(authorization.MarketplaceScope), roleassignments.DefaultListForScopeOperationOptions())
		if err != nil {
			t.Fatalf("loading Role Assignment List: %+v", err)
		}

		if roleAssignments.Model == nil || len(*roleAssignments.Model) == 0 {
			return
		}

		for _, roleAssignment := range *roleAssignments.Model {
			if roleAssignment.Id == nil || roleAssignment.Properties == nil || roleAssignment.Properties.RoleDefinitionId != *(*roleDefinitions.Model)[0].Id || pointer.From(roleAssignment.Properties.Scope) != authorization.MarketplaceScope {
				continue
			}

			id, err := parse.ScopedRoleAssignmentID(*roleAssignment.Id)
			if err != nil {
				t.Fatalf("parsing scoped role assignment id: %+v", err)
			}

			options := roleassignments.DefaultDeleteOperationOptions()
			if id.TenantId != "" {
				options.TenantId = &id.TenantId
			}

			resp, err := roleAssignmentsClient.Delete(ctx, id.ScopedId, options)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					t.Fatalf("deleting role assignment: %+v", err)
				}
			}
		}
	}
}
