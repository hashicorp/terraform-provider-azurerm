// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IotHubEndpointCosmosDBAccountResource struct{}

func TestAccIotHubEndpointCosmosDBAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_cosmosdb_account", "test")
	r := IotHubEndpointCosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func TestAccIotHubEndpointCosmosDBAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_cosmosdb_account", "test")
	r := IotHubEndpointCosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub_endpoint_cosmosdb_account"),
		},
	})
}

func TestAccIotHubEndpointCosmosDBAccount_authenticationTypeSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_cosmosdb_account", "test")
	r := IotHubEndpointCosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointCosmosDBAccount_authenticationTypeUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_cosmosdb_account", "test")
	r := IotHubEndpointCosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointCosmosDBAccount_authenticationTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_cosmosdb_account", "test")
	r := IotHubEndpointCosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
		{
			Config: r.authenticationTypeUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.authenticationTypeSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.authenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func TestAccIotHubEndpointCosmosDBAccount_partitionKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_cosmosdb_account", "test")
	r := IotHubEndpointCosmosDBAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithPartitionKey(data, "keyName1", "{deviceid}-{YYYY}-{MM}"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
		{
			Config: r.basicWithPartitionKey(data, "keyName2", "{deviceid}-{MM}-{YYYY}"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("primary_key", "secondary_key"),
	})
}

func (r IotHubEndpointCosmosDBAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_cosmosdb_account" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.endpoint.name
  iothub_id           = azurerm_iothub.test.id
  container_name      = azurerm_cosmosdb_sql_container.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  endpoint_uri        = azurerm_cosmosdb_account.test.endpoint
  primary_key         = azurerm_cosmosdb_account.test.primary_key
  secondary_key       = azurerm_cosmosdb_account.test.secondary_key
}
`, r.template(data))
}

func (r IotHubEndpointCosmosDBAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_cosmosdb_account" "import" {
  name                = azurerm_iothub_endpoint_cosmosdb_account.test.name
  resource_group_name = azurerm_iothub_endpoint_cosmosdb_account.test.resource_group_name
  iothub_id           = azurerm_iothub_endpoint_cosmosdb_account.test.iothub_id
  container_name      = azurerm_iothub_endpoint_cosmosdb_account.test.container_name
  database_name       = azurerm_iothub_endpoint_cosmosdb_account.test.database_name
  endpoint_uri        = azurerm_iothub_endpoint_cosmosdb_account.test.endpoint_uri
  primary_key         = azurerm_iothub_endpoint_cosmosdb_account.test.primary_key
  secondary_key       = azurerm_iothub_endpoint_cosmosdb_account.test.secondary_key
}
`, r.basic(data))
}

func (r IotHubEndpointCosmosDBAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.EndpointCosmosDBAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	iothub, err := clients.IoTHub.ResourceClient.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil || iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	if endpoints := iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections; endpoints != nil {
		for _, endpoint := range pointer.From(endpoints) {
			if strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
				return pointer.To(true), nil
			}
		}
	}

	return pointer.To(false), nil
}

func (r IotHubEndpointCosmosDBAccountResource) authenticationTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_cosmosdb_account" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.endpoint.name
  iothub_id           = azurerm_iothub.test.id
  container_name      = azurerm_cosmosdb_sql_container.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  endpoint_uri        = azurerm_cosmosdb_account.test.endpoint

  primary_key   = azurerm_cosmosdb_account.test.primary_key
  secondary_key = azurerm_cosmosdb_account.test.secondary_key
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointCosmosDBAccountResource) authenticationTypeSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_cosmosdb_account" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.endpoint.name
  iothub_id           = azurerm_iothub.test.id
  container_name      = azurerm_cosmosdb_sql_container.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  endpoint_uri        = azurerm_cosmosdb_account.test.endpoint

  authentication_type = "identityBased"

  depends_on = [
    azurerm_cosmosdb_sql_role_assignment.system,
  ]
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointCosmosDBAccountResource) authenticationTypeUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_cosmosdb_account" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.endpoint.name
  iothub_id           = azurerm_iothub.test.id
  container_name      = azurerm_cosmosdb_sql_container.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  endpoint_uri        = azurerm_cosmosdb_account.test.endpoint

  authentication_type = "identityBased"
  identity_id         = azurerm_user_assigned_identity.test.id

  depends_on = [
    azurerm_cosmosdb_sql_role_assignment.user,
  ]
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointCosmosDBAccountResource) authenticationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.iothub.name
  location            = azurerm_resource_group.iothub.location
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[2]d"
  resource_group_name = azurerm_resource_group.iothub.name
  location            = azurerm_resource_group.iothub.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  lifecycle {
    ignore_changes = [endpoint]
  }
}

resource "azurerm_cosmosdb_sql_role_definition" "test" {
  name                = "acctestsqlrole%[3]s"
  resource_group_name = azurerm_resource_group.endpoint.name
  account_name        = azurerm_cosmosdb_account.test.name
  assignable_scopes = [
    azurerm_cosmosdb_account.test.id,
  ]

  permissions {
    data_actions = [
      "Microsoft.DocumentDB/databaseAccounts/readMetadata",
    ]
  }
}

resource "azurerm_cosmosdb_sql_role_assignment" "system" {
  resource_group_name = azurerm_resource_group.endpoint.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = azurerm_iothub.test.identity[0].principal_id
  scope               = azurerm_cosmosdb_account.test.id
}

resource "azurerm_cosmosdb_sql_role_assignment" "user" {
  resource_group_name = azurerm_resource_group.endpoint.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = azurerm_user_assigned_identity.test.principal_id
  scope               = azurerm_cosmosdb_account.test.id
}
`, r.dependencies(data), data.RandomInteger, data.RandomString)
}

func (r IotHubEndpointCosmosDBAccountResource) basicWithPartitionKey(data acceptance.TestData, partitionKeyName string, partitionKeyTemplate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_cosmosdb_account" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.endpoint.name
  iothub_id           = azurerm_iothub.test.id
  container_name      = azurerm_cosmosdb_sql_container.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  endpoint_uri        = azurerm_cosmosdb_account.test.endpoint
  primary_key         = azurerm_cosmosdb_account.test.primary_key
  secondary_key       = azurerm_cosmosdb_account.test.secondary_key

  partition_key_name     = "%s"
  partition_key_template = "%s"
}
`, r.template(data), partitionKeyName, partitionKeyTemplate)
}

func (r IotHubEndpointCosmosDBAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[2]d"
  resource_group_name = azurerm_resource_group.iothub.name
  location            = azurerm_resource_group.iothub.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  lifecycle {
    ignore_changes = [endpoint]
  }
}
`, r.dependencies(data), data.RandomInteger)
}

func (IotHubEndpointCosmosDBAccountResource) dependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "iothub" {
  name     = "acctestRG-iothub-%[2]d"
  location = "%[1]s"
}

resource "azurerm_resource_group" "endpoint" {
  name     = "acctestRG-iothub-endpoint-%[2]d"
  location = "%[1]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.endpoint.location
  resource_group_name = azurerm_resource_group.endpoint.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.endpoint.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
}
`, data.Locations.Primary, data.RandomInteger)
}
