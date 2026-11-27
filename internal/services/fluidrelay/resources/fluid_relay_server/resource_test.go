// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FluidRelayServerResource struct{}

func TestAccFluidRelayServer_regressionTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayServerResource{}

	data.ResourceRegressionTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
		},
	}, "")
}

func TestAccFluidRelayServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayServerResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
				check.That(data.ResourceName).Key("frs_tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("service_endpoints.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFluidRelayServer_storageBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayServerResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.storageBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
				check.That(data.ResourceName).Key("frs_tenant_id").IsUUID(),
			),
		},
		data.ImportStep("storage_sku"),
	})
}

func TestAccFluidRelayServer_ami(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	f := FluidRelayServerResource{}

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.userAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.systemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.systemAndUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFluidRelayServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	var f FluidRelayServerResource

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.RequiresImportErrorStep(f.requiresImport),
	})
}

func TestAccFluidRelayServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	var f FluidRelayServerResource

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
	})
}

func TestAccFluidRelayServer_customerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	var f FluidRelayServerResource

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.customerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.customerManagedKeyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFluidRelayServer_customerManagedKeyVersionless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fluid_relay_server", "test")
	var f FluidRelayServerResource

	data.ResourceTest(t, f, []acceptance.TestStep{
		{
			Config: f.customerManagedKeyVersionless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
		{
			Config: f.customerManagedKeyVersionlessUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(f),
			),
		},
		data.ImportStep(),
	})
}

func (f FluidRelayServerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fluidrelayservers.ParseFluidRelayServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.FluidRelay.FluidRelayServers.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (f FluidRelayServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fluidrelay-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) templateWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-userAssignedIdentity-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, f.template(data), data.RandomInteger)
}

func (f FluidRelayServerResource) templateWithCMK(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-userAssignedIdentity-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"
    ]
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}
`, f.template(data), data.RandomInteger, data.RandomString)
}

func (f FluidRelayServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  tags = {
    foo = "bar"
  }
}
`, f.template(data), data.RandomInteger, data.Locations.Primary)
}

// basic storage sku only work with east asia and south-ease-asia
func (f FluidRelayServerResource) storageBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "SouthEastAsia"
  storage_sku         = "basic"
  tags = {
    foo = "bar"
  }
}
`, f.template(data), data.RandomInteger)
}

func (f FluidRelayServerResource) userAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, f.templateWithIdentity(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) systemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  tags = {
    foo = "bar"
  }
}
`, f.templateWithIdentity(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) systemAndUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  tags = {
    foo = "bar"
  }
}
`, f.templateWithIdentity(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_fluid_relay_server" "import" {
  name                = azurerm_fluid_relay_server.test.name
  resource_group_name = azurerm_fluid_relay_server.test.resource_group_name
  location            = azurerm_fluid_relay_server.test.location
}
`, f.basic(data))
}

func (f FluidRelayServerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  tags = {
    foo = "bar2"
  }
}
`, f.template(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) customerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, f.templateWithCMK(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) customerManagedKeyVersionless(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, f.templateWithCMK(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) customerManagedKeyVersionlessUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test2.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, f.templateWithCMK(data), data.RandomInteger, data.Locations.Primary)
}

func (f FluidRelayServerResource) customerManagedKeyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_fluid_relay_server" "test" {
  name                = "acctestRG-fuildRelayServer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test2.versionless_id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, f.templateWithCMK(data), data.RandomInteger, data.Locations.Primary)
}
