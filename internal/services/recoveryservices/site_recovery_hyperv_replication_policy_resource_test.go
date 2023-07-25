// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryReplicationPolicyHyperVResource struct{}

func TestAccSiteRecoveryReplicationPolicyHyperV_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_hyperv_replication_policy", "test")
	r := SiteRecoveryReplicationPolicyHyperVResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 2, 1, 300),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationPolicyHyperV_noSnapshots(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_hyperv_replication_policy", "test")
	r := SiteRecoveryReplicationPolicyHyperVResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 2, 0, 30),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("recovery_point_retention_in_hours").HasValue("2"),
				check.That(data.ResourceName).Key("application_consistent_snapshot_frequency_in_hours").HasValue("0"),
				check.That(data.ResourceName).Key("replication_interval_in_seconds").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryReplicationPolicyHyperVResource) basic(data acceptance.TestData, retentionInHours int, snapshotFrequencyInHours int, replicationInterval int) string {
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

resource "azurerm_site_recovery_hyperv_replication_policy" "test" {
  recovery_vault_id                                  = azurerm_recovery_services_vault.test.id
  name                                               = "acctest-policy-%d"
  recovery_point_retention_in_hours                  = %d
  application_consistent_snapshot_frequency_in_hours = %d
  replication_interval_in_seconds                    = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, retentionInHours, snapshotFrequencyInHours, replicationInterval)
}

func (t SiteRecoveryReplicationPolicyHyperVResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationpolicies.ParseReplicationPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
