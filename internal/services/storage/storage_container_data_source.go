// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func dataSourceStorageContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageContainerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"container_access_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_encryption_scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_scope_override_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"metadata": MetaDataComputedSchema(),

			"has_immutability_policy": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewStorageContainerID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, containerName)

	container, err := containerClient.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	if model := container.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("name", containerName)
			d.Set("container_access_type", containerAccessTypeConversionMap[string(pointer.From(props.PublicAccess))])

			d.Set("default_encryption_scope", props.DefaultEncryptionScope)
			d.Set("encryption_scope_override_enabled", !pointer.From(props.DenyEncryptionScopeOverride))

			if err = d.Set("metadata", FlattenMetaData(pointer.From(props.Metadata))); err != nil {
				return fmt.Errorf("setting `metadata`: %v", err)
			}

			d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
			d.Set("has_legal_hold", props.HasLegalHold)
		}
	}

	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Container %q: %v", id, err)
	}

	// Determine the blob endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeBlob)
	if err != nil {
		return fmt.Errorf("determining Blob endpoint: %v", err)
	}

	// Parse the blob endpoint as a data plane account ID
	accountDpId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	d.Set("url", containers.NewContainerID(*accountDpId, id.ContainerName).ID())

	d.SetId(id.ID())

	return nil
}
