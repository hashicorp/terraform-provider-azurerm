// Copyright (c) HashiCorp, Inc.
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

type MachineLearningRegistry struct{}

func TestAccMachineLearningRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistry{}

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
	r := MachineLearningRegistry{}

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
	r := MachineLearningRegistry{}

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
	r := MachineLearningRegistry{}

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

func TestAccMachineLearningRegistry_privateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_registry", "test")
	r := MachineLearningRegistry{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MachineLearningRegistry) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
		return nil, fmt.Errorf("retrieving Machine Learning Data Store %q: %+v", state.ID, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MachineLearningRegistry) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlreg-%[2]d"

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r MachineLearningRegistry) requiresImport(data acceptance.TestData) string {
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

func (r MachineLearningRegistry) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  name                          = "accmlreg-%[2]d"
  public_network_access_enabled = false

  primary_region {
    system_created_storage_account_type = "Standard_ZRS"
    hns_enabled          = true
  }

  replication_region {
    location             = "%[3]s"
    system_created_storage_account_type = "Standard_ZRS"
    hns_enabled          = true
  }

  replication_region {
    location             = "%[4]s"
    system_created_storage_account_type = "Standard_ZRS"
    hns_enabled          = true
  }

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}

func (r MachineLearningRegistry) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_machine_learning_registry" "test" {
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  name                          = "accmlreg-%[2]d"
  public_network_access_enabled = true
  primary_region {
    system_created_storage_account_type = "Standard_ZRS"
    hns_enabled          = true
  }

  replication_region {
    location             = "%[3]s"
    system_created_storage_account_type = "Standard_ZRS"
    hns_enabled          = true
  }

  replication_region {
    location             = "%[4]s"
    system_created_storage_account_type = "Standard_ZRS"
    hns_enabled          = true
  }

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, data.Locations.Ternary)
}

func (r MachineLearningRegistry) privateEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsnetendpoint-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.2.0/24"]
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.api.azureml.ms"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctest-link-%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-privatelink-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "mlregistry-privateserviceconnection"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_machine_learning_registry.test.id
    subresource_names              = ["amlregistry"]
  }

  private_dns_zone_group {
    name                 = "mlregistry-dns-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.test.id]
  }
}

resource "azurerm_machine_learning_registry" "test" {
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  name                          = "accmlreg-%[2]d"
  public_network_access_enabled = false

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MachineLearningRegistry) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}
