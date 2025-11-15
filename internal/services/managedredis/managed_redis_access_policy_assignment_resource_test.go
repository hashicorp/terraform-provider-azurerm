// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedRedisAccessPolicyAssignmentResource struct{}

func TestAccManagedRedisAccessPolicyAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_access_policy_assignment", "test")
	r := ManagedRedisAccessPolicyAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedRedisAccessPolicyAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_redis_access_policy_assignment", "test")
	r := ManagedRedisAccessPolicyAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ManagedRedisAccessPolicyAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databases.ParseAccessPolicyAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ManagedRedis.DatabaseClient.AccessPolicyAssignmentGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ManagedRedisAccessPolicyAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azuread_client_config" "current" {}

resource "azurerm_managed_redis_access_policy_assignment" "test" {
  name               = "acctestassignment%d"
  managed_redis_id   = azurerm_managed_redis.test.id
  access_policy_name = "default"
  object_id          = data.azuread_client_config.current.object_id
}
`, r.template(data), data.RandomInteger)
}

func (r ManagedRedisAccessPolicyAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_redis_access_policy_assignment" "import" {
  name               = azurerm_managed_redis_access_policy_assignment.test.name
  managed_redis_id   = azurerm_managed_redis_access_policy_assignment.test.managed_redis_id
  access_policy_name = azurerm_managed_redis_access_policy_assignment.test.access_policy_name
  object_id          = azurerm_managed_redis_access_policy_assignment.test.object_id
}
`, r.basic(data))
}

func (r ManagedRedisAccessPolicyAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redis-%d"
  location = "%s"
}

resource "azurerm_managed_redis" "test" {
  name                = "acctest-redis-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Balanced_B1"

  default_database {
    access_keys_authentication_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
