// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	batchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Batch Account.",
		Read:        dataSourceBatchAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Batch Account.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"location":            commonschema.LocationComputed(),
			"storage_account_id": {
				Description: "The ID of the Storage Account used for this Batch Account.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"pool_allocation_mode": {
				Description: "The pool allocation mode configured for this Batch Account.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"key_vault_reference": {
				Description: "The `key_vault_reference` block that describes the Azure KeyVault reference to use when deploying the Azure Batch Account using the `UserSubscription` pool allocation mode.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Description: "The Azure identifier of the Azure KeyVault reference.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"url": {
							Description: "The HTTPS URL of the Azure KeyVault reference.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"primary_access_key": {
				Description: "The Batch Account primary access key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"secondary_access_key": {
				Description: "The Batch Account secondary access key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"account_endpoint": {
				Description: "The account endpoint used to interact with the Batch service.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"encryption": {
				Description: "The `encryption` block that describes the Azure KeyVault key reference used to encrypt data for the Azure Batch Account.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Description: "The full URL path of the Key Vault Key used to encrypt data for this Batch Account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceBatchAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := batchClient.From(meta).AccountClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := batchaccount.NewBatchAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("account_endpoint", props.AccountEndpoint)

			if autoStorage := props.AutoStorage; autoStorage != nil {
				d.Set("storage_account_id", autoStorage.StorageAccountId)
			}
			d.Set("pool_allocation_mode", pointer.FromEnum(props.PoolAllocationMode))
			poolAllocationMode := d.Get("pool_allocation_mode").(string)

			if encryption := props.Encryption; encryption != nil {
				d.Set("encryption", flattenEncryption(encryption))
			}

			if poolAllocationMode == string(batchaccount.PoolAllocationModeBatchService) {
				keys, err := client.GetKeys(ctx, id)
				if err != nil {
					return fmt.Errorf("cannot read keys for %s: %v", id, err)
				}

				if keysModel := keys.Model; keysModel != nil {
					d.Set("primary_access_key", keysModel.Primary)
					d.Set("secondary_access_key", keysModel.Secondary)
				}

				// set empty keyvault reference which is not needed in Batch Service allocation mode.
				d.Set("key_vault_reference", []any{})
			} else if poolAllocationMode == string(batchaccount.PoolAllocationModeUserSubscription) {
				if err := d.Set("key_vault_reference", flattenBatchAccountKeyvaultReference(props.KeyVaultReference)); err != nil {
					return fmt.Errorf("flattening `key_vault_reference`: %+v", err)
				}
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
