// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicationPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationpolicies.NewReplicationPolicyID(subscriptionId, resGroup, vaultName, name)

	recoveryPoint := int64(d.Get("recovery_point_retention_in_minutes").(int))
	appConsistency := int64(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	if appConsistency > recoveryPoint {
		return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
	}

	parameters := replicationpolicies.UpdatePolicyInput{
		Properties: &replicationpolicies.UpdatePolicyInputProperties{
			ReplicationProviderSettings: &replicationpolicies.A2APolicyCreationInput{
				RecoveryPointHistory:            &recoveryPoint,
				AppConsistentFrequencyInMinutes: &appConsistency,
				MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
			},
		},
	}
	if err := client.UpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}
