// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultCertificateContactsResource struct{}

func TestAccKeyVaultCertificateContacts_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contacts", "test")
	r := KeyVaultCertificateContactsResource{}

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

func TestAccKeyVaultCertificateContacts_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contacts", "test")
	r := KeyVaultCertificateContactsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_certificate_contacts"),
		},
	})
}

func TestAccKeyVaultCertificateContacts_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contacts", "test")
	r := KeyVaultCertificateContactsResource{}

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

func TestAccKeyVaultCertificateContacts_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contacts", "test")
	r := KeyVaultCertificateContactsResource{}

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
	})
}

func TestAccKeyVaultCertificateContacts_nonExistentVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contacts", "test")
	r := KeyVaultCertificateContactsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:             r.nonExistentVault(data),
			ExpectNonEmptyPlan: true,
			ExpectError:        regexp.MustCompile(`not found`),
		},
	})
}

func (r KeyVaultCertificateContactsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CertificateContactsID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.KeyVault.ManagementClient.GetCertificateContacts(ctx, id.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}

	return utils.Bool(resp.ContactList != nil && len(*resp.ContactList) != 0), nil
}

func (r KeyVaultCertificateContactsResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_contacts" "test" {
  key_vault_id = azurerm_key_vault.test.id

  contact {
    email = "example@example.com"
  }

  depends_on = [
    azurerm_key_vault_access_policy.test
  ]
}
`, template)
}

func (r KeyVaultCertificateContactsResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_contacts" "import" {
  key_vault_id = azurerm_key_vault_certificate_contacts.test.key_vault_id

  contact {
    email = "example@example.com"
  }

  depends_on = [
    azurerm_key_vault_access_policy.test
  ]
}
`, template)
}

func (r KeyVaultCertificateContactsResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_contacts" "test" {
  key_vault_id = azurerm_key_vault.test.id

  contact {
    email = "example@example.com"
    name  = "example"
    phone = "01234567890"
  }

  contact {
    email = "example2@example.com"
    name  = "example2"
  }

  depends_on = [
    azurerm_key_vault_access_policy.test
  ]
}
`, template)
}

func (r KeyVaultCertificateContactsResource) nonExistentVault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_certificate_contacts" "test" {
  # Must appear to be URL, but not actually exist - appending a string works
  key_vault_id = "${azurerm_key_vault.test.id}NOPE"

  contact {
    email = "example@example.com"
  }

  depends_on = [
    azurerm_key_vault_access_policy.test
  ]
}
`, template)
}

func (KeyVaultCertificateContactsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  lifecycle {
    ignore_changes = [
      contact
    ]
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  certificate_permissions = [
    "ManageContacts",
  ]

  key_permissions = [
    "Create",
  ]

  secret_permissions = [
    "Set",
  ]
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
