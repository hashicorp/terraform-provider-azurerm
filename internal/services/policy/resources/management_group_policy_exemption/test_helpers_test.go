// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_exemption_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ManagementGroupAssignmentTestResource struct{}

func (r ManagementGroupAssignmentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "Acceptance Test MgmtGroup %[1]d"
}
`, data.RandomInteger)
}

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicySetBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}
