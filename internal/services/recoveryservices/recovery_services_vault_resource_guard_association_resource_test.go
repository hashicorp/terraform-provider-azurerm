// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VaultResourceGuardAssociationResource struct{}

func TestAccSiteRecoveryVaultResourceGuardAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault_resource_guard_association", "test")
	r := VaultResourceGuardAssociationResource{}

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

func (VaultResourceGuardAssociationResource) basic(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_data_protection_resource_guard" "test" {
  name                = "acctest-dprg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_recovery_services_vault_resource_guard_association" "test" {
  name              = "VaultProxy"
  vault_id          = azurerm_recovery_services_vault.test.id
  resource_guard_id = azurerm_data_protection_resource_guard.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_data_protection_resource_guard" "test" {
  name                = "acctest-dprg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_recovery_services_vault_resource_guard_association" "test" {
  vault_id          = azurerm_recovery_services_vault.test.id
  resource_guard_id = azurerm_data_protection_resource_guard.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t VaultResourceGuardAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ResourceGuardProxyClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}
