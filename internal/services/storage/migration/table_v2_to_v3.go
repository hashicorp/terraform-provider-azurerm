package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = TableV2ToV3{}

type TableV2ToV3 struct{}

func (TableV2ToV3) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"acl": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"access_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"expiry": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"permissions": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"resource_manager_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (TableV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		tableName, ok := rawState["name"].(string)
		if !ok || tableName == "" {
			return rawState, nil
		}

		accountName, ok := rawState["storage_account_name"].(string)
		if !ok || accountName == "" {
			return rawState, nil
		}

		// If resource_manager_id is already present and valid, use it
		if rmIdRaw, ok := rawState["resource_manager_id"].(string); ok && rmIdRaw != "" {
			log.Printf("[DEBUG] Updating ID from %q to Management Plane ID %q", rawState["id"], rmIdRaw)
			rawState["id"] = rmIdRaw
			
			// Extract subscription and resource group from rmId
			rmId, err := parse.StorageTableResourceManagerID(rmIdRaw)
			if err == nil {
				rawState["storage_account_id"] = commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName).ID()
			}
			
			delete(rawState, "storage_account_name")
			return rawState, nil
		}

		// Fallback: If resource_manager_id is missing, look up the account to get the Resource Group
		client := meta.(*clients.Client)
		storageClient := client.Storage
		subscriptionId := client.Account.SubscriptionId

		account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
		}
		if account == nil {
			return nil, fmt.Errorf("locating Storage Account %q", accountName)
		}

		newID := parse.NewStorageTableResourceManagerID(subscriptionId, account.StorageAccountId.ResourceGroupName, accountName, "default", tableName).ID()
		log.Printf("[DEBUG] Updating ID from %q to Management Plane ID %q", rawState["id"], newID)
		
		rawState["id"] = newID
		rawState["storage_account_id"] = commonids.NewStorageAccountID(subscriptionId, account.StorageAccountId.ResourceGroupName, accountName).ID()
		delete(rawState, "storage_account_name")
		return rawState, nil
	}
}
