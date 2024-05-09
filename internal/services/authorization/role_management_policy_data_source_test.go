// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type RoleManagementPolicyDataSource struct{}

func TestRoleManagementPolicyDataSource_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_management_policy", "test")
	r := RoleManagementPolicyDataSource{}

	// Ignore the dangling resource post-test as the policy remains while the group is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.resourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestRoleManagementPolicyDataSource_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_management_policy", "test")
	r := RoleManagementPolicyDataSource{}

	// Ignore the dangling resource post-test as the policy remains while the group is in a pending deletion state
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.managementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (RoleManagementPolicyDataSource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Authorization.RoleManagementPoliciesClient

	id, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("failed to retrieve role management policy with ID %q: %+v", state.ID, err)
	}

	return pointer.To(true), nil
}

func (RoleManagementPolicyDataSource) resourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]s"
  location = "%[2]s"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_resource_group.test.id
}

data "azurerm_role_management_policy" "test" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id
}
`, data.RandomString, data.Locations.Primary)
}

func (RoleManagementPolicyDataSource) managementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {}

data "azurerm_client_config" "current" {
}

resource "azurerm_management_group" "test" {
  name = "acctest%[1]s"
}

data "azurerm_role_definition" "contributor" {
  name  = "Contributor"
  scope = azurerm_management_group.test.id
}

data "azurerm_role_management_policy" "test" {
  scope              = azurerm_management_group.test.id
  role_definition_id = data.azurerm_role_definition.contributor.id
}
`, data.RandomString)
}
