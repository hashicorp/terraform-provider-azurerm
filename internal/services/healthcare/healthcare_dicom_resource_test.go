// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthCareDicomResource struct{}

func TestAccHealthCareDicomResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicomResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (HealthCareDicomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dicomservices.ParseDicomServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HealthCare.HealthcareWorkspaceDicomServiceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HealthCareDicomResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "acctestdicom%[2]s"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location
  depends_on   = [azurerm_healthcare_workspace.test]
}
`, r.template(data), data.RandomString)
}

func (r HealthCareDicomResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctestfs%[2]s"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "acctestdicom%[2]s"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  cors {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = 500
    allow_credentials  = true
  }

  data_partitions_enabled = true

  encryption_key_url = azurerm_key_vault_key.test.id

  storage {
    storage_account_id = azurerm_storage_account.test.id
    file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "None"
  }
}
`, r.template(data), data.RandomString)
}

func (r HealthCareDicomResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "acctestdicom%[2]s"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = azurerm_resource_group.test.location

  cors {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = 500
    allow_credentials  = true
  }

  encryption_key_url = azurerm_key_vault_key.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "Prod"
  }
}
`, r.template(data), data.RandomString)
}

func (r HealthCareDicomResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_dicom_service" "import" {
  name         = azurerm_healthcare_dicom_service.test.name
  workspace_id = azurerm_healthcare_dicom_service.test.workspace_id
  location     = azurerm_healthcare_dicom_service.test.location
}
`, r.basic(data))
}

func (HealthCareDicomResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg%[2]s"
  location = "%[1]s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv%[2]s"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "standard"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days = 7
  purge_protection_enabled   = true

  access_policy {
    tenant_id       = data.azurerm_client_config.current.tenant_id
    object_id       = data.azurerm_client_config.current.object_id
    key_permissions = ["Create", "Delete", "Get", "Purge", "Recover", "Update", "GetRotationPolicy", "SetRotationPolicy"]
  }

  access_policy {
    tenant_id       = data.azurerm_client_config.current.tenant_id
    object_id       = azurerm_user_assigned_identity.test.principal_id
    key_permissions = ["Create", "Delete", "Get", "Import", "Purge", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[2]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "acctesthw%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.Locations.Primary, data.RandomString)
}
