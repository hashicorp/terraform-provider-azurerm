// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/localusers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LocalUserResource struct{}

func TestAccLocalUser_passwordOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.passwordOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccLocalUser_sshKeyOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sshKeyOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ssh_authorized_key"),
		{
			Config: r.sshKeyOnlyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ssh_authorized_key"),
		{
			Config: r.sshKeyOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ssh_authorized_key"),
	})
}

func TestAccLocalUser_sshKeyED25519(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sshKeyED25519(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("ssh_authorized_key"),
	})
}

func TestAccLocalUser_passwordAndSSHKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.passwordAndSSHKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("password").IsNotEmpty(),
			),
		},
		data.ImportStep("password", "ssh_authorized_key"),
		{
			Config: r.passwordAndSSHKeyMoreAuthKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("password").IsNotEmpty(),
			),
		},
		data.ImportStep("password", "ssh_authorized_key"),
		{
			Config: r.sshKeyOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("password").IsEmpty(),
			),
		},
		data.ImportStep("ssh_authorized_key"),
		{
			Config: r.passwordOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccLocalUser_homeDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.homeDirectory(data, "foo"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.homeDirectory(data, "bar"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccLocalUser_permissionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.noPermissionScope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.permissionScope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.permissionScopeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccLocalUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_local_user", "test")
	r := LocalUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.passwordOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r LocalUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Storage.ResourceManager.LocalUsers

	id, err := localusers.ParseLocalUserID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r LocalUserResource) passwordOnly(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  ssh_password_enabled = true
}
`, template)
}

func (r LocalUserResource) sshKeyOnly(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name               = "user"
  storage_account_id = azurerm_storage_account.test.id
  ssh_key_enabled    = true
  ssh_authorized_key {
    description = "key1"
    key         = local.first_public_key
  }
}
`, template)
}

func (r LocalUserResource) sshKeyED25519(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name               = "user"
  storage_account_id = azurerm_storage_account.test.id
  ssh_key_enabled    = true
  ssh_authorized_key {
    description = "key1"
    key         = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINreJb2zmSALgQUjzq/vKsf05fp5kM0ZGl8YDsP1FdWc"
  }
}
`, template)
}

func (r LocalUserResource) sshKeyOnlyUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name               = "user"
  storage_account_id = azurerm_storage_account.test.id
  ssh_key_enabled    = true
  ssh_authorized_key {
    description = "key1"
    key         = local.first_public_key
  }
  ssh_authorized_key {
    description = "key2"
    key         = local.first_public_key
  }
}
`, template)
}

func (r LocalUserResource) passwordAndSSHKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  ssh_key_enabled      = true
  ssh_password_enabled = true
  ssh_authorized_key {
    description = "key1"
    key         = local.first_public_key
  }
}
`, template)
}

func (r LocalUserResource) passwordAndSSHKeyMoreAuthKeys(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  ssh_key_enabled      = true
  ssh_password_enabled = true
  ssh_authorized_key {
    description = "key1"
    key         = local.first_public_key
  }
  ssh_authorized_key {
    description = "key2"
    key         = local.second_public_key
  }
}
`, template)
}

func (r LocalUserResource) requiresImport(data acceptance.TestData) string {
	template := r.passwordOnly(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "import" {
  name                 = azurerm_storage_account_local_user.test.name
  storage_account_id   = azurerm_storage_account_local_user.test.storage_account_id
  ssh_password_enabled = true
}
`, template)
}

func (r LocalUserResource) noPermissionScope(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  ssh_password_enabled = true
}
`, template)
}

func (r LocalUserResource) permissionScope(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  ssh_password_enabled = true
  permission_scope {
    permissions {
      read = true
    }
    service       = "blob"
    resource_name = azurerm_storage_container.test.name
  }
}
`, template)
}

func (r LocalUserResource) permissionScopeUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  ssh_password_enabled = true
  permission_scope {
    permissions {
      read   = true
      write  = true
      create = true
      delete = true
      list   = true
    }
    service       = "blob"
    resource_name = azurerm_storage_container.test.name
  }
}
`, template)
}

func (r LocalUserResource) homeDirectory(data acceptance.TestData, directory string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_local_user" "test" {
  name                 = "user"
  storage_account_id   = azurerm_storage_account.test.id
  home_directory       = "%s"
  ssh_password_enabled = true
}
`, template, directory)
}

func (r LocalUserResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}
resource "azurerm_storage_container" "test" {
  name                 = "acctestcontainer"
  storage_account_name = azurerm_storage_account.test.name
}

locals {
  first_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
