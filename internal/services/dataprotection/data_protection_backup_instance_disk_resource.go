// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backuppolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupInstanceDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupInstanceDiskCreateUpdate,
		Read:   resourceDataProtectionBackupInstanceDiskRead,
		Update: resourceDataProtectionBackupInstanceDiskCreateUpdate,
		Delete: resourceDataProtectionBackupInstanceDiskDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := backupinstances.ParseBackupInstanceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: backupinstances.ValidateBackupVaultID,
			},

			"disk_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: disks.ValidateDiskID,
			},

			"snapshot_resource_group_name": commonschema.ResourceGroupName(),

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: backuppolicies.ValidateBackupPolicyID,
			},
		},
	}
}

func resourceDataProtectionBackupInstanceDiskCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := backupinstances.ParseBackupVaultID(d.Get("vault_id").(string))
	id := backupinstances.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_disk", id.ID())
		}
	}

	diskId, _ := disks.ParseDiskID(d.Get("disk_id").(string))
	location := location.Normalize(d.Get("location").(string))
	policyId, _ := backuppolicies.ParseBackupPolicyID(d.Get("backup_policy_id").(string))
	snapshotResourceGroupId := resourceParse.NewResourceGroupID(subscriptionId, d.Get("snapshot_resource_group_name").(string))

	parameters := backupinstances.BackupInstanceResource{
		Properties: &backupinstances.BackupInstance{
			DataSourceInfo: backupinstances.Datasource{
				DatasourceType:   utils.String("Microsoft.Compute/disks"),
				ObjectType:       utils.String("Datasource"),
				ResourceID:       diskId.ID(),
				ResourceLocation: utils.String(location),
				ResourceName:     utils.String(diskId.DiskName),
				ResourceType:     utils.String("Microsoft.Compute/disks"),
				ResourceUri:      utils.String(diskId.ID()),
			},
			FriendlyName: utils.String(id.BackupInstanceName),
			PolicyInfo: backupinstances.PolicyInfo{
				PolicyId: policyId.ID(),
				PolicyParameters: &backupinstances.PolicyParameters{
					DataStoreParametersList: &[]backupinstances.DataStoreParameters{
						backupinstances.AzureOperationalStoreParameters{
							ResourceGroupId: utils.String(snapshotResourceGroupId.ID()),
							DataStoreType:   backupinstances.DataStoreTypesOperationalStore,
						},
					},
				},
			},
		},
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupInstance (%q): %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(backupinstances.StatusConfiguringProtection)},
		Target:     []string{string(backupinstances.StatusProtectionConfigured)},
		Refresh:    policyProtectionStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for BackupInstance(%q) policy protection to be completed: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupInstanceDiskRead(d, meta)
}

func resourceDataProtectionBackupInstanceDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstances.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := backupinstances.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupInstanceName)
	d.Set("vault_id", vaultId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("disk_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)

			d.Set("backup_policy_id", props.PolicyInfo.PolicyId)
			if props.PolicyInfo.PolicyParameters != nil && props.PolicyInfo.PolicyParameters.DataStoreParametersList != nil && len(*props.PolicyInfo.PolicyParameters.DataStoreParametersList) > 0 {
				parameter := (*props.PolicyInfo.PolicyParameters.DataStoreParametersList)[0].(backupinstances.AzureOperationalStoreParameters)

				if parameter.ResourceGroupId != nil {
					resourceGroupId, err := resourceParse.ResourceGroupIDInsensitively(*parameter.ResourceGroupId)
					if err != nil {
						return err
					}
					d.Set("snapshot_resource_group_name", resourceGroupId.ResourceGroup)
				}
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupInstanceDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstances.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting DataProtection BackupInstance (%q): %+v", id, err)
	}

	return nil
}
