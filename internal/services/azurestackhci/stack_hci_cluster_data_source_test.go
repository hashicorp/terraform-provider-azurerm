// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StackHCIClusterDataSource struct{}

func TestAccStackHCIClusterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("client_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("cloud_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("service_endpoint").IsNotEmpty(),
				check.That(data.ResourceName).Key("resource_provider_object_id").IsNotEmpty(),
			),
		},
	})
}

func TestAccStackHCIClusterDataSource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsNotEmpty(),
			),
		},
	})
}

func (d StackHCIClusterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_stack_hci_cluster" "test" {
  name                = azurerm_stack_hci_cluster.test.name
  resource_group_name = azurerm_stack_hci_cluster.test.resource_group_name
}
`, StackHCIClusterResource{}.basic(data))
}

func (d StackHCIClusterDataSource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_stack_hci_cluster" "test" {
  name                = azurerm_stack_hci_cluster.test.name
  resource_group_name = azurerm_stack_hci_cluster.test.resource_group_name
}
`, StackHCIClusterResource{}.systemAssignedIdentity(data))
}
