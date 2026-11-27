// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageShareCreate(d *pluginsdk.ResourceData, meta any) error {
	sharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := fileshares.NewShareID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := sharesClient.Get(ctx, id, fileshares.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_share", id.ID())
		}
	}

	payload := fileshares.FileShare{
		Properties: &fileshares.FileShareProperties{
			EnabledProtocols:  pointer.ToEnum[fileshares.EnabledProtocols](d.Get("enabled_protocol").(string)),
			Metadata:          pointer.To(helpers.ExpandMetaData(d.Get("metadata").(map[string]any))),
			ShareQuota:        pointer.To(int64(d.Get("quota").(int))),
			SignedIdentifiers: expandStorageShareACLs(d.Get("acl").(*pluginsdk.Set).List()),
		},
	}

	if sharedAccessTier, ok := d.GetOk("access_tier"); ok && sharedAccessTier.(string) != "" {
		payload.Properties.AccessTier = pointer.ToEnum[fileshares.ShareAccessTier](sharedAccessTier.(string))
	}

	pollerType := custompollers.NewStorageShareCreatePoller(sharesClient, id, payload)
	poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if err = poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageShareRead(d, meta)
}

func expandStorageShareACLs(input []any) *[]fileshares.SignedIdentifier {
	results := make([]fileshares.SignedIdentifier, 0)

	for _, v := range input {
		acl := v.(map[string]any)

		policies := acl["access_policy"].([]any)
		policy := policies[0].(map[string]any)

		identifier := fileshares.SignedIdentifier{
			Id: pointer.To(acl["id"].(string)),
			AccessPolicy: &fileshares.AccessPolicy{
				StartTime:  pointer.To(policy["start"].(string)),
				ExpiryTime: pointer.To(policy["expiry"].(string)),
				Permission: pointer.To(policy["permissions"].(string)),
			},
		}
		results = append(results, identifier)
	}

	return pointer.To(results)
}
