// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementNamedValueResource struct{}

func TestAccApiManagementNamedValue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

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

func TestAccApiManagementNamedValue_keyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementNamedValue_keyVaultSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultSystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementNamedValue_keyVaultInvalidSecretValue(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.keyVaultWithInvalidSecretValue(data),
			ExpectError: regexp.MustCompile("`secret` must be true when `value_from_key_vault` is set"),
		},
	})
}

func TestAccApiManagementNamedValue_keyVaultUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultUpdateToValue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementNamedValue_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")
	r := ApiManagementNamedValueResource{}

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
		data.ImportStep("value"),
	})
}

func (ApiManagementNamedValueResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namedvalue.ParseNamedValueID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.NamedValueClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementNamedValueResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementNamedValueResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty%d"
  value               = "Test Value"
  tags                = ["tag1", "tag2"]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementNamedValueResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty2%d"
  value               = "Test Value2"
  secret              = true
  tags                = ["tag3", "tag4"]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementNamedValueResource) keyVaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "DeleteIssuers",
    "Get",
    "GetIssuers",
    "Import",
    "List",
    "ListIssuers",
    "ManageContacts",
    "ManageIssuers",
    "SetIssuers",
    "Update",
    "Purge",
  ]
  secret_permissions = [
    "Get",
    "Delete",
    "List",
    "Purge",
    "Recover",
    "Set",
  ]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.test.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_access_policy.test]
}

resource "azurerm_key_vault_secret" "test2" {
  name         = "secret2-%[3]s"
  value        = "rick-and-morty2"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_access_policy.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r ApiManagementNamedValueResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestKeyVault%[2]d"
  secret              = true
  value_from_key_vault {
    secret_id          = azurerm_key_vault_secret.test.id
    identity_client_id = azurerm_user_assigned_identity.test.client_id
  }

  tags = ["tag1", "tag2"]

  depends_on = [azurerm_key_vault_access_policy.test2]
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementNamedValueResource) keyVaultUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestKeyVault%[2]d"
  secret              = true
  value_from_key_vault {
    secret_id          = azurerm_key_vault_secret.test2.id
    identity_client_id = azurerm_user_assigned_identity.test.client_id
  }
  tags = ["tag3", "tag4"]

  depends_on = [azurerm_key_vault_access_policy.test2]
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementNamedValueResource) keyVaultUpdateToValue(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestKeyVault%[2]d"
  secret              = false
  value               = "Key Vault to Value"
  tags                = ["tag5", "tag6"]
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementNamedValueResource) keyVaultWithInvalidSecretValue(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestKeyVault%[2]d"
  secret              = false
  value_from_key_vault {
    secret_id          = azurerm_key_vault_secret.test.id
    identity_client_id = azurerm_user_assigned_identity.test.client_id
  }

  tags = ["tag1", "tag2"]

  depends_on = [azurerm_key_vault_access_policy.test2]
}
`, r.keyVaultTemplate(data), data.RandomInteger)
}

func (r ApiManagementNamedValueResource) keyVaultSystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Apim-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Consumption_0"

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "DeleteIssuers",
    "Get",
    "GetIssuers",
    "Import",
    "List",
    "ListIssuers",
    "ManageContacts",
    "ManageIssuers",
    "SetIssuers",
    "Update",
    "Purge",
  ]
  secret_permissions = [
    "Get",
    "Delete",
    "List",
    "Purge",
    "Recover",
    "Set",
  ]
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_api_management.test.identity.0.tenant_id
  object_id    = azurerm_api_management.test.identity.0.principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[3]s"
  value        = "rick-and-morty"
  key_vault_id = azurerm_key_vault.test.id

  depends_on = [azurerm_key_vault_access_policy.test]
}

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestKeyVault%[1]d"
  secret              = true
  value_from_key_vault {
    secret_id = azurerm_key_vault_secret.test.id
  }

  depends_on = [azurerm_key_vault_access_policy.test2]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
