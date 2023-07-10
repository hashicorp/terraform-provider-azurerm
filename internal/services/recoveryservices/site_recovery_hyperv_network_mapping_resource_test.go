// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworkmappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryHyperVNetworkMappingResource struct{}

func TestAccSiteRecoveryHyperVNetworkMapping_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_hyperv_network_mapping", "test")
	r := SiteRecoveryHyperVNetworkMappingResource{}

	vaultId := os.Getenv("ARM_TEST_HYPERV_VAULT_ID")
	vmmName := os.Getenv("ARM_TEST_HYPERV_VMM_NAME")
	networkName := os.Getenv("ARM_TEST_HYPERV_VMM_NETWORK_NAME")

	if vaultId == "" || vmmName == "" || networkName == "" {
		t.Skip("Skipping as ARM_TEST_HYPERV_VAULT_ID, ARM_TEST_HYPERV_VMM_NAME or ARM_TEST_HYPERV_VMM_NETWORK_NAME not set")
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, vaultId, vmmName, networkName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryHyperVNetworkMappingResource) basic(data acceptance.TestData, vaultId, vmmName, networkName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "target" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "target" {
  name                = "network-%[1]d"
  resource_group_name = azurerm_resource_group.target.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_resource_group.target.location
}

resource "azurerm_site_recovery_hyperv_network_mapping" "test" {
  name                                              = "mapping-%[1]d"
  recovery_vault_id                                 = "%[3]s"
  source_system_center_virtual_machine_manager_name = "%[4]s"
  source_network_name                               = "%[5]s"
  target_network_id                                 = azurerm_virtual_network.target.id
}
`, data.RandomInteger, data.Locations.Primary, vaultId, vmmName, networkName)
}

func (t SiteRecoveryHyperVNetworkMappingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.NetworkMappingClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Network Mapping %q: %+v", id, err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving Recovery Service Network Mapping %q: `model` was nil", id)
	}

	return utils.Bool(resp.Model.Id != nil), nil
}
