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
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/table/tables"
)

func dataSourceStorageTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageTableRead,

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

			"acl": {
				Type:     pluginsdk.TypeSet,
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

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tableName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", accountName, tableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q for Table %q", accountName, tableName)
	}

	// Determine the table endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
	if err != nil {
		return fmt.Errorf("determining Table endpoint: %v", err)
	}

	// Parse the table endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := tables.NewTableID(*accountId, tableName)

	aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	acls, err := aclClient.GetACLs(ctx, id.TableName)
	if err != nil {
		return fmt.Errorf("retrieving ACLs for %s: %v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", tableName)
	d.Set("storage_account_name", accountName)

	if err = d.Set("acl", flattenStorageTableACLs(acls)); err != nil {
		return fmt.Errorf("setting acl: %v", err)
	}

	resourceManagerId := parse.NewStorageTableResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, "default", tableName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}
