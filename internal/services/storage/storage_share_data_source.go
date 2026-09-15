// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
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

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
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

			"enabled_protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"rbac_scope_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	sharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareName := d.Get("name").(string)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := fileshares.NewShareID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, shareName)

	share, err := sharesClient.Get(ctx, id, fileshares.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("name", shareName)
	d.Set("storage_account_id", accountId.ID())

	if model := share.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("quota", props.ShareQuota)
			enabledProtocols := fileshares.EnabledProtocolsSMB
			if props.EnabledProtocols != nil {
				enabledProtocols = *props.EnabledProtocols
			}
			d.Set("enabled_protocol", string(enabledProtocols))
			d.Set("access_tier", string(pointer.From(props.AccessTier)))
			d.Set("acl", flattenStorageShareACLs(pointer.From(props.SignedIdentifiers)))
			d.Set("metadata", FlattenMetaData(pointer.From(props.Metadata)))
		}
	}

	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Share %q: %v", id, err)
	}

	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
	if err != nil {
		return fmt.Errorf("determining File endpoint: %v", err)
	}

	accountDataPlaneId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	d.Set("url", shares.NewShareID(*accountDataPlaneId, id.ShareName).ID())

	d.Set("rbac_scope_id", parse.NewStorageShareResourceManagerID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, "default", id.ShareName).ID())

	d.SetId(id.ID())

	return nil
}
