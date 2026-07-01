// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MachineLearningRegistryResource struct{}

func TestAccMachineLearningRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryResource{}

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

func TestAccMachineLearningRegistry_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryResource{}

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

func TestAccMachineLearningRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryResource{}

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

func TestAccMachineLearningRegistry_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistryResource{}

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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MachineLearningRegistryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	registryClient := client.MachineLearning.RegistryManagement
	id, err := registrymanagement.ParseRegistryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := registryClient.RegistriesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Registry %q: %+v", state.ID, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MachineLearningRegistryResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestMLR-%[2]d"

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r MachineLearningRegistryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_registry" "import" {
  name                = azurerm_machine_learning_registry.test.name
  location            = azurerm_machine_learning_registry.test.location
  resource_group_name = azurerm_machine_learning_registry.test.resource_group_name

  identity {
    type = azurerm_machine_learning_registry.test.identity.0.type
  }
}
`, r.basic(data))
}

func (r MachineLearningRegistryResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  name                          = "acctestMLR-%[2]d"
  public_network_access_enabled = false

  system_created_storage_account_type                           = "Standard_ZRS"
  system_created_storage_account_hierarchical_namespace_enabled = true
  system_created_container_registry_sku                         = "Premium"

  replication_region {
    location                                                      = "%[3]s"
    system_created_storage_account_type                           = "Standard_ZRS"
    system_created_storage_account_hierarchical_namespace_enabled = true
    system_created_container_registry_sku                         = "Premium"
  }

  replication_region {
    location                                                      = "%[4]s"
    system_created_storage_account_type                           = "Standard_ZRS"
    system_created_storage_account_hierarchical_namespace_enabled = true
    system_created_container_registry_sku                         = "Premium"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "test"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}

func (r MachineLearningRegistryResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  name                          = "acctestMLR-%[2]d"
  public_network_access_enabled = true

  system_created_storage_account_type                           = "Standard_ZRS"
  system_created_storage_account_hierarchical_namespace_enabled = true
  system_created_container_registry_sku                         = "Premium"

  replication_region {
    location                                                      = "%[3]s"
    system_created_storage_account_type                           = "Standard_ZRS"
    system_created_storage_account_hierarchical_namespace_enabled = true
    system_created_container_registry_sku                         = "Premium"
  }

  replication_region {
    location                                                      = "%[4]s"
    system_created_storage_account_type                           = "Standard_ZRS"
    system_created_storage_account_hierarchical_namespace_enabled = true
    system_created_container_registry_sku                         = "Premium"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "staging"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}

func (r MachineLearningRegistryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

`, data.RandomInteger, data.Locations.Primary)
}
