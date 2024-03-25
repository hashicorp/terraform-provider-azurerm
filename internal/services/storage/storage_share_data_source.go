// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/file/shares"
)

func dataSourceStorageShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageShareRead,

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

			"metadata": MetaDataComputedSchema(),

			"acl": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"access_policy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"expiry": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"permissions": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"quota": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Share %q: %s", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q for Share %q", accountName, shareName)
	}

	// The files API does not support bearer tokens (@manicminer, 2024-02-15)
	sharesDataPlaneClient, err := storageClient.FileSharesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
	if err != nil {
		return fmt.Errorf("building FileShares Client: %v", err)
	}

	// Determine the file endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
	if err != nil {
		return fmt.Errorf("determining File endpoint: %v", err)
	}

	// Parse the file endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := shares.NewShareID(*accountId, shareName).ID()

	props, err := sharesDataPlaneClient.Get(ctx, shareName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}
	if props == nil {
		return fmt.Errorf("%s was not found", id)
	}
	d.SetId(id)

	d.Set("name", shareName)
	d.Set("storage_account_name", accountName)
	d.Set("quota", props.QuotaGB)
	if err = d.Set("acl", flattenStorageShareACLs(props.ACLs)); err != nil {
		return fmt.Errorf("setting `acl`: %v", err)
	}

	if err = d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %v", err)
	}

	resourceManagerId := parse.NewStorageShareResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, "default", shareName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}
