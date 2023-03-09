package recoveryservices_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryReplicationPolicyV2AResource struct{}

func TestAccSiteRecoveryReplicationPolicyV2A_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_vmware_replication_policy_classic", "test")
	r := SiteRecoveryReplicationPolicyV2AResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 24*60, 4*60, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, 24*60, 4*60, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationPolicyV2A_noSnapshots(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_vmware_replication_policy_classic", "test")
	r := SiteRecoveryReplicationPolicyV2AResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 48*60, 0, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("recovery_point_retention_in_minutes").HasValue("2880"),
				check.That(data.ResourceName).Key("application_consistent_snapshot_frequency_in_minutes").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationPolicyV2A_wrongSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_vmware_replication_policy_classic", "test")
	r := SiteRecoveryReplicationPolicyV2AResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basic(data, 0, 3, false),
			ExpectError: regexp.MustCompile("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero"),
		},
	})
}

func (SiteRecoveryReplicationPolicyV2AResource) basic(data acceptance.TestData, retentionInMinutes int, snapshotFrequencyInMinutes int, isFailback bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                               = "acctest-vault-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  sku                                = "Standard"
  classic_vmware_replication_enabled = true

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_vmware_replication_policy_classic" "test" {
  recovery_vault_id                                    = azurerm_recovery_services_vault.test.id
  name                                                 = "acctest-policy-%d"
  recovery_point_retention_in_minutes                  = %d
  application_consistent_snapshot_frequency_in_minutes = %d
  recovery_point_threshold_in_minutes                  = 60
  is_failback                                          = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, retentionInMinutes, snapshotFrequencyInMinutes, isFailback)
}

func (t SiteRecoveryReplicationPolicyV2AResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationpolicies.ParseReplicationPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery replication policy (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}
