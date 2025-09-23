// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupContainerVMAppModel struct {
	ResourceGroupName string `tfschema:"resource_group_name"`
	RecoveryVaultName string `tfschema:"recovery_vault_name"`
	SourceResourceId  string `tfschema:"source_resource_id"`
	WorkloadType	  string `tfschema:"workload_type"`
}

type BackupContainerVMAppResource struct{}

var _ sdk.Resource = BackupContainerVMAppResource{}

func (r BackupContainerVMAppResource) ResourceType() string {
	return "azurerm_backup_container_vm_app"
}

func (r BackupContainerVMAppResource) ModelObject() interface{} {
	return &BackupContainerVMAppModel{}
}

func (r BackupContainerVMAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return protectioncontainers.ValidateProtectionContainerID
}

func (r BackupContainerVMAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"source_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workload_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(protectioncontainers.WorkloadTypeSAPAseDatabase),
			}, false),
		},
	}
}

func (r BackupContainerVMAppResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackupContainerVMAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupContainerVMAppModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.RecoveryServices.BackupProtectionContainersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vmId, err := commonids.ParseVirtualMachineID(model.SourceResourceId)
			if err != nil {
				return fmt.Errorf("parsing source_resource_id %q: %+v", model.SourceResourceId, err)
			}

			containerName := fmt.Sprintf("VMAppContainer;Compute;%s;%s", vmId.ResourceGroupName, vmId.VirtualMachineName)
			id := protectioncontainers.NewProtectionContainerID(subscriptionId, model.ResourceGroupName, model.RecoveryVaultName, "Azure", containerName)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_backup_container_vm_app", id.ID())
			}

			item := protectioncontainers.ProtectionContainerResource{
				Properties: &protectioncontainers.AzureVMAppContainerProtectionContainer{
					SourceResourceId:     &model.SourceResourceId,
					BackupManagementType: pointer.To(protectioncontainers.BackupManagementTypeAzureWorkload),
					WorkloadType:         (*protectioncontainers.WorkloadType)(&model.WorkloadType),
				},
			}

			if err = client.RegisterThenPoll(ctx, id, item); err != nil {
				return fmt.Errorf("registering %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BackupContainerVMAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.BackupProtectionContainersClient

			id, err := protectioncontainers.ParseProtectionContainerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("making Read request on backup protection container %s : %+v", id.String(), err)
			}

			state := BackupContainerVMAppModel{
				ResourceGroupName: id.ResourceGroupName,
				RecoveryVaultName: id.VaultName,
			}

			// Read workload type from config, because it's missing from GET response
			state.WorkloadType = metadata.ResourceData.Get("workload_type").(string)

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(protectioncontainers.AzureVMAppContainerProtectionContainer); ok {
					state.SourceResourceId = *properties.SourceResourceId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BackupContainerVMAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 80 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.BackupProtectionContainersClient
			opClient := metadata.Client.RecoveryServices.BackupOperationStatusesClient

			id, err := protectioncontainers.ParseProtectionContainerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Unregister(ctx, *id)
			if err != nil {
				return fmt.Errorf("deregistering %s: %+v", id, err)
			}

			locationURL, err := resp.HttpResponse.Location()
			if err != nil || locationURL == nil {
				return fmt.Errorf("unregistering backup protection container %s : Location header missing or empty", id.String())
			}

			opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

			parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
			if err != nil {
				return err
			}
			operationID := parsedLocation.Path["backupOperationResults"]

			deleteTimeout := metadata.ResourceData.Timeout(pluginsdk.TimeoutDelete)
			if err = resourceBackupProtectionContainerWorkloadWaitForOperation(ctx, opClient, id.VaultName, id.ResourceGroupName, operationID, deleteTimeout); err != nil {
				return err
			}

			return nil
		},
	}
}

// nolint unused - linter mistakenly things this function isn't used?
func resourceBackupProtectionContainerWorkloadWaitForOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string, timeout time.Duration) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		Delay:                     10 * time.Second,
		Pending:                   []string{"InProgress"},
		Target:                    []string{"Succeeded"},
		Refresh:                   resourceBackupProtectionContainerWorkloadCheckOperation(ctx, client, vaultName, resourceGroup, operationID),
		ContinuousTargetOccurence: 5,
		Timeout: 				   timeout,
	}

	log.Printf("[DEBUG] Waiting for backup container operation %q (Vault %q) to complete", operationID, vaultName)
	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func resourceBackupProtectionContainerWorkloadCheckOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, operationID)
		if err != nil {
			return resp, "Error", fmt.Errorf("making Read request on Recovery Service Protection Container operation %q (Vault %q in Resource Group %q): %+v", operationID, vaultName, resourceGroup, err)
		}

		if opErr := resp.Error; opErr != nil {
			errMsg := "No upstream error message"
			if opErr.Message != nil {
				errMsg = *opErr.Message
			}
			err = fmt.Errorf("recovery Service Protection Container operation status failed with status %q (Vault %q Resource Group %q Operation ID %q): %+v", resp.Status, vaultName, resourceGroup, operationID, errMsg)
		}

		return resp, string(resp.Status), err
	}
}
