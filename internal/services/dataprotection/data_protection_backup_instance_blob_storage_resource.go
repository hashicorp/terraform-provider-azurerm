// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupinstanceresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupInstanceBlobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupInstanceBlobStorageCreateUpdate,
		Read:   resourceDataProtectionBackupInstanceBlobStorageRead,
		Update: resourceDataProtectionBackupInstanceBlobStorageCreateUpdate,
		Delete: resourceDataProtectionBackupInstanceBlobStorageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := backupinstanceresources.ParseBackupInstanceID(id)
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
				ValidateFunc: backupvaults.ValidateBackupVaultID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: backuppolicies.ValidateBackupPolicyID,
			},

			"storage_account_container_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"protection_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// The `storage_account_container_names` can not be removed once specified.
			pluginsdk.ForceNewIfChange("storage_account_container_names", func(ctx context.Context, old, new, _ interface{}) bool {
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),
		),
	}
}

func resourceDataProtectionBackupInstanceBlobStorageCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := backupvaults.ParseBackupVaultID(d.Get("vault_id").(string))
	id := backupinstanceresources.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	if d.IsNewResource() {
		existing, err := client.BackupInstancesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_blob_storage", id.ID())
		}
	}

	storageAccountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}
	location := location.Normalize(d.Get("location").(string))
	policyId, err := backuppolicies.ParseBackupPolicyID(d.Get("backup_policy_id").(string))
	if err != nil {
		return err
	}

	parameters := backupinstanceresources.BackupInstanceResource{
		Properties: &backupinstanceresources.BackupInstance{
			DataSourceInfo: backupinstanceresources.Datasource{
				DatasourceType:   pointer.To("Microsoft.Storage/storageAccounts/blobServices"),
				ObjectType:       pointer.To("Datasource"),
				ResourceID:       storageAccountId.ID(),
				ResourceLocation: pointer.To(location),
				ResourceName:     pointer.To(storageAccountId.StorageAccountName),
				ResourceType:     pointer.To("Microsoft.Storage/storageAccounts"),
				ResourceUri:      pointer.To(storageAccountId.ID()),
			},
			FriendlyName: pointer.To(id.BackupInstanceName),
			PolicyInfo: backupinstanceresources.PolicyInfo{
				PolicyId: policyId.ID(),
			},
		},
	}

	if v, ok := d.GetOk("storage_account_container_names"); ok {
		parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{
			BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
				backupinstanceresources.BlobBackupDatasourceParameters{
					ContainersList: pointer.From(utils.ExpandStringSlice(v.([]interface{}))),
				},
			},
		}
	}

	if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupInstance (%q): %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(backupinstanceresources.StatusConfiguringProtection), "UpdatingProtection"},
		Target:     []string{string(backupinstanceresources.StatusProtectionConfigured)},
		Refresh:    policyProtectionStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for BackupInstance(%q) policy protection to be completed: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupInstanceBlobStorageRead(d, meta)
}

func resourceDataProtectionBackupInstanceBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstanceresources.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.BackupInstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := backupvaults.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupInstanceName)
	d.Set("vault_id", vaultId.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("storage_account_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)
			d.Set("backup_policy_id", props.PolicyInfo.PolicyId)
			d.Set("protection_state", pointer.FromEnum(props.CurrentProtectionState))
			if policyParas := props.PolicyInfo.PolicyParameters; policyParas != nil {
				if dataStoreParas := policyParas.BackupDatasourceParametersList; dataStoreParas != nil {
					if dsp := pointer.From(dataStoreParas); len(dsp) > 0 {
						if parameter, ok := dsp[0].(backupinstanceresources.BlobBackupDatasourceParameters); ok {
							d.Set("storage_account_container_names", utils.FlattenStringSlice(&parameter.ContainersList))
						}
					}
				}
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupInstanceBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstanceresources.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	err = client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions())
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
