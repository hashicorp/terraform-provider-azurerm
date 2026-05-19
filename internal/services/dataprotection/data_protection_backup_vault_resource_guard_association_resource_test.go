// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	resourceguardproxy "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupVaultResourceGuardAssociationResource struct{}

func TestAccDataProtectionBackupVaultResourceGuardAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_resource_guard_association", "test")
	r := DataProtectionBackupVaultResourceGuardAssociationResource{}

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

func TestAccDataProtectionBackupVaultResourceGuardAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_vault_resource_guard_association", "test")
	r := DataProtectionBackupVaultResourceGuardAssociationResource{}

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

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DataProtection.ResourceGuardProxyClient.DppResourceGuardProxyGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dpbv-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"
}

resource "azurerm_data_protection_resource_guard" "test" {
  name                = "acctest-dprg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_protection_backup_vault_resource_guard_association" "test" {
  data_protection_backup_vault_id  = azurerm_data_protection_backup_vault.test.id
  data_protection_resource_guard_id = azurerm_data_protection_resource_guard.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_vault_resource_guard_association" "import" {
  data_protection_backup_vault_id  = azurerm_data_protection_backup_vault_resource_guard_association.test.data_protection_backup_vault_id
  data_protection_resource_guard_id = azurerm_data_protection_backup_vault_resource_guard_association.test.data_protection_resource_guard_id
}
`, r.basic(data))
}
