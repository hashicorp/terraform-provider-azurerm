// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IotHubFileUploadResource struct{}

func TestAccIotHubFileUpload_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authentication_type").HasValue("keyBased"),
				check.That(data.ResourceName).Key("default_ttl").HasValue("PT1H"),
				check.That(data.ResourceName).Key("lock_duration").HasValue("PT1M"),
				check.That(data.ResourceName).Key("max_delivery_count").HasValue("10"),
				check.That(data.ResourceName).Key("notifications_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("sas_ttl").HasValue("PT1H"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iothub_file_upload"),
		},
	})
}

func TestAccIotHubFileUpload_authenticationType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
			Config: r.authenticationTypeKeyBased(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_connectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithConnectionStringUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_containerName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithContainerNameUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_defaultTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultTTL(data, "PT2H"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultTTL(data, "PT3H"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_lockDuration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.lockDuration(data, "PT2M"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.lockDuration(data, "PT3M"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_maxDeliveryCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.maxDeliveryCount(data, 11),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.maxDeliveryCount(data, 12),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_notificationsEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.notificationsEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.notificationsEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_sasTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sasTTL(data, "PT2H"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.sasTTL(data, "PT3H"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IotHubFileUploadResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IotHubID(state.ID)
	if err != nil {
		return nil, err
	}

	iotHub, err := clients.IoTHub.ResourceClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if iotHub.Properties != nil && iotHub.Properties.MessagingEndpoints != nil {
		if storageEndpoint, ok := iotHub.Properties.StorageEndpoints["$default"]; ok {
			if storageEndpoint.ConnectionString != nil && *storageEndpoint.ConnectionString != "" && storageEndpoint.ContainerName != nil && *storageEndpoint.ContainerName != "" {
				return utils.Bool(true), nil
			}
		}
	}

	return utils.Bool(false), nil
}

func (r IotHubFileUploadResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name
}
`, r.template(data))
}

func (r IotHubFileUploadResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "import" {
  iothub_id         = azurerm_iothub_file_upload.test.iothub_id
  connection_string = azurerm_iothub_file_upload.test.connection_string
  container_name    = azurerm_iothub_file_upload.test.container_name
}
`, r.basic(data))
}

func (r IotHubFileUploadResource) authenticationTypeUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  authentication_type = "identityBased"
  identity_id         = azurerm_user_assigned_identity.test.id
}
`, r.authenticationTypeTemplate(data))
}

func (r IotHubFileUploadResource) authenticationTypeSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  authentication_type = "identityBased"

  depends_on = [
    azurerm_role_assignment.test_storage_blob_data_contrib_system
  ]
}
`, r.authenticationTypeTemplate(data))
}

func (r IotHubFileUploadResource) authenticationTypeKeyBased(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  authentication_type = "keyBased"
}
`, r.authenticationTypeTemplate(data))
}

func (IotHubFileUploadResource) authenticationTypeTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_user" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  lifecycle {
    ignore_changes = [
      file_upload
    ]
  }
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_system" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r IotHubFileUploadResource) basicWithConnectionStringUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test2" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test2" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test2.name
  container_access_type = "private"
}

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test2.primary_blob_connection_string
  container_name    = azurerm_storage_container.test2.name
}
`, r.template(data), data.RandomString)
}

func (r IotHubFileUploadResource) basicWithContainerNameUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test2" {
  name                  = "test2"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test2.name
}
`, r.template(data))
}

func (r IotHubFileUploadResource) defaultTTL(data acceptance.TestData, defaultTTL string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  default_ttl = "%s"
}
`, r.template(data), defaultTTL)
}

func (r IotHubFileUploadResource) lockDuration(data acceptance.TestData, lockDuration string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  lock_duration = "%s"
}
`, r.template(data), lockDuration)
}

func (r IotHubFileUploadResource) maxDeliveryCount(data acceptance.TestData, maxDeliveryCount int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  max_delivery_count = %d
}
`, r.template(data), maxDeliveryCount)
}

func (r IotHubFileUploadResource) notificationsEnabled(data acceptance.TestData, notificationsEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  notifications_enabled = %t
}
`, r.template(data), notificationsEnabled)
}

func (r IotHubFileUploadResource) sasTTL(data acceptance.TestData, sasTTL string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name

  sas_ttl = "%s"
}
`, r.template(data), sasTTL)
}

func (IotHubFileUploadResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  lifecycle {
    ignore_changes = [
      file_upload
    ]
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
