// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIClusterConfigurationAssignmentResource struct{}

func TestAccStackHCIClusterConfigurationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster_automanage_configuration_assignment", "test")
	r := StackHCIClusterConfigurationAssignmentResource{}
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

func TestAccStackHCIClusterConfigurationAssignment_requireImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster_automanage_configuration_assignment", "test")
	r := StackHCIClusterConfigurationAssignmentResource{}
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

func (r StackHCIClusterConfigurationAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Automanage.ConfigurationProfileHCIAssignmentsClient

	id, err := configurationprofilehciassignments.ParseConfigurationProfileAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIClusterConfigurationAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_stack_hci_cluster_automanage_configuration_assignment" "test" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.test.id
  configuration_id     = azurerm_automanage_configuration.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r StackHCIClusterConfigurationAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_stack_hci_cluster_automanage_configuration_assignment" "import" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.test.id
  configuration_id     = azurerm_automanage_configuration.test.id
}
`, config)
}
