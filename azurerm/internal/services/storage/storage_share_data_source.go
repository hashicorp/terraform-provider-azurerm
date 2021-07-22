package storage

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"access_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"permissions": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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
		return fmt.Errorf("error retrieving Account %q for Share %q: %s", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Account %q for Share %q", accountName, shareName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building FileShares Client for Storage Account %q (Resource Group %q): %s", accountName, account.ResourceGroup, err)
	}

	id := parse.NewStorageShareDataPlaneId(accountName, storageClient.Environment.StorageEndpointSuffix, shareName).ID()
	d.SetId(id)
	props, err := client.Get(ctx, account.ResourceGroup, accountName, shareName)
	if err != nil {
		return fmt.Errorf("retrieving Share %q (Account %q / Resource Group %q): %s", shareName, accountName, account.ResourceGroup, err)
	}
	if props == nil {
		return fmt.Errorf("share %q was not found in Account %q / Resource Group %q", shareName, accountName, account.ResourceGroup)
	}

	d.Set("name", shareName)
	d.Set("storage_account_name", accountName)
	d.Set("quota", props.QuotaGB)
	d.Set("acl", flattenStorageShareACLs(props.ACLs))

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("error setting `metadata`: %+v", err)
	}

	resourceManagerId := parse.NewStorageShareResourceManagerID(storageClient.SubscriptionId, account.ResourceGroup, accountName, "default", shareName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}
