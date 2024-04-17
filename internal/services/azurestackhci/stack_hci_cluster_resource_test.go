// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StackHCIClusterResource struct{}

func TestAccStackHCICluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cloud_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("service_endpoint").IsNotEmpty(),
				check.That(data.ResourceName).Key("resource_provider_object_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCICluster_basicWithoutClientId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithoutClientId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cloud_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("service_endpoint").IsNotEmpty(),
				check.That(data.ResourceName).Key("resource_provider_object_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCICluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

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

func TestAccStackHCICluster_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCICluster_systemAssignedIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCICluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCICluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCICluster_automanageConfigurationAssignment(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_cluster", "test")
	r := StackHCIClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automanageConfigurationAssignment(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("automanage_configuration_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeAutomanageConfigurationAssignment(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StackHCIClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.Clusters
	id, err := clusters.ParseClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r StackHCIClusterResource) basicWithoutClientId(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "import" {
  name                = azurerm_stack_hci_cluster.test.name
  resource_group_name = azurerm_stack_hci_cluster.test.resource_group_name
  location            = azurerm_stack_hci_cluster.test.location
  client_id           = azurerm_stack_hci_cluster.test.client_id
  tenant_id           = azurerm_stack_hci_cluster.test.tenant_id
}
`, config)
}

func (r StackHCIClusterResource) systemAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = azuread_application.test.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENv = "Test2"
  }
}
`, template, data.RandomInteger)
}

func (r StackHCIClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "test" {
  display_name = "acctestspa-%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hci-%d"
  location = "%s"
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r StackHCIClusterResource) automanageConfigurationAssignment(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_stack_hci_cluster" "test" {
  name                        = "acctest-StackHCICluster-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  client_id                   = data.azurerm_client_config.current.client_id
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  automanage_configuration_id = azurerm_automanage_configuration.test.id

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r StackHCIClusterResource) removeAutomanageConfigurationAssignment(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
