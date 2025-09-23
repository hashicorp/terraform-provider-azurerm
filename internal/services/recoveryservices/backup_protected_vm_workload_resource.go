// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupProtectedVMWorkloadModel struct {
	ResourceGroupName string `tfschema:"resource_group_name"`
	RecoveryVaultName string `tfschema:"recovery_vault_name"`
	BackupPolicyId    string `tfschema:"backup_policy_id"`
	WorkloadType      string `tfschema:"workload_type"`
	SourceVMId        string `tfschema:"source_vm_id"`
	ProtectedItemName string `tfschema:"protected_item_name"`
	ProtectionState   string `tfschema:"protection_state"`
}

type BackupProtectedVMWorkloadResource struct{}

var _ sdk.Resource = BackupProtectedVMWorkloadResource{}

func (r BackupProtectedVMWorkloadResource) ResourceType() string {
	return "azurerm_backup_protected_vm_workload"
}

func (r BackupProtectedVMWorkloadResource) ModelObject() interface{} {
	return &BackupProtectedVMWorkloadModel{}
}

func (r BackupProtectedVMWorkloadResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return protecteditems.ValidateProtectedItemID
}

func (r BackupProtectedVMWorkloadResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protected_item_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		// TODO: Test to see if we need to deassociate source vm id as well
		"source_vm_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		// TODO: If we need to support protection_state, the backup_policy_id needs to be optional
		"backup_policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: protectionpolicies.ValidateBackupPolicyID,
		},

		"workload_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(protecteditems.DataSourceTypeSAPAseDatabase),
			}, false),
		},

		// TODO: Double check with the service team if we can suspend vm workload backup, or use empty policyId to stop protection instead
		"protection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because `protection_state` is set by Azure and may not be a persistent value.
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				// While not a persistent state, `Protected` is an option to allow a path from `ProtectionStopped` to a protected state.
				string(protecteditems.ProtectionStateProtected),
				string(protecteditems.ProtectionStateProtectionStopped),
			}, false),
			DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				// We suppress the diff if the only change is from "IRPending" or "ProtectionPaused" to "Protected".
				// These states are not persistent and are set by Azure based on the current protection state.
				// While `Invalid` and `ProtectionError` are also not configurable, we're opting to output this in the diff
				// as these states should indicate to the user that there is an error with the backup protected VM resource requiring attention.
				suppressStates := []string{
					string(protecteditems.ProtectedItemStateIRPending),
					string(protecteditems.ProtectedItemStateProtectionPaused),
				}

				if new == string(protecteditems.ProtectionStateProtected) && slices.Contains(suppressStates, old) {
					return true
				}

				return false
			},
		},
	}
}

func (r BackupProtectedVMWorkloadResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackupProtectedVMWorkloadResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupProtectedVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.RecoveryServices.ProtectedItemsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vmId, err := commonids.ParseVirtualMachineID(model.SourceVMId)
			if err != nil {
				return fmt.Errorf("parsing source_vm_id %q: %+v", model.SourceVMId, err)
			}

			containerName := fmt.Sprintf("VMAppContainer;compute;%s;%s", vmId.ResourceGroupName, vmId.VirtualMachineName)
			protectedItemName := model.ProtectedItemName

			id := protecteditems.NewProtectedItemID(subscriptionId, model.ResourceGroupName, model.RecoveryVaultName, "Azure", containerName, protectedItemName)

			existing, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				isSoftDeleted := false
				if existing.Model != nil && existing.Model.Properties != nil {
					if props, ok := existing.Model.Properties.(protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem); ok {
						isSoftDeleted = pointer.From(props.IsScheduledForDeferredDelete)
					}
				}

				if isSoftDeleted {
					if metadata.Client.Features.RecoveryServicesVault.RecoverSoftDeletedBackupProtectedVMWorkload {
						err = resourceRecoveryServicesVaultBackupProtectedVMWorkloadRecoverSoftDeleted(ctx, client, id, model.WorkloadType)
						if err != nil {
							return fmt.Errorf("recovering soft deleted %s: %+v", id, err)
						}
					} else {
						return errors.New(optedOutOfRecoveringSoftDeletedBackupProtectedVMWorkloadFmt(vmId.ID(), model.RecoveryVaultName))
					}
				}

				if !isSoftDeleted {
					return tf.ImportAsExistsError("azurerm_backup_protected_vm_workload", id.ID())
				}
			}

			item := protecteditems.ProtectedItemResource{}

			switch model.WorkloadType {
			case string(protecteditems.DataSourceTypeSAPAseDatabase):
				item.Properties = &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
					PolicyId:             &model.BackupPolicyId,
					BackupManagementType: pointer.To(protecteditems.BackupManagementTypeAzureWorkload),
					WorkloadType:         pointer.To(protecteditems.DataSourceType(model.WorkloadType)),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, item, protecteditems.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// the protection state cannot be set during initial creation.
			protectionState := model.ProtectionState
			protectionStateUpdateRequired := slices.Contains([]string{
				string(protecteditems.ProtectionStateProtectionStopped),
			}, protectionState)

			if protectionStateUpdateRequired {
				updateInput := protecteditems.ProtectedItemResource{}

				switch model.WorkloadType {
				case string(protecteditems.DataSourceTypeSAPAseDatabase):
					updateInput.Properties = &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
						ProtectionState: pointer.To(protecteditems.ProtectionState(protectionState)),
					}
				}

				if err := client.CreateOrUpdateThenPoll(ctx, id, updateInput, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating protection state for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r BackupProtectedVMWorkloadResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectedItemsClient

			id, err := protecteditems.ParseProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := BackupProtectedVMWorkloadModel{
				ResourceGroupName: id.ResourceGroupName,
				RecoveryVaultName: id.VaultName,
				ProtectedItemName: id.ProtectedItemName,
			}

			if model := resp.Model; model != nil && model.Properties != nil {
				switch props := model.Properties.(type) {
				case protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem:
					// Check for soft delete
					if pointer.From(props.IsScheduledForDeferredDelete) {
						return metadata.MarkAsGone(id)
					}

					if props.PolicyId != nil {
						state.BackupPolicyId = *props.PolicyId
					}
					if props.SourceResourceId != nil {
						state.SourceVMId = *props.SourceResourceId
					}
					if props.WorkloadType != nil {
						state.WorkloadType = string(*props.WorkloadType)
					}
					if props.ProtectionState != nil {
						state.ProtectionState = string(*props.ProtectionState)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BackupProtectedVMWorkloadResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectedItemsClient

			id, err := protecteditems.ParseProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectedVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("protection_state") || metadata.ResourceData.HasChange("backup_policy_id") {
				updateInput := protecteditems.ProtectedItemResource{}

				switch model.WorkloadType {
				case string(protecteditems.DataSourceTypeSAPAseDatabase):
					properties := &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
						BackupManagementType: pointer.To(protecteditems.BackupManagementTypeAzureWorkload),
						WorkloadType:         pointer.To(protecteditems.DataSourceType(model.WorkloadType)),
					}

					if metadata.ResourceData.HasChange("protection_state") {
						properties.ProtectionState = pointer.To(protecteditems.ProtectionState(model.ProtectionState))
					}

					if metadata.ResourceData.HasChange("backup_policy_id") {
						properties.PolicyId = pointer.To(model.BackupPolicyId)
					}

					updateInput.Properties = properties
				}

				if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r BackupProtectedVMWorkloadResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 80 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectedItemsClient

			id, err := protecteditems.ParseProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectedVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			features := metadata.Client.Features.RecoveryService

			if features.VMWorkloadBackupStopProtectionAndRetainDataOnDestroy {
				log.Printf("[DEBUG] Retaining Data and Stopping Protection for %s", id)

				existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
				if err != nil {
					if response.WasNotFound(existing.HttpResponse) {
						return nil
					}

					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}

				if existing.Model != nil && existing.Model.Properties != nil {
					updateInput := protecteditems.ProtectedItemResource{}

					switch model.WorkloadType {
					case string(protecteditems.DataSourceTypeSAPAseDatabase):
						updateInput.Properties = &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
							ProtectionState: pointer.To(protecteditems.ProtectionStateProtectionStopped),
						}
					}

					if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
						return fmt.Errorf("setting protection to %s and retaining data for %s: %+v", protecteditems.ProtectionStateProtectionStopped, *id, err)
					}

					return nil
				}
			}

			log.Printf("[DEBUG] Deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func resourceRecoveryServicesVaultBackupProtectedVMWorkloadRecoverSoftDeleted(ctx context.Context, client *protecteditems.ProtectedItemsClient, id protecteditems.ProtectedItemId, workloadType string) error {
	var payload protecteditems.ProtectedItemResource

	switch workloadType {
	case string(protecteditems.DataSourceTypeSAPAseDatabase):
		payload = protecteditems.ProtectedItemResource{
			Properties: &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
				IsRehydrate: pointer.To(true),
			},
		}
	default:
		return fmt.Errorf("unsupported workload type for soft delete recovery: %s", workloadType)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("recovering soft-deleted %s: %+v", id, err)
	}

	return nil
}

func optedOutOfRecoveringSoftDeletedBackupProtectedVMWorkloadFmt(vmId string, vaultName string) string {
	return fmt.Sprintf(`
An existing soft-deleted Backup Protected VM Workload exists with the source VM %q in the recovery services
vault %q, however automatically recovering this Backup Protected VM workload has been disabled via the 
"features" block.

Terraform can automatically recover the soft-deleted Backup Protected VM Workload when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import".
`, vmId, vaultName)
}
