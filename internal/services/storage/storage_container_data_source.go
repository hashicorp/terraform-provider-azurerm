// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/containers"
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

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
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

			// TODO: support for ACL's, Legal Holds and Immutability Policies
			"has_immutability_policy": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

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
