// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryVMWareReplicationPolicyAssociationResource struct{}

func TestAccSiteRecoveryVMWareReplicationPolicyAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_vmware_replication_policy_association", "test")
	r := SiteRecoveryVMWareReplicationPolicyAssociationResource{}

	vaultId := os.Getenv("ARM_TEST_VMWARE_VAULT_ID")
	if vaultId == "" {
		t.Skip("Skipping as ARM_TEST_VMWARE_VAULT_ID is not specified")
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, vaultId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// association policy requires there to be a VMWare Server connected to the vault.
func (SiteRecoveryVMWareReplicationPolicyAssociationResource) basic(data acceptance.TestData, vaultId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_site_recovery_vmware_replication_policy" "test" {
  recovery_vault_id                                    = "%s"
  name                                                 = "acctest-policy-%d"
  recovery_point_retention_in_minutes                  = %d
  application_consistent_snapshot_frequency_in_minutes = %d
}

resource "azurerm_site_recovery_vmware_replication_policy_association" "test" {
  name              = "acctest-%d"
  recovery_vault_id = "%s"
  policy_id         = azurerm_site_recovery_vmware_replication_policy.test.id
}
`, vaultId, data.RandomInteger, 24*60, 4*60, data.RandomInteger, vaultId)
}

func (t SiteRecoveryVMWareReplicationPolicyAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ContainerMappingClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}
