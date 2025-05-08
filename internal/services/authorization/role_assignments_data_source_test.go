// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RoleAssignmentsDataSource struct{}

func TestAccRoleAssignmentsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_assignments", "test")
	d := RoleAssignmentsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("principal_id").Exists(),
				check.That(data.ResourceName).Key("role_assignments.#").HasValue("1"),
			),
		},
	})
}

func (d RoleAssignmentsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_role_assignments" "test" {
  scope        = azurerm_resource_group.test.id
  principal_id = azurerm_user_assigned_identity.test.principal_id

  limit_at_scope = true

  // Account for eventual consistency in Role Assignments List operation after creating a new Role Assignment
  depends_on = [time_sleep.wait]
}
`, d.template(data))
}

func (RoleAssignmentsDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-role-assignments-%[1]d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "time_sleep" "wait" {
  create_duration = "30s"

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary)
}
