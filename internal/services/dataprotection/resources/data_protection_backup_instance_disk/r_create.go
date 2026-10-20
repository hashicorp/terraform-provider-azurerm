// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_disk

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupInstanceDiskCreateUpdate(d *schema.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := backupvaultresources.ParseBackupVaultID(d.Get("vault_id").(string))
	id := backupinstanceresources.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.BackupInstancesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_disk", id.ID())
			}
		}
	}

	diskId, err := commonids.ParseManagedDiskID(d.Get("disk_id").(string))
	if err != nil {
		return err
	}
	location := location.Normalize(d.Get("location").(string))
	policyId, err := basebackuppolicyresources.ParseBackupPolicyID(d.Get("backup_policy_id").(string))
	if err != nil {
		return err
	}

	snapshotSubscriptionId := subscriptionId
	if v := d.Get("snapshot_subscription_id").(string); v != "" {
		snapshotSubscriptionId = v
	}
	snapshotResourceGroupId := commonids.NewResourceGroupID(snapshotSubscriptionId, d.Get("snapshot_resource_group_name").(string))

	parameters := backupinstanceresources.BackupInstanceResource{
		Properties: &backupinstanceresources.BackupInstance{
			DataSourceInfo: backupinstanceresources.Datasource{
				DatasourceType:   pointer.To("Microsoft.Compute/disks"),
				ObjectType:       pointer.To("Datasource"),
				ResourceID:       diskId.ID(),
				ResourceLocation: pointer.To(location),
				ResourceName:     pointer.To(diskId.DiskName),
				ResourceType:     pointer.To("Microsoft.Compute/disks"),
				ResourceUri:      pointer.To(diskId.ID()),
			},
			FriendlyName: pointer.To(id.BackupInstanceName),
			PolicyInfo: backupinstanceresources.PolicyInfo{
				PolicyId: policyId.ID(),
				PolicyParameters: &backupinstanceresources.PolicyParameters{
					DataStoreParametersList: &[]backupinstanceresources.DataStoreParameters{
						backupinstanceresources.AzureOperationalStoreParameters{
							ResourceGroupId: pointer.To(snapshotResourceGroupId.ID()),
							DataStoreType:   backupinstanceresources.DataStoreTypesOperationalStore,
						},
					},
				},
			},
		},
	}

	if d.IsNewResource() {
		if err := client.BackupInstancesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating DataProtection BackupInstance (%q): %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating DataProtection BackupInstance (%q): %+v", id, err)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(backupinstanceresources.StatusConfiguringProtection)},
		Target:     []string{string(backupinstanceresources.StatusProtectionConfigured)},
		Refresh:    helpers.PolicyProtectionStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for BackupInstance(%q) policy protection to be completed: %+v", id, err)
	}

	return resourceDataProtectionBackupInstanceDiskRead(d, meta)
}
