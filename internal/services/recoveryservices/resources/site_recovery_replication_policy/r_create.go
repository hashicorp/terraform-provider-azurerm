// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicationPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationpolicies.NewReplicationPolicyID(subscriptionId, resGroup, vaultName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
			if !response.WasNotFound(existing.HttpResponse) && !helpers.WasBadRequestWithNotExist(existing.HttpResponse, err) {
				return fmt.Errorf("checking for presence of existing site recovery replication policy %s: %+v", name, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replication_policy", *existing.Model.Id)
		}
	}

	recoveryPoint := int64(d.Get("recovery_point_retention_in_minutes").(int))
	appConsistency := int64(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	if appConsistency > recoveryPoint {
		return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
	}
	parameters := replicationpolicies.CreatePolicyInput{
		Properties: &replicationpolicies.CreatePolicyInputProperties{
			ProviderSpecificInput: &replicationpolicies.A2APolicyCreationInput{
				RecoveryPointHistory:            &recoveryPoint,
				AppConsistentFrequencyInMinutes: &appConsistency,
				MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
			},
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func expandA2APolicyDetail(input *replicationpolicies.Policy) (out *replicationpolicies.A2APolicyDetails, isA2A bool) {
	if input.Properties == nil {
		return nil, false
	}
	if input.Properties.ProviderSpecificDetails == nil {
		return nil, false
	}
	detail, isA2A := input.Properties.ProviderSpecificDetails.(replicationpolicies.A2APolicyDetails)
	if isA2A {
		out = &detail
	}
	return
}
