// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %s", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Account %q for Share %q", accountName, shareName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building FileShares Client for Storage Account %q (Resource Group %q): %s", accountName, account.ResourceGroup, err)
	}

	id := parse.NewStorageShareDataPlaneId(accountName, storageClient.Environment.StorageEndpointSuffix, shareName).ID()
	props, err := client.Get(ctx, account.ResourceGroup, accountName, shareName)
	if err != nil {
		return fmt.Errorf("retrieving Share %q (Account %q / Resource Group %q): %s", shareName, accountName, account.ResourceGroup, err)
	}
	if props == nil {
		return fmt.Errorf("share %q was not found in Account %q / Resource Group %q", shareName, accountName, account.ResourceGroup)
	}
	d.SetId(id)

	d.Set("name", shareName)
	d.Set("storage_account_name", accountName)
	d.Set("quota", props.QuotaGB)
	if err := d.Set("acl", flattenStorageShareACLs(props.ACLs)); err != nil {
		return fmt.Errorf("setting `acl`: %+v", err)
	}

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %+v", err)
	}

	resourceManagerId := parse.NewStorageShareResourceManagerID(storageClient.SubscriptionId, account.ResourceGroup, accountName, "default", shareName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}
