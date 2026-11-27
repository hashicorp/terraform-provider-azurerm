// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareRead(d *pluginsdk.ResourceData, meta any) error {
	sharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	existing, err := sharesClient.Get(ctx, *id, fileshares.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found, removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", *id, err)
	}

	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())
	d.Set("name", id.ShareName)

	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("quota", props.ShareQuota)
			// Resource Manager treats nil and "SMB" as the same and we may not get a full response here
			enabledProtocols := fileshares.EnabledProtocolsSMB
			if props.EnabledProtocols != nil {
				enabledProtocols = *props.EnabledProtocols
			}
			d.Set("enabled_protocol", string(enabledProtocols))
			d.Set("access_tier", pointer.FromEnum(props.AccessTier))
			d.Set("acl", flattenStorageShareACLs(pointer.From(props.SignedIdentifiers)))
			d.Set("metadata", helpers.FlattenMetaData(pointer.From(props.Metadata)))
		}
	}

	// TODO - The following section for `url` will need to be updated to go-azure-sdk when the Giovanni Deprecation process has been completed
	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Share %q: %v", id, err)
	}

	// Determine the file endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
	if err != nil {
		return fmt.Errorf("determining File endpoint: %v", err)
	}

	// Parse the file endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	d.Set("url", shares.NewShareID(*accountId, id.ShareName).ID())
	d.Set("rbac_scope_id", parse.NewStorageShareResourceManagerID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, "default", id.ShareName).ID())

	return nil
}

func flattenStorageShareACLs(input []fileshares.SignedIdentifier) []any {
	result := make([]any, 0)

	for _, v := range input {
		output := map[string]any{
			"id": v.Id,
			"access_policy": []any{
				map[string]any{
					"start":       v.AccessPolicy.StartTime,
					"expiry":      v.AccessPolicy.ExpiryTime,
					"permissions": v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result
}
