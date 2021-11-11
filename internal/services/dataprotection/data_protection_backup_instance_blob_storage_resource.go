package dataprotection

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
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
			_, err := parse.BackupInstanceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": location.Schema(),

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BackupVaultID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.BackupPolicyID,
			},
		},
	}
}
func resourceDataProtectionBackupInstanceBlobStorageCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := parse.BackupVaultID(d.Get("vault_id").(string))
	id := parse.NewBackupInstanceID(subscriptionId, vaultId.ResourceGroup, vaultId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing DataProtection BackupInstance (%q): %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_instance_blob_storage", id.ID())
		}
	}

	storageAccountId, _ := storageParse.StorageAccountID(d.Get("storage_account_id").(string))
	location := location.Normalize(d.Get("location").(string))
	policyId, _ := parse.BackupPolicyID(d.Get("backup_policy_id").(string))

	parameters := dataprotection.BackupInstanceResource{
		Properties: &dataprotection.BackupInstance{
			DataSourceInfo: &dataprotection.Datasource{
				DatasourceType:   utils.String("Microsoft.Storage/storageAccounts/blobServices"),
				ObjectType:       utils.String("Datasource"),
				ResourceID:       utils.String(storageAccountId.ID()),
				ResourceLocation: utils.String(location),
				ResourceName:     utils.String(storageAccountId.Name),
				ResourceType:     utils.String("Microsoft.Storage/storageAccounts"),
				ResourceURI:      utils.String(storageAccountId.ID()),
			},
			FriendlyName: utils.String(id.Name),
			PolicyInfo: &dataprotection.PolicyInfo{
				PolicyID: utils.String(policyId.ID()),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.BackupVaultName, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupInstance (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of the DataProtection BackupInstance (%q): %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(dataprotection.StatusConfiguringProtection), "UpdatingProtection"},
		Target:     []string{string(dataprotection.StatusProtectionConfigured)},
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

	id, err := parse.BackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := parse.NewBackupVaultID(id.SubscriptionId, id.ResourceGroup, id.BackupVaultName)
	d.Set("name", id.Name)
	d.Set("vault_id", vaultId.ID())
	if props := resp.Properties; props != nil {
		if props.DataSourceInfo != nil {
			d.Set("storage_account_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)
		}
		if props.PolicyInfo != nil {
			d.Set("backup_policy_id", props.PolicyInfo.PolicyID)
		}
	}
	return nil
}

func resourceDataProtectionBackupInstanceBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting DataProtection BackupInstance (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the DataProtection BackupInstance (%q): %+v", id.Name, err)
	}
	return nil
}
