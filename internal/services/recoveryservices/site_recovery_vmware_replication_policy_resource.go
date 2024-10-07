// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const EnableMultiVMSyncEnabled string = "True"

type VMWareReplicationPolicyResource struct{}

var _ sdk.ResourceWithCustomizeDiff = VMWareReplicationPolicyResource{}

type SiteRecoveryReplicationPolicyVmwareModel struct {
	Name                                            string `tfschema:"name"`
	RecoveryVaultID                                 string `tfschema:"recovery_vault_id"`
	ApplicationConsistentSnapshotFrequencyInMinutes int64  `tfschema:"application_consistent_snapshot_frequency_in_minutes"`
	RecoveryPointRetentionInMinutes                 int64  `tfschema:"recovery_point_retention_in_minutes"`
}

func (r VMWareReplicationPolicyResource) ModelObject() interface{} {
	return &SiteRecoveryReplicationPolicyVmwareModel{}
}

func (r VMWareReplicationPolicyResource) ResourceType() string {
	return "azurerm_site_recovery_vmware_replication_policy"
}
func (r VMWareReplicationPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ReplicationPolicyID
}

func (r VMWareReplicationPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringDoesNotContainAny("<>~',\"!@$^#_ (){}[]|`+=;,.*%&:\\?/"),
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: vaults.ValidateVaultID,
		},

		"recovery_point_retention_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntBetween(0, 15*24*60),
		},

		"application_consistent_snapshot_frequency_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntBetween(0, 12*60),
		},
	}
}

func (r VMWareReplicationPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VMWareReplicationPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan SiteRecoveryReplicationPolicyVmwareModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			parsedVaultId, err := vaults.ParseVaultID(plan.RecoveryVaultID)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", plan.RecoveryVaultID, err)
			}

			id := replicationpolicies.NewReplicationPolicyID(parsedVaultId.SubscriptionId, parsedVaultId.ResourceGroupName, parsedVaultId.VaultName, plan.Name)

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			existing, err := client.Get(ctx, id)
			if err != nil {
				// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
				if !response.WasNotFound(existing.HttpResponse) && !wasBadRequestWithNotExist(existing.HttpResponse, err) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_site_recovery_vmware_replication_policy", id.ID())
			}

			recoveryPoint := plan.RecoveryPointRetentionInMinutes
			appConsistency := plan.ApplicationConsistentSnapshotFrequencyInMinutes
			if appConsistency > recoveryPoint {
				return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
			}

			parameters := replicationpolicies.CreatePolicyInput{
				Properties: &replicationpolicies.CreatePolicyInputProperties{
					ProviderSpecificInput: &replicationpolicies.InMageRcmPolicyCreationInput{
						RecoveryPointHistoryInMinutes:     &recoveryPoint,
						AppConsistentFrequencyInMinutes:   &appConsistency,
						CrashConsistentFrequencyInMinutes: utils.Int64(10),
						EnableMultiVMSync:                 utils.String(EnableMultiVMSyncEnabled),
					},
				},
			}

			err = client.CreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r VMWareReplicationPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s : %+v", id, err)
			}

			vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

			state := SiteRecoveryReplicationPolicyVmwareModel{
				Name:            id.ReplicationPolicyName,
				RecoveryVaultID: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if resp.Model.Properties != nil && resp.Model.Properties.ProviderSpecificDetails != nil {
					if inMageRcm, ok := resp.Model.Properties.ProviderSpecificDetails.(replicationpolicies.InMageRcmPolicyDetails); ok {
						if inMageRcm.RecoveryPointHistoryInMinutes != nil {
							state.RecoveryPointRetentionInMinutes = *inMageRcm.RecoveryPointHistoryInMinutes
						}
						if inMageRcm.AppConsistentFrequencyInMinutes != nil {
							state.ApplicationConsistentSnapshotFrequencyInMinutes = *inMageRcm.AppConsistentFrequencyInMinutes
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r VMWareReplicationPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan SiteRecoveryReplicationPolicyVmwareModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			recoveryPoint := plan.RecoveryPointRetentionInMinutes
			appConsistency := plan.ApplicationConsistentSnapshotFrequencyInMinutes
			if appConsistency > recoveryPoint {
				return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
			}

			parameters := replicationpolicies.UpdatePolicyInput{
				Properties: &replicationpolicies.UpdatePolicyInputProperties{
					ReplicationProviderSettings: &replicationpolicies.InMageRcmPolicyCreationInput{
						RecoveryPointHistoryInMinutes:     &recoveryPoint,
						AppConsistentFrequencyInMinutes:   &appConsistency,
						EnableMultiVMSync:                 utils.String(EnableMultiVMSyncEnabled),
						CrashConsistentFrequencyInMinutes: utils.Int64(10),
					},
				},
			}
			err = client.UpdateThenPoll(ctx, *id, parameters)
			if err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}

}

func (r VMWareReplicationPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s : %+v", id, err)
			}

			return nil
		},
	}
}

func (r VMWareReplicationPolicyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan SiteRecoveryReplicationPolicyVmwareModel
			if err := metadata.DecodeDiff(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			if plan.RecoveryPointRetentionInMinutes == 0 && plan.ApplicationConsistentSnapshotFrequencyInMinutes > 0 {
				return fmt.Errorf("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero")
			}

			return nil
		},
	}
}
