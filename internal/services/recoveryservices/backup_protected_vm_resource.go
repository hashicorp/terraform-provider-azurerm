// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRecoveryServicesBackupProtectedVM() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesBackupProtectedVMCreate,
		Read:   resourceRecoveryServicesBackupProtectedVMRead,
		Update: resourceRecoveryServicesBackupProtectedVMUpdate,
		Delete: resourceRecoveryServicesBackupProtectedVMDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protecteditems.ParseProtectedItemID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(120 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(80 * time.Minute),
		},

		Schema: resourceRecoveryServicesBackupProtectedVMSchema(),

		// It's possible to remove the associated vm from the protected backup so we'll only ForceNew this attribute if it's
		// changing to something other than empty.
		CustomizeDiff: pluginsdk.ForceNewIfChange("source_vm_id", func(ctx context.Context, old, new, meta interface{}) bool {
			return new.(string) != "" && old.(string) != new.(string)
		}),
	}
}

func resourceRecoveryServicesBackupProtectedVMCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opClient := meta.(*clients.Client).RecoveryServices.ProtectedItemOperationResultsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	// source_vm_id must be specified at creation time but can be removed during update
	if _, ok := d.GetOk("source_vm_id"); !ok {
		return fmt.Errorf("`source_vm_id` must be specified when creating")
	}

	vmId := d.Get("source_vm_id").(string)
	policyId := d.Get("backup_policy_id").(string)

	if policyId == "" {
		return fmt.Errorf("`backup_policy_id` must be specified during creation")
	}

	// get VM name from id
	parsedVmId, err := commonids.ParseVirtualMachineID(vmId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_vm_id '%s': %+v", vmId, err)
	}

	protectedItemName := fmt.Sprintf("VM;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroupName, parsedVmId.VirtualMachineName)
	containerName := fmt.Sprintf("iaasvmcontainer;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroupName, parsedVmId.VirtualMachineName)

	log.Printf("[DEBUG] Creating Azure Backup Protected VM %s (resource group %q)", protectedItemName, resourceGroup)

	id := protecteditems.NewProtectedItemID(subscriptionId, resourceGroup, vaultName, "Azure", containerName, protectedItemName)

	existing, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		isSoftDeleted := false
		if existing.Model != nil && existing.Model.Properties != nil {
			if prop, ok := existing.Model.Properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
				isSoftDeleted = pointer.From(prop.IsScheduledForDeferredDelete)
			}
		}

		if isSoftDeleted {
			if meta.(*clients.Client).Features.RecoveryServicesVault.RecoverSoftDeletedBackupProtectedVM {
				err = resourceRecoveryServicesVaultBackupProtectedVMRecoverSoftDeleted(ctx, client, opClient, id)
				if err != nil {
					return fmt.Errorf("recovering soft deleted %s: %+v", id, err)
				}
			} else {
				return fmt.Errorf(optedOutOfRecoveringSoftDeletedBackupProtectedVMFmt(parsedVmId.ID(), vaultName))
			}
		}

		if !isSoftDeleted {
			return tf.ImportAsExistsError("azurerm_backup_protected_vm", id.ID())
		}
	}

	item := protecteditems.ProtectedItemResource{
		Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
			PolicyId:           &policyId,
			WorkloadType:       pointer.To(protecteditems.DataSourceTypeVM),
			SourceResourceId:   pointer.To(vmId),
			FriendlyName:       pointer.To(parsedVmId.VirtualMachineName),
			ExtendedProperties: expandDiskExclusion(d),
			VirtualMachineId:   pointer.To(vmId),
		},
	}

	protectionState, ok := d.GetOk("protection_state")
	protectionStopped := strings.EqualFold(protectionState.(string), string(protecteditems.ProtectionStateProtectionStopped))
	requireUpdateProtectionState := ok && protectionStopped

	resp, err := client.CreateOrUpdate(ctx, id, item)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	operationId, err := parseBackupOperationId(resp.HttpResponse)
	if err != nil {
		return fmt.Errorf("issuing create request for %s: %+v", id, err)
	}

	if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, opClient, id, operationId); err != nil {
		return err
	}

	d.SetId(id.ID())

	// the protection state will be updated in the additional update.
	if requireUpdateProtectionState {
		p := protecteditems.ProtectionState(protectionState.(string))
		updateInput := protecteditems.ProtectedItemResource{
			Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
				ProtectionState:  &p,
				SourceResourceId: utils.String(vmId),
			},
		}

		resp, err = client.CreateOrUpdate(ctx, id, updateInput)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		operationId, err = parseBackupOperationId(resp.HttpResponse)
		if err != nil {
			return fmt.Errorf("issuing create request for %s: %+v", id, err)
		}

		if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, opClient, id, operationId); err != nil {
			return err
		}
	}

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func resourceRecoveryServicesBackupProtectedVMRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading %s", id)

	resp, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			if vm, ok := properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
				if vm.IsScheduledForDeferredDelete != nil && *vm.IsScheduledForDeferredDelete {
					d.SetId("")
					return nil
				}

				d.Set("source_vm_id", vm.SourceResourceId)
				d.Set("protection_state", pointer.From(vm.ProtectionState))

				backupPolicyId := ""
				if policyId := pointer.From(vm.PolicyId); policyId != "" {
					parsedPolicyId, err := protectionpolicies.ParseBackupPolicyIDInsensitively(policyId)
					if err != nil {
						return fmt.Errorf("parsing policy ID %q: %+v", policyId, err)
					}
					backupPolicyId = parsedPolicyId.ID()
				}
				d.Set("backup_policy_id", backupPolicyId)

				if v := vm.ExtendedProperties; v != nil && v.DiskExclusionProperties != nil {
					if *v.DiskExclusionProperties.IsInclusionList {
						if err := d.Set("include_disk_luns", utils.FlattenInt64Slice(v.DiskExclusionProperties.DiskLunList)); err != nil {
							return fmt.Errorf("setting include_disk_luns: %+v", err)
						}
					} else {
						if err := d.Set("exclude_disk_luns", utils.FlattenInt64Slice(v.DiskExclusionProperties.DiskLunList)); err != nil {
							return fmt.Errorf("setting exclude_disk_luns: %+v", err)
						}
					}
				}
			}
		}
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	return nil
}

func resourceRecoveryServicesBackupProtectedVMUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opClient := meta.(*clients.Client).RecoveryServices.ProtectedItemOperationResultsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		return err
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if _, ok := existing.Model.Properties.(protecteditems.AzureIaaSComputeVMProtectedItem); !ok {
		return fmt.Errorf("retrieving %s: `properties` was not a VM Protected Item", id)
	}

	model := *existing.Model
	properties := existing.Model.Properties.(protecteditems.AzureIaaSComputeVMProtectedItem)
	updateProtectedBackup := false

	if d.HasChange("backup_policy_id") {
		properties.PolicyId = pointer.To(d.Get("backup_policy_id").(string))
		updateProtectedBackup = true
	}

	if d.HasChange("exclude_disk_luns") || d.HasChange("include_disk_luns") {
		properties.ExtendedProperties = expandDiskExclusion(d)
		updateProtectedBackup = true
	}

	model.Properties = properties

	if updateProtectedBackup {
		resp, err := client.CreateOrUpdate(ctx, *id, model)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		operationId, err := parseBackupOperationId(resp.HttpResponse)
		if err != nil {
			return fmt.Errorf("issuing update request for %s: %+v", id, err)
		}

		if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, opClient, *id, operationId); err != nil {
			return err
		}
	}

	protectionState := string(pointer.From(properties.ProtectionState))
	protectionStopped := false
	if d.HasChange("protection_state") {
		protectionState = d.Get("protection_state").(string)
		protectionStopped = strings.EqualFold(protectionState, string(protecteditems.ProtectionStateProtectionStopped))
	}

	// the protection state will be updated in the additional update.
	if protectionStopped {
		p := protecteditems.ProtectionState(protectionState)
		vmId := d.Get("source_vm_id").(string)
		updateInput := protecteditems.ProtectedItemResource{
			Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
				ProtectionState:  &p,
				SourceResourceId: utils.String(vmId),
			},
		}

		resp, err := client.CreateOrUpdate(ctx, *id, updateInput)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		operationId, err := parseBackupOperationId(resp.HttpResponse)
		if err != nil {
			return fmt.Errorf("issuing update request for %s: %+v", id, err)
		}

		if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, opClient, *id, operationId); err != nil {
			return err
		}
	}

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func resourceRecoveryServicesBackupProtectedVMDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opResultClient := meta.(*clients.Client).RecoveryServices.BackupOperationResultsClient
	opClient := meta.(*clients.Client).RecoveryServices.ProtectedItemOperationResultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	if meta.(*clients.Client).Features.RecoveryService.VMBackupStopProtectionAndRetainDataOnDestroy {
		log.Printf("[DEBUG] Retaining Data and Stopping Protection for %s", id)

		existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(existing.HttpResponse) {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("making Read request on %s: %+v", id, err)
		}

		if model := existing.Model; model != nil {
			if properties := model.Properties; properties != nil {
				if vm, ok := properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
					updateInput := protecteditems.ProtectedItemResource{
						Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
							ProtectionState:  pointer.To(protecteditems.ProtectionStateProtectionStopped),
							SourceResourceId: vm.SourceResourceId,
						},
					}

					resp, err := client.CreateOrUpdate(ctx, *id, updateInput)
					if err != nil {
						return fmt.Errorf("stopping protection and retaining data for %s: %+v", id, err)
					}

					operationId, err := parseBackupOperationId(resp.HttpResponse)
					if err != nil {
						return fmt.Errorf("issuing creating/updating request for %s: %+v", id, err)
					}

					if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, opClient, *id, operationId); err != nil {
						return err
					}

					return nil
				}
			}
		}
	}

	log.Printf("[DEBUG] Deleting %s", id)

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("issuing delete request for %s: %+v", id, err)
		}
	}

	operationId, err := parseBackupOperationId(resp.HttpResponse)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx, client, opResultClient, *id, operationId); err != nil {
		return err
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx context.Context, opClient *backup.ProtectedItemOperationResultsClient, id protecteditems.ProtectedItemId, operationId string) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context was missing a deadline")
	}

	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"202"},
		Target:     []string{"200", "204"},
		Timeout:    time.Until(deadline),
		Refresh: func() (interface{}, string, error) {
			resp, err := opClient.Get(ctx, id.VaultName, id.ResourceGroupName, id.BackupFabricName, id.ProtectionContainerName, id.ProtectedItemName, operationId)
			if err != nil {
				return nil, "Error", fmt.Errorf("making Read request on Recovery Service Protected Item operation %q for %s: %+v", operationId, id, err)
			}
			return resp, strconv.Itoa(resp.StatusCode), err
		},
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx context.Context, client *protecteditems.ProtectedItemsClient, opResultClient *backup.OperationResultsClient, id protecteditems.ProtectedItemId, operationId string) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context was missing a deadline")
	}

	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Pending"},
		Target:     []string{"NotFound", "Stopped"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, "NotFound", nil
				}

				return resp, "Error", fmt.Errorf("making Read request on %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					if vm, ok := properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
						if vm.ProtectionState != nil && strings.EqualFold(string(*vm.ProtectionState), string(backup.ProtectionStateProtectionStopped)) {
							return resp, "Stopped", nil
						}
					}
				}
			}
			return resp, "Pending", nil
		},

		Timeout: time.Until(deadline),
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s: %+v", id, err)
	}

	// we should also wait for the operation to complete, or it will fail when creating a new backup vm with the same vm in different vault immediately.
	opState := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"202"},
		Target:     []string{"200", "204"},
		Refresh: func() (interface{}, string, error) {
			resp, err := opResultClient.Get(ctx, id.VaultName, id.ResourceGroupName, operationId)
			if err != nil {
				return nil, "Error", fmt.Errorf("making Read request on Recovery Service Protected Item operation %q for %s: %+v", operationId, id, err)
			}
			return resp, strconv.Itoa(resp.StatusCode), err
		},

		Timeout: time.Until(deadline),
	}

	_, err = opState.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the Recovery Service Protected Item operation to be deleted for %s: %+v", id, err)
	}

	return nil
}

func parseBackupOperationId(resp *http.Response) (operationId string, err error) {
	if resp == nil {
		return "", fmt.Errorf("Response is nil")
	}

	locationURL, err := resp.Location()
	if err != nil || locationURL == nil {
		return "", fmt.Errorf("Location header missing or empty")
	}

	parsedLocation, err := azure.ParseAzureResourceID(handleAzureSdkForGoBug2824(locationURL.Path))
	if err != nil {
		return "", err
	}

	if l, ok := parsedLocation.Path["backupOperationResults"]; ok {
		return l, nil
	}

	if l, ok := parsedLocation.Path["operationResults"]; ok {
		return l, nil
	}

	return "", fmt.Errorf("Location header missing backupOperationResults")
}

func expandDiskExclusion(d *pluginsdk.ResourceData) *protecteditems.ExtendedProperties {
	if v, ok := d.GetOk("include_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &protecteditems.ExtendedProperties{
			DiskExclusionProperties: &protecteditems.DiskExclusionProperties{
				DiskLunList:     utils.ExpandInt64Slice(diskLun),
				IsInclusionList: utils.Bool(true),
			},
		}
	}

	if v, ok := d.GetOk("exclude_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &protecteditems.ExtendedProperties{
			DiskExclusionProperties: &protecteditems.DiskExclusionProperties{
				DiskLunList:     utils.ExpandInt64Slice(diskLun),
				IsInclusionList: utils.Bool(false),
			},
		}
	}
	return nil
}

func expandDiskLunList(input []interface{}) []interface{} {
	result := make([]interface{}, 0, len(input))
	for _, v := range input {
		result = append(result, v.(int))
	}
	return result
}

func resourceRecoveryServicesVaultBackupProtectedVMRecoverSoftDeleted(ctx context.Context, client *protecteditems.ProtectedItemsClient, opClient *backup.ProtectedItemOperationResultsClient, id protecteditems.ProtectedItemId) (err error) {
	resp, err := client.CreateOrUpdate(ctx, id, protecteditems.ProtectedItemResource{
		Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
			IsRehydrate: pointer.To(true),
		},
	},
	)
	if err != nil {
		return fmt.Errorf("issuing request for %s: %+v", id, err)
	}

	operationId, err := parseBackupOperationId(resp.HttpResponse)
	if err != nil {
		return err
	}

	if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, opClient, id, operationId); err != nil {
		return err
	}

	return nil
}

func optedOutOfRecoveringSoftDeletedBackupProtectedVMFmt(vmId string, vaultName string) string {
	return fmt.Sprintf(`
An existing soft-deleted Backup Protected VM exists with the source VM %q in the recovery services
vault %q, however automatically recovering this Backup Protected VM has been disabled via the 
"features" block.

Terraform can automatically recover the soft-deleted Backup Protected VM when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import".
`, vmId, vaultName)
}

func resourceRecoveryServicesBackupProtectedVMSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"source_vm_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validation.StringIsEmpty,
				azure.ValidateResourceID,
			),
			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/10357
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"backup_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: protectionpolicies.ValidateBackupPolicyID,
		},

		"exclude_disk_luns": {
			Type:          pluginsdk.TypeSet,
			ConflictsWith: []string{"include_disk_luns"},
			Optional:      true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},

		"include_disk_luns": {
			Type:          pluginsdk.TypeSet,
			ConflictsWith: []string{"exclude_disk_luns"},
			Optional:      true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},

		"protection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(backup.ProtectedItemStateIRPending),
				string(backup.ProtectedItemStateProtected),
				string(backup.ProtectedItemStateProtectionError),
				string(backup.ProtectedItemStateProtectionStopped),
				string(backup.ProtectedItemStateProtectionPaused),
				string(backup.ProtectionStateInvalid),
			}, false),
		},
	}
}
