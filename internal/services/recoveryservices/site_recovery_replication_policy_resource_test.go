// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryReplicationPolicyResource struct{}

func TestAccSiteRecoveryReplicationPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_policy", "test")
	r := SiteRecoveryReplicationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 24*60, 4*60),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationPolicy_noSnapshots(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_policy", "test")
	r := SiteRecoveryReplicationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 48*60, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("recovery_point_retention_in_minutes").HasValue("2880"),
				check.That(data.ResourceName).Key("application_consistent_snapshot_frequency_in_minutes").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationPolicy_wrongSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_policy", "test")
	r := SiteRecoveryReplicationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basic(data, 0, 3),
			ExpectError: regexp.MustCompile("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero"),
		},
	})
}

func (SiteRecoveryReplicationPolicyResource) basic(data acceptance.TestData, retentionInMinutes int, snapshotFrequencyInMinutes int) string {
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

resource "azurerm_site_recovery_replication_policy" "test" {
  resource_group_name                                  = azurerm_resource_group.test.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.test.name
  name                                                 = "acctest-policy-%d"
  recovery_point_retention_in_minutes                  = %d
  application_consistent_snapshot_frequency_in_minutes = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, retentionInMinutes, snapshotFrequencyInMinutes)
}

func (t SiteRecoveryReplicationPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationpolicies.ParseReplicationPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery replication policy (%s): %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("reading site recovery replication policy (%s): model is nil", id.String())
	}

	return utils.Bool(model.Id != nil), nil
}
