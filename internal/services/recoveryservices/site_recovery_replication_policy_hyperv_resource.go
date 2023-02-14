package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ReplicationPolicyHyperVModel struct {
	Name                                          string `tfschema:"name"`
	ResourceGroupName                             string `tfschema:"resource_group_name"`
	RecoveryVaultName                             string `tfschema:"recovery_vault_name"`
	RecoveryPointRetentionInHours                 int64  `tfschema:"recovery_point_retention_in_hours"`
	ApplicationConsistentSnapshotFrequencyInHours int64  `tfschema:"application_consistent_snapshot_frequency_in_hours"`
	CopyFrequency                                 int64  `tfschema:"replication_interval_in_seconds"`
}

type ReplicationPolicyHyperVResource struct{}

var _ sdk.ResourceWithUpdate = ReplicationPolicyHyperVResource{}

func (r ReplicationPolicyHyperVResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},
		"recovery_point_retention_in_hours": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntBetween(0, 24),
		},
		"application_consistent_snapshot_frequency_in_hours": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntBetween(0, 12),
		},
		"replication_interval_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntInSlice([]int{30, 300}),
		},
	}
}

func (r ReplicationPolicyHyperVResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r ReplicationPolicyHyperVResource) ModelObject() interface{} {
	return &ReplicationPolicyHyperVModel{}
}

func (r ReplicationPolicyHyperVResource) ResourceType() string {
	return "azurerm_site_recovery_replication_policy_hyperv"
}

func (r ReplicationPolicyHyperVResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan ReplicationPolicyHyperVModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			id := replicationpolicies.NewReplicationPolicyID(subscriptionId, plan.ResourceGroupName, plan.RecoveryVaultName, plan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
				if !response.WasNotFound(existing.HttpResponse) && !wasBadRequestWithNotExist(existing.HttpResponse, err) {
					return fmt.Errorf("checking presence %s: %+v", plan.Name, err)
				}
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := replicationpolicies.CreatePolicyInput{
				Properties: &replicationpolicies.CreatePolicyInputProperties{
					ProviderSpecificInput: &replicationpolicies.HyperVReplicaAzurePolicyInput{
						RecoveryPointHistoryDuration:                  &plan.RecoveryPointRetentionInHours,
						ApplicationConsistentSnapshotFrequencyInHours: &plan.ApplicationConsistentSnapshotFrequencyInHours,
						ReplicationInterval:                           &plan.CopyFrequency,
					},
				},
			}
			err = client.CreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating site recovery replication policy %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ReplicationPolicyHyperVResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.MarkAsGone(id)
					return nil
				}
				return fmt.Errorf("reading %s : %+v", id.String(), err)
			}

			state := ReplicationPolicyHyperVModel{
				Name:              id.ReplicationPolicyName,
				ResourceGroupName: id.ResourceGroupName,
				RecoveryVaultName: id.VaultName,
			}

			if model := resp.Model; model != nil {
				if detail, isA2A := expandH2APolicyDetail(resp.Model); isA2A {
					if detail.ApplicationConsistentSnapshotFrequencyInHours != nil {
						state.ApplicationConsistentSnapshotFrequencyInHours = *detail.ApplicationConsistentSnapshotFrequencyInHours
					}
					if detail.ReplicationInterval != nil {
						state.CopyFrequency = *detail.ReplicationInterval
					}
					if detail.RecoveryPointHistoryDurationInHours != nil {
						state.RecoveryPointRetentionInHours = *detail.RecoveryPointHistoryDurationInHours
					}

				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r ReplicationPolicyHyperVResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ReplicationPolicyHyperVModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			parameters := replicationpolicies.UpdatePolicyInput{
				Properties: &replicationpolicies.UpdatePolicyInputProperties{
					ReplicationProviderSettings: &replicationpolicies.HyperVReplicaAzurePolicyInput{
						RecoveryPointHistoryDuration:                  &model.RecoveryPointRetentionInHours,
						ApplicationConsistentSnapshotFrequencyInHours: &model.ApplicationConsistentSnapshotFrequencyInHours,
						ReplicationInterval:                           &model.CopyFrequency,
					},
				},
			}

			err = client.UpdateThenPoll(ctx, *id, parameters)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ReplicationPolicyHyperVResource) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("deleting  %s : %+v", id, err)
			}

			return nil
		},
	}
}

func (r ReplicationPolicyHyperVResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationpolicies.ValidateReplicationPolicyID
}

func expandH2APolicyDetail(input *replicationpolicies.Policy) (out *replicationpolicies.HyperVReplicaAzurePolicyDetails, isA2A bool) {
	if input.Properties == nil {
		return nil, false
	}
	if input.Properties.ProviderSpecificDetails == nil {
		return nil, false
	}
	detail, isA2A := input.Properties.ProviderSpecificDetails.(replicationpolicies.HyperVReplicaAzurePolicyDetails)
	if isA2A {
		out = &detail
	}
	return
}
