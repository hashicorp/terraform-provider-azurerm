package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type VMWareReplicationPolicyClassicResource struct{}

var _ sdk.ResourceWithCustomizeDiff = VMWareReplicationPolicyClassicResource{}

type VMWareReplicationPolicyClassicModel struct {
	Name                                            string `tfschema:"name"`
	RecoveryVaultID                                 string `tfschema:"recovery_vault_id"`
	ApplicationConsistentSnapshotFrequencyInMinutes int64  `tfschema:"application_consistent_snapshot_frequency_in_minutes"`
	RecoveryPointRetentionInMinutes                 int64  `tfschema:"recovery_point_retention_in_minutes"`
	RecoveryPointThresholdInMinutes                 int64  `tfschema:"recovery_point_threshold_in_minutes"`
	IsFailBack                                      bool   `tfschema:"is_failback"`
}

func (r VMWareReplicationPolicyClassicResource) ModelObject() interface{} {
	return &VMWareReplicationPolicyClassicModel{}
}

func (r VMWareReplicationPolicyClassicResource) ResourceType() string {
	return "azurerm_site_recovery_vmware_replication_policy_classic"
}
func (r VMWareReplicationPolicyClassicResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ReplicationPolicyID
}

func (r VMWareReplicationPolicyClassicResource) Arguments() map[string]*pluginsdk.Schema {
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

		"recovery_point_threshold_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(15, 4*60),
		},

		"application_consistent_snapshot_frequency_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntBetween(0, 12*60),
		},

		"is_failback": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
	}
}

func (r VMWareReplicationPolicyClassicResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VMWareReplicationPolicyClassicResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model VMWareReplicationPolicyClassicModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			parsedVaultId, err := vaults.ParseVaultID(model.RecoveryVaultID)
			if err != nil {
				return fmt.Errorf("paring %s: %+v", model.RecoveryVaultID, err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient
			id := replicationpolicies.NewReplicationPolicyID(parsedVaultId.SubscriptionId, parsedVaultId.ResourceGroupName, parsedVaultId.VaultName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
				if !response.WasNotFound(existing.HttpResponse) && !wasBadRequestWithNotExist(existing.HttpResponse, err) {
					return fmt.Errorf("checking for presence of existing site recovery replication policy %s: %+v", model.Name, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_site_recovery_vmware_replication_policy_classic", id.ID())
			}

			recoveryPoint := model.RecoveryPointRetentionInMinutes
			appConsistency := model.ApplicationConsistentSnapshotFrequencyInMinutes
			if appConsistency > recoveryPoint {
				return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
			}
			rpoThreshold := model.RecoveryPointThresholdInMinutes

			var specifcInput replicationpolicies.PolicyProviderSpecificInput
			if model.IsFailBack {
				specifcInput = replicationpolicies.InMagePolicyInput{
					RecoveryPointHistory:            &recoveryPoint,
					AppConsistentFrequencyInMinutes: &appConsistency,
					RecoveryPointThresholdInMinutes: &rpoThreshold,
					MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
				}
			} else {
				specifcInput = replicationpolicies.InMageAzureV2PolicyInput{
					RecoveryPointHistory:            &recoveryPoint,
					AppConsistentFrequencyInMinutes: &appConsistency,
					RecoveryPointThresholdInMinutes: &rpoThreshold,
					MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
				}
			}

			parameters := replicationpolicies.CreatePolicyInput{
				Properties: &replicationpolicies.CreatePolicyInputProperties{
					ProviderSpecificInput: &specifcInput,
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

func (r VMWareReplicationPolicyClassicResource) Read() sdk.ResourceFunc {
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

			state := VMWareReplicationPolicyClassicModel{
				Name:            id.ReplicationPolicyName,
				RecoveryVaultID: vaultId.ID(),
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.ProviderSpecificDetails != nil {
				if detail, ok := resp.Model.Properties.ProviderSpecificDetails.(replicationpolicies.InMageAzureV2PolicyDetails); ok {
					if detail.RecoveryPointHistory != nil {
						state.RecoveryPointRetentionInMinutes = *detail.RecoveryPointHistory
					}
					if detail.AppConsistentFrequencyInMinutes != nil {
						state.ApplicationConsistentSnapshotFrequencyInMinutes = *detail.AppConsistentFrequencyInMinutes
					}
					if detail.RecoveryPointThresholdInMinutes != nil {
						state.RecoveryPointThresholdInMinutes = *detail.RecoveryPointThresholdInMinutes
					}
					state.IsFailBack = false
				}

				if detail, ok := resp.Model.Properties.ProviderSpecificDetails.(replicationpolicies.InMagePolicyDetails); ok {
					if detail.RecoveryPointHistory != nil {
						state.RecoveryPointRetentionInMinutes = *detail.RecoveryPointHistory
					}
					if detail.AppConsistentFrequencyInMinutes != nil {
						state.ApplicationConsistentSnapshotFrequencyInMinutes = *detail.AppConsistentFrequencyInMinutes
					}
					if detail.RecoveryPointThresholdInMinutes != nil {
						state.RecoveryPointThresholdInMinutes = *detail.RecoveryPointThresholdInMinutes
					}
					state.IsFailBack = true
				}

			}

			return metadata.Encode(&state)
		},
	}
}

func (r VMWareReplicationPolicyClassicResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model VMWareReplicationPolicyClassicModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", id, err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			recoveryPoint := model.RecoveryPointRetentionInMinutes
			appConsistency := model.ApplicationConsistentSnapshotFrequencyInMinutes
			if appConsistency > recoveryPoint {
				return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
			}
			rpoThreshold := model.RecoveryPointThresholdInMinutes
			parameters := replicationpolicies.UpdatePolicyInput{
				Properties: &replicationpolicies.UpdatePolicyInputProperties{
					ReplicationProviderSettings: &replicationpolicies.InMageAzureV2PolicyInput{
						RecoveryPointHistory:            &recoveryPoint,
						AppConsistentFrequencyInMinutes: &appConsistency,
						RecoveryPointThresholdInMinutes: &rpoThreshold,
						MultiVMSyncStatus:               replicationpolicies.SetMultiVMSyncStatusEnable,
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

func (r VMWareReplicationPolicyClassicResource) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("deleting site recovery replication policy %s : %+v", id.String(), err)
			}

			return nil
		},
	}
}

func (r VMWareReplicationPolicyClassicResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model VMWareReplicationPolicyClassicModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			if model.RecoveryPointRetentionInMinutes == 0 && model.ApplicationConsistentSnapshotFrequencyInMinutes > 0 {
				return fmt.Errorf("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero")
			}

			return nil
		},
	}
}
