package batch

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2022-01-01/batch"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceBatchAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBatchAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"location":            commonschema.LocationComputed(),
			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"pool_allocation_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"key_vault_reference": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"account_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"encryption": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemId,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceBatchAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("account_endpoint", resp.AccountEndpoint)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.AccountProperties; props != nil {
		if autoStorage := props.AutoStorage; autoStorage != nil {
			d.Set("storage_account_id", autoStorage.StorageAccountID)
		}
		d.Set("pool_allocation_mode", props.PoolAllocationMode)
		poolAllocationMode := d.Get("pool_allocation_mode").(string)

		if encryption := props.Encryption; encryption != nil {
			d.Set("encryption", flattenEncryption(encryption))
		}

		if poolAllocationMode == string(batch.PoolAllocationModeBatchService) {
			keys, err := client.GetKeys(ctx, id.ResourceGroup, id.BatchAccountName)
			if err != nil {
				return fmt.Errorf("cannot read keys for %s: %v", id, err)
			}

			d.Set("primary_access_key", keys.Primary)
			d.Set("secondary_access_key", keys.Secondary)

			// set empty keyvault reference which is not needed in Batch Service allocation mode.
			d.Set("key_vault_reference", []interface{}{})
		} else if poolAllocationMode == string(batch.PoolAllocationModeUserSubscription) {
			if err := d.Set("key_vault_reference", flattenBatchAccountKeyvaultReference(props.KeyVaultReference)); err != nil {
				return fmt.Errorf("flattening `key_vault_reference`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
