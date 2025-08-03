// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AssignmentDataSource struct{}

func TestAccDataSourceAssignment_builtinPolicyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_assignment", "test")
	d := AssignmentDataSource{}
	r := ResourceGroupAssignmentTestResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtinPolicy(data, r),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("parameters").Exists(),
			),
		},
	})
}

func TestAccDataSourceAssignment_builtinPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_assignment", "test")
	d := AssignmentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.builtinPolicyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("policy_definition_id").Exists(),
				check.That(data.ResourceName).Key("description").HasValue("Description"),
				check.That(data.ResourceName).Key("display_name").HasValue("My Assignment"),
				check.That(data.ResourceName).Key("enforce").HasValue("true"),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
				check.That(data.ResourceName).Key("non_compliance_message.0.policy_definition_reference_id").HasValue("AINE_MinimumPasswordLength"),
				check.That(data.ResourceName).Key("not_scopes.#").HasValue("1"),
				check.That(data.ResourceName).Key("metadata").Exists(),
			),
		},
	})
}

func TestAccDataSourceAssignment_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_assignment", "test")
	d := AssignmentDataSource{}
	r := ResourceGroupAssignmentTestResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.systemAssignedIdentity(data, r),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		{
			Config: d.userAssignedIdentity(data, r),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			),
		},
	})
}

func (d AssignmentDataSource) builtinPolicy(data acceptance.TestData, r ResourceGroupAssignmentTestResource) string {
	config := r.withBuiltInPolicyBasic(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_assignment" "test" {
  name     = azurerm_resource_group_policy_assignment.test.name
  scope_id = azurerm_resource_group.test.id
}
`, config)
}

func (d AssignmentDataSource) builtinPolicyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest%[1]d"
  location = %[2]q
}

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_group_policy_assignment" "test" {
  name                 = "acctestpa-%[1]d"
  resource_group_id    = azurerm_resource_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  description  = "Description"
  display_name = "My Assignment"
  enforce      = true

  non_compliance_message {
    content                        = "test"
    policy_definition_reference_id = "AINE_MinimumPasswordLength"
  }

  not_scopes = [
    format("%%s/virtualMachines/testvm1", azurerm_resource_group.test.id)
  ]

  metadata = jsonencode({
    "category" : "Testing"
  })
}

data "azurerm_policy_assignment" "test" {
  name     = azurerm_resource_group_policy_assignment.test.name
  scope_id = azurerm_resource_group.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (d AssignmentDataSource) systemAssignedIdentity(data acceptance.TestData, r ResourceGroupAssignmentTestResource) string {
	config := r.systemAssignedIdentity(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_assignment" "test" {
  name     = azurerm_resource_group_policy_assignment.test.name
  scope_id = azurerm_resource_group.test.id
}
`, config)
}

func (d AssignmentDataSource) userAssignedIdentity(data acceptance.TestData, r ResourceGroupAssignmentTestResource) string {
	config := r.userAssignedIdentity(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_assignment" "test" {
  name     = azurerm_resource_group_policy_assignment.test.name
  scope_id = azurerm_resource_group.test.id
}
`, config)
}
