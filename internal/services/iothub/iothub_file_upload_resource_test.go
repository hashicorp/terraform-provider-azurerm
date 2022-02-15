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

// NOTE: this resource intentionally doesn't support Requires Import
//       since a File Upload is created by default

func TestAccIotHubFileUpload_basic(t *testing.T) {
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
	})
}

func TestAccIotHubFileUpload_update(t *testing.T) {
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
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIotHubFileUpload_authenticationTypeSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

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

func TestAccIotHubFileUpload_authenticationTypeUserAssignedIdentity(t *testing.T) {
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
	})
}

func TestAccIotHubFileUpload_authenticationTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_file_upload", "test")
	r := IotHubFileUploadResource{}

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

func (IotHubFileUploadResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FileUploadID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTHub.ResourceClient.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil || resp.Properties == nil || resp.Properties.StorageEndpoints["$default"] == nil || resp.Properties.MessagingEndpoints["fileNotifications"] == nil {
		return nil, fmt.Errorf("reading IotHub File Upload (%s): %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (IotHubFileUploadResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (IotHubFileUploadResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_file_upload" "test" {
  iothub_id          = azurerm_iothub.test.id
  connection_string  = azurerm_storage_account.test.primary_blob_connection_string
  container_name     = azurerm_storage_container.test.name
  notifications      = true
  max_delivery_count = 100
  sas_ttl            = "P1D"
  default_ttl        = "P2D"
  lock_duration      = "PT5M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r IotHubFileUploadResource) authenticationTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_file_upload" "test" {
  iothub_id         = azurerm_iothub.test.id
  connection_string = azurerm_storage_account.test.primary_blob_connection_string
  container_name    = azurerm_storage_container.test.name
}
`, r.authenticationTemplate(data))
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
    azurerm_role_assignment.test_storage_blob_data_contrib_system,
  ]
}
`, r.authenticationTemplate(data))
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
`, r.authenticationTemplate(data))
}

func (IotHubFileUploadResource) authenticationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
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

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
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
}

resource "azurerm_role_assignment" "test_storage_blob_data_contrib_system" {
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_account.test.id
  principal_id         = azurerm_iothub.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
