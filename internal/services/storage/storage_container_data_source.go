// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func dataSourceStorageContainer() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
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
		},
	}

	if !features.FivePointOhBeta() {
		r.Schema["resource_manager_id"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Computed:   true,
			Deprecated: "this property has been deprecated in favour of `id` and will be removed in version 5.0 of the Provider.",
		}

		r.Schema["storage_account_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_account_id",
			},
		}

		r.Schema["storage_account_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_account_id",
			},
		}
	}

	return r
}

func dataSourceStorageContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)

	if !features.FivePointOhBeta() {
		storageClient := meta.(*clients.Client).Storage
		accountName := d.Get("storage_account_name").(string)
		if accountName != "" {
			account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
			if err != nil {
				return fmt.Errorf("retrieving Storage Account %q for Container %q: %v", accountName, containerName, err)
			}
			if account == nil {
				return fmt.Errorf("locating Storage Account %q for Container %q", accountName, containerName)
			}

			containersDataPlaneClient, err := storageClient.ContainersDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Containers Client: %v", err)
			}

			// Determine the blob endpoint, so we can build a data plane ID
			endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeBlob)
			if err != nil {
				return fmt.Errorf("determining Blob endpoint: %v", err)
			}

			// Parse the blob endpoint as a data plane account ID
			accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Account ID: %v", err)
			}

			id := containers.NewContainerID(*accountId, containerName)

			props, err := containersDataPlaneClient.Get(ctx, containerName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			if props == nil {
				return fmt.Errorf("retrieving %s: result was nil", id)
			}

			d.SetId(id.ID())

			d.Set("name", containerName)
			d.Set("storage_account_name", accountName)
			d.Set("container_access_type", flattenStorageContainerAccessLevel(props.AccessLevel))

			d.Set("default_encryption_scope", props.DefaultEncryptionScope)
			d.Set("encryption_scope_override_enabled", !props.EncryptionScopeOverrideDisabled)

			if err = d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
				return fmt.Errorf("setting `metadata`: %v", err)
			}

			d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
			d.Set("has_legal_hold", props.HasLegalHold)

			resourceManagerId := commonids.NewStorageContainerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, containerName)
			d.Set("resource_manager_id", resourceManagerId.ID())

			return nil
		}
	}

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

			if !features.FivePointOhBeta() {
				d.Set("resource_manager_id", id.ID())
			}
		}
	}

	d.SetId(id.ID())

	return nil
}
