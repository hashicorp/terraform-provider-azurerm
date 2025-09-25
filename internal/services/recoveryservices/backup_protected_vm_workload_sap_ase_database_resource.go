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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupProtectedVMWorkloadSAPAseDatabaseModel struct {
	ResourceGroupName    string `tfschema:"resource_group_name"`
	RecoveryVaultName    string `tfschema:"recovery_vault_name"`
	BackupPolicyId       string `tfschema:"backup_policy_id"`
	SourceVMId           string `tfschema:"source_vm_id"`
	DatabaseName         string `tfschema:"database_name"`
	DatabaseInstanceName string `tfschema:"database_instance_name"`
	ProtectionState      string `tfschema:"protection_state"`
}

type BackupProtectedVMWorkloadSAPAseDatabaseResource struct{}

var _ sdk.Resource = BackupProtectedVMWorkloadSAPAseDatabaseResource{}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) ResourceType() string {
	return "azurerm_backup_protected_vm_workload_sap_ase_database"
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) ModelObject() interface{} {
	return &BackupProtectedVMWorkloadSAPAseDatabaseModel{}
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return protecteditems.ValidateProtectedItemID
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"database_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			//TODO: ValidateFunc
		},

		"database_instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			//TODO: ValidateFunc
		},

		"source_vm_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validation.StringIsEmpty,
				azure.ValidateResourceID,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"backup_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: protectionpolicies.ValidateBackupPolicyID,
		},

		// TODO: Double check with the service team if we can suspend vm workload backup
		"protection_state": helpers.BackupProtectedVMWorkloadProtectionStateSchema(),
	}
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupProtectedVMWorkloadSAPAseDatabaseModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.RecoveryServices.ProtectedItemsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vmId, err := commonids.ParseVirtualMachineID(model.SourceVMId)
			if err != nil {
				return fmt.Errorf("parsing source_vm_id %q: %+v", model.SourceVMId, err)
			}

			policyId := model.BackupPolicyId
			if policyId == "" {
				return fmt.Errorf("`backup_policy_id` must be specified during creation")
			}

			containerName := fmt.Sprintf("VMAppContainer;compute;%s;%s", vmId.ResourceGroupName, vmId.VirtualMachineName)
			protectedItemName := fmt.Sprintf("SAPAseDatabase;%s;%s", model.DatabaseInstanceName, model.DatabaseName)

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
					if metadata.Client.Features.RecoveryServicesVault.RecoverSoftDeletedBackupProtectedVMWorkloadSAPAseDatabase {
						err = resourceRecoveryServicesVaultBackupProtectedVMWorkloadRecoverSoftDeleted(ctx, client, id)
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
			item.Properties = &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
				PolicyId:             &policyId,
				BackupManagementType: pointer.To(protecteditems.BackupManagementTypeAzureWorkload),
				WorkloadType:         pointer.To(protecteditems.DataSourceTypeSAPAseDatabase),
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
				updateInput.Properties = &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
					ProtectionState: pointer.To(protecteditems.ProtectionState(protectionState)),
				}

				if err := client.CreateOrUpdateThenPoll(ctx, id, updateInput, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating protection state for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) Read() sdk.ResourceFunc {
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

			state := BackupProtectedVMWorkloadSAPAseDatabaseModel{
				ResourceGroupName: id.ResourceGroupName,
				RecoveryVaultName: id.VaultName,
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					if props, ok := properties.(protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem); ok {
						if pointer.From(props.IsScheduledForDeferredDelete) {
							return metadata.MarkAsGone(id)
						}

						state.DatabaseName = *props.FriendlyName
						state.DatabaseInstanceName = *props.ParentName

						backupPolicyId := ""
						if policyId := pointer.From(props.PolicyId); policyId != "" {
							parsedPolicyId, err := protectionpolicies.ParseBackupPolicyIDInsensitively(policyId)
							if err != nil {
								return fmt.Errorf("parsing policy ID %q: %+v", policyId, err)
							}
							backupPolicyId = parsedPolicyId.ID()
						}
						state.BackupPolicyId = backupPolicyId

						state.SourceVMId = *props.SourceResourceId
						state.ProtectionState = string(*props.ProtectionState)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectedItemsClient

			id, err := protecteditems.ParseProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectedVMWorkloadSAPAseDatabaseModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("protection_state") || metadata.ResourceData.HasChange("backup_policy_id") {
				updateInput := protecteditems.ProtectedItemResource{}

				properties := &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
					BackupManagementType: pointer.To(protecteditems.BackupManagementTypeAzureWorkload),
					WorkloadType:         pointer.To(protecteditems.DataSourceTypeSAPAseDatabase),
				}

				if metadata.ResourceData.HasChange("protection_state") {
					properties.ProtectionState = pointer.To(protecteditems.ProtectionState(model.ProtectionState))
				}

				if metadata.ResourceData.HasChange("backup_policy_id") {
					properties.PolicyId = pointer.To(model.BackupPolicyId)
				}

				updateInput.Properties = properties

				if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r BackupProtectedVMWorkloadSAPAseDatabaseResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 80 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectedItemsClient

			id, err := protecteditems.ParseProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectedVMWorkloadSAPAseDatabaseModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			features := metadata.Client.Features.RecoveryService

			if features.VMWorkloadSAPAseDatabaseBackupStopProtectionAndRetainDataOnDestroy {
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
					updateInput.Properties = &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
						ProtectionState: pointer.To(protecteditems.ProtectionStateProtectionStopped),
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

func resourceRecoveryServicesVaultBackupProtectedVMWorkloadRecoverSoftDeleted(ctx context.Context, client *protecteditems.ProtectedItemsClient, id protecteditems.ProtectedItemId) error {
	payload := protecteditems.ProtectedItemResource{
		Properties: &protecteditems.AzureVMWorkloadSAPAseDatabaseProtectedItem{
			IsRehydrate: pointer.To(true),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload, protecteditems.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("recovering soft-deleted %s: %+v", id, err)
	}

	return nil
}

func optedOutOfRecoveringSoftDeletedBackupProtectedVMWorkloadFmt(vmId string, vaultName string) string {
	return fmt.Sprintf(`
An existing soft-deleted Backup Protected VM Workload SAP Ase Database exists with the source VM %q in the recovery services
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
