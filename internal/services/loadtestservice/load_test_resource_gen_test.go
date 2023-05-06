package loadtestservice_test

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LoadTestTestResource struct{}

func TestAccLoadTest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_load_test", "test")
	r := LoadTestTestResource{}

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

func TestAccLoadTest_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_load_test", "test")
	r := LoadTestTestResource{}

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

func TestAccLoadTest_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_load_test", "test")
	r := LoadTestTestResource{}

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

func TestAccLoadTest_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_load_test", "test")
	r := LoadTestTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}
func (r LoadTestTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadtests.ParseLoadTestID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LoadTestService.LoadTests.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
func (r LoadTestTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_load_test" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data))
}

func (r LoadTestTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_load_test" "import" {
  location            = azurerm_resource_group.test.location
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.basic(data))
}

func (r LoadTestTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-${local.random_str}"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "SetRotationPolicy",
      "GetRotationPolicy",
      "Rotate",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "Get",
      "WrapKey",
      "UnwrapKey",
      "GetRotationPolicy",
      "Decrypt",
      "Encrypt",
    ]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-${local.random_str}"
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

resource "azurerm_load_test" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  description         = "foo"

  encryption {

    identity {
      resource_id = azurerm_user_assigned_identity.test.id
      type        = "UserAssigned"
    }

    key_url = azurerm_key_vault_key.test.id
  }


  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    env  = "Production"
    test = "Acceptance"
  }
}
`, r.template(data))
}

func (r LoadTestTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  random_integer   = %[1]d
  primary_location = %[2]q
  random_str       = %[3]q
}


resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${local.random_integer}"
  location = local.primary_location
}



resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-${local.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
