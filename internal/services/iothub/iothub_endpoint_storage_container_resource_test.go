// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IotHubEndpointStorageContainerResource struct{}

func TestAccIotHubEndpointStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("file_name_format").HasValue("{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}"),
				check.That(data.ResourceName).Key("batch_frequency_in_seconds").HasValue("300"),
				check.That(data.ResourceName).Key("max_chunk_size_in_bytes").HasValue("314572800"),
				check.That(data.ResourceName).Key("encoding").HasValue("Avro"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointStorageContainer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

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

func TestAccIotHubEndpointStorageContainer_IotHubIdAndTwoResourceGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIotHubIdAndTwoResourceGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubEndpointStorageContainer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub_endpoint_storage_container"),
		},
	})
}

func TestAccIotHubEndpointStorageContainer_AuthenticationTypeSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

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

func TestAccIotHubEndpointStorageContainer_AuthenticationTypeUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

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

func TestAccIotHubEndpointStorageContainer_AuthenticationTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_endpoint_storage_container", "test")
	r := IotHubEndpointStorageContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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
		data.ImportStep(),
	})
}

func (IotHubEndpointStorageContainerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcont"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }

  lifecycle {
    ignore_changes = [endpoint]
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  container_name    = "acctestcont"
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IotHubEndpointStorageContainerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcont"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }

  lifecycle {
    ignore_changes = [endpoint]
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  container_name    = "acctestcont"
  connection_string = azurerm_storage_account.test.primary_blob_connection_string

  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  encoding                   = "JSON"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r IotHubEndpointStorageContainerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_storage_container" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  container_name    = "acctestcont"
  connection_string = azurerm_storage_account.test.primary_blob_connection_string

  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  encoding                   = "JSON"
}
`, r.basic(data))
}

func (IotHubEndpointStorageContainerResource) withIotHubIdAndTwoResourceGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcont"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }

  lifecycle {
    ignore_changes = [endpoint]
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctest"
  iothub_id           = azurerm_iothub.test.id

  container_name    = "acctestcont"
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t IotHubEndpointStorageContainerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.EndpointStorageContainerID(state.ID)
	if err != nil {
		return nil, err
	}

	iothub, err := clients.IoTHub.ResourceClient.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil || iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	if endpoints := iothub.Properties.Routing.Endpoints.StorageContainers; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					return utils.Bool(true), nil
				}
			}
		}
	}

	return utils.Bool(false), nil
}

func (r IotHubEndpointStorageContainerResource) authenticationTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  container_name    = azurerm_storage_container.test.name
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointStorageContainerResource) authenticationTypeSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  authentication_type = "identityBased"
  container_name      = azurerm_storage_container.test.name
  endpoint_uri        = azurerm_storage_account.test.primary_blob_endpoint

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_system,
  ]
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointStorageContainerResource) authenticationTypeUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_id           = azurerm_iothub.test.id
  name                = "acctest"

  authentication_type = "identityBased"
  identity_id         = azurerm_user_assigned_identity.test.id
  container_name      = azurerm_storage_container.test.name
  endpoint_uri        = azurerm_storage_account.test.primary_blob_endpoint
}
`, r.authenticationTemplate(data))
}

func (r IotHubEndpointStorageContainerResource) authenticationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestcont"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_user" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_user,
  ]

  lifecycle {
    ignore_changes = [endpoint]
  }
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_system" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary)
}
