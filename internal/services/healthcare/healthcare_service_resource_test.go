// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	service "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthCareServiceResource struct{}

func TestAccHealthCareService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

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

func TestAccHealthCareService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

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

func TestAccHealthCareService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

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

func TestAccHealthCareService_publicNetworkAccessDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_service", "test")
	r := HealthCareServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (HealthCareServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := service.ParseServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HealthCare.HealthcareServiceClient.ServicesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Healthcare service %q (resource group: %q): %+v", id.ServiceName, id.ResourceGroupName, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (HealthCareServiceResource) basic(data acceptance.TestData) string {
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_healthcare_service" "test" {
  name                = "testacc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  access_policy_object_ids = [
    data.azurerm_client_config.current.object_id,
  ]
}
`, data.RandomInteger, location, data.RandomIntOfLength(17)) // name can only be 24 chars long
}

func (r HealthCareServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_healthcare_service" "import" {
  name                = azurerm_healthcare_service.test.name
  location            = azurerm_healthcare_service.test.location
  resource_group_name = azurerm_healthcare_service.test.resource_group_name

  access_policy_object_ids = [
    "${data.azurerm_client_config.current.object_id}",
  ]
}
`, r.basic(data))
}

func (HealthCareServiceResource) complete(data acceptance.TestData) string {
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
      "GetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "Get",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_healthcare_service" "test" {
  name                = "testacc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
    purpose     = "AcceptanceTests"
  }

  access_policy_object_ids = [
    data.azurerm_client_config.current.object_id,
  ]

  authentication_configuration {
    authority           = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}"
    audience            = "https://azurehealthcareapis.com"
    smart_proxy_enabled = true
  }

  cors_configuration {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = 500
    allow_credentials  = true
  }

  cosmosdb_throughput                   = 400
  cosmosdb_key_vault_key_versionless_id = azurerm_key_vault_key.test.versionless_id
}
`, data.RandomInteger, location, data.RandomString, data.RandomIntOfLength(17)) // name can only be 24 chars long
}

func (HealthCareServiceResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	// currently only supported in "ukwest", "northcentralus", "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

provider "azuread" {}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
      "GetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "Get",
      "UnwrapKey",
      "WrapKey",
      "GetRotationPolicy"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_healthcare_service" "test" {
  name                = "testacc%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "production"
    purpose     = "AcceptanceTests"
  }

  access_policy_object_ids = [
    data.azurerm_client_config.current.object_id,
  ]

  authentication_configuration {
    authority           = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}"
    audience            = "https://azurehealthcareapis.com"
    smart_proxy_enabled = true
  }

  cors_configuration {
    allowed_origins    = ["http://www.example.com", "http://www.example2.com"]
    allowed_headers    = ["*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = 500
    allow_credentials  = true
  }

  cosmosdb_throughput                   = 400
  cosmosdb_key_vault_key_versionless_id = azurerm_key_vault_key.test.versionless_id

  public_network_access_enabled = false
}
`, data.RandomInteger, location, data.RandomString, data.RandomIntOfLength(17)) // name can only be 24 chars long
}
