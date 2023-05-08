package appconfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppConfigurationReplicaTestResource struct{}

func TestAccAppConfigurationReplica_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_replica", "test")
	r := AppConfigurationReplicaTestResource{}

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

func TestAccAppConfigurationReplica_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration_replica", "test")
	r := AppConfigurationReplicaTestResource{}

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

func (r AppConfigurationReplicaTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicas.ParseReplicaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppConfiguration.ReplicasClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AppConfigurationReplicaTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_replica" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  location                 = local.secondary_location
  name                     = "replica${local.random_integer}"
}
`, r.template(data))
}

func (r AppConfigurationReplicaTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration_replica" "import" {
  configuration_store_id   = azurerm_app_configuration_replica.test.configuration_store_id
  location                 = azurerm_app_configuration_replica.test.location
  name                     = azurerm_app_configuration_replica.test.name
}
`, r.basic(data))
}

func (r AppConfigurationReplicaTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
    random_integer = %[1]d
	primary_location	   = %[2]q
    secondary_location = %[3]q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appconfig-replica-${local.random_integer}"
  location = local.primary_location
}

resource "azurerm_app_configuration" "test" {
  name                       = "testaccappconf%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku                        = "standard"
  soft_delete_retention_days = 1
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
