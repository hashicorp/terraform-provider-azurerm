// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryProtectionContainerResource struct{}

func TestAccSiteRecoveryProtectionContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_protection_container", "test")
	r := SiteRecoveryProtectionContainerResource{}

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

func (SiteRecoveryProtectionContainerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_fabric" "test" {
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric-%d"
  location            = azurerm_resource_group.test.location
}

resource "azurerm_site_recovery_protection_container" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test.name
  name                 = "acctest-protection-cont-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t SiteRecoveryProtectionContainerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotectioncontainers.ParseReplicationProtectionContainerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectionContainerClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery protection container (%s): %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("reading site recovery protection container (%s): model is nil", id.String())
	}

	return utils.Bool(model.Id != nil), nil
}
