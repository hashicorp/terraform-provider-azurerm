// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedRedisAccessPolicyAssignmentDataSource struct{}

func TestAccManagedRedisAccessPolicyAssignmentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_redis_access_policy_assignment", "test")
	r := ManagedRedisAccessPolicyAssignmentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestassignment%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("access_policy_name").HasValue("default"),
				check.That(data.ResourceName).Key("object_id").IsNotEmpty(),
			),
		},
	})
}

func (r ManagedRedisAccessPolicyAssignmentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_managed_redis_access_policy_assignment" "test" {
  name                = azurerm_managed_redis_access_policy_assignment.test.name
  managed_redis_name  = azurerm_managed_redis.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, ManagedRedisAccessPolicyAssignmentResource{}.basic(data))
}
