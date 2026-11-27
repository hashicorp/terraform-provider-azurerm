// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobcontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

// Note: when adding new properties, ensure the `(StorageContainerListResource).listItemToBlobContainer` function is updated as well.
func resourceStorageContainerFlatten(ctx context.Context, d *pluginsdk.ResourceData, id *commonids.StorageContainerId, model *blobcontainers.BlobContainer, meta any, includeResource bool) error {
	d.Set("name", id.ContainerName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("container_access_type", containerAccessTypeConversionMap[pointer.FromEnum(props.PublicAccess)])
			d.Set("default_encryption_scope", props.DefaultEncryptionScope)
			d.Set("encryption_scope_override_enabled", !pointer.From(props.DenyEncryptionScopeOverride))
			d.Set("metadata", helpers.FlattenMetaData(pointer.From(props.Metadata)))

			d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
			d.Set("has_legal_hold", props.HasLegalHold)
		}
	}

	if includeResource {
		account, err := client.From(meta).GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
		if err != nil {
			return fmt.Errorf("retrieving Account for Container %q: %v", id, err)
		}

		// Determine the blob endpoint, so we can build a data plane ID
		endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeBlob)
		if err != nil {
			return fmt.Errorf("determining Blob endpoint: %v", err)
		}

		// Parse the blob endpoint as a data plane account ID
		accountId, err := accounts.ParseAccountID(*endpoint, client.From(meta).StorageDomainSuffix)
		if err != nil {
			return fmt.Errorf("parsing Account ID: %v", err)
		}

		d.Set("url", containers.NewContainerID(*accountId, id.ContainerName).ID())
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
