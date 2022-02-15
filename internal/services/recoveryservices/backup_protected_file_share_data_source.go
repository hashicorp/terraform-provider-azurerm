package recoveryservices

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	recoveryServicesValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: These tests fail because enabling backup on file shares with no content
func dataSourceBackupProtectedFileShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBackupProtectedFileShareRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"source_file_share_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageShareName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: recoveryServicesValidate.RecoveryServicesVaultName,
			},

			"source_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func dataSourceBackupProtectedFileShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	// get storage account name from id
	parsedStorageAccountID, err := storageParse.StorageAccountID(d.Get("source_storage_account_id").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_storage_account_id '%s': %+v", d.Get("source_storage_account_id").(string), err)
	}
	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, parsedStorageAccountID.Name)

	id := parse.NewProtectedItemID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), "Azure", containerName, d.Get("source_file_share_name").(string))

	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.VaultName, id.ResourceGroup, "Azure", id.ProtectionContainerName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on Recovery Service Protected File Share %q (Vault %q Resource Group %q): %+v", id.Name, id.VaultName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)

	if properties := resp.Properties; properties != nil {
		if item, ok := properties.AsAzureFileshareProtectedItem(); ok {
			d.Set("source_storage_account_id", parsedStorageAccountID.ID())
			d.Set("source_file_share_name", item.FriendlyName)

			if v := item.PolicyID; v != nil {
				d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
			}
		}
	}

	return nil

}
