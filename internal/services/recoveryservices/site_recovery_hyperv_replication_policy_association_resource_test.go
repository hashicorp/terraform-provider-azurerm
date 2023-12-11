// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoverHyperVReplicationPolicyAssociationResource struct{}

func TestAccSiteRecoveryHyperVReplicationPolicyAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_hyperv_replication_policy_association", "test")
	r := SiteRecoverHyperVReplicationPolicyAssociationResource{}
	hostResource := HyperVHostTestResource{}
	adminPwd := GenerateRandomPassword(10)

	data.ResourceTest(t, r, append(hostResource.PrepareHostTestSteps(data, adminPwd), []acceptance.TestStep{
		{
			Config: r.basic(data, adminPwd),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	}...))
}

func (SiteRecoverHyperVReplicationPolicyAssociationResource) basic(data acceptance.TestData, adminPwd string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_hyperv_replication_policy" "test" {
  recovery_vault_id                                  = azurerm_recovery_services_vault.test.id
  name                                               = "acctest-policy-%d"
  recovery_point_retention_in_hours                  = 2
  application_consistent_snapshot_frequency_in_hours = 1
  replication_interval_in_seconds                    = 300
}

resource "azurerm_site_recovery_hyperv_replication_policy_association" "test" {
  name           = "test-association"
  hyperv_site_id = azurerm_site_recovery_services_vault_hyperv_site.test.id
  policy_id      = azurerm_site_recovery_hyperv_replication_policy.test.id
}
`, HyperVHostTestResource{}.template(data, adminPwd), data.RandomInteger)
}

func (t SiteRecoverHyperVReplicationPolicyAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ContainerMappingClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
