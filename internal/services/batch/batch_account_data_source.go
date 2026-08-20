// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBatchAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBatchAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"allowed_authentication_modes": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"account_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"encryption": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

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

			"location": commonschema.LocationComputed(),

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_access": dataSourceBatchAccountEndpointAccessProfileSchema(),

						"node_management_access": dataSourceBatchAccountEndpointAccessProfileSchema(),
					},
				},
			},

			"pool_allocation_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"storage_account_authentication_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_node_identity": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceBatchAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("account_endpoint", props.AccountEndpoint)

			if err := d.Set("allowed_authentication_modes", flattenAllowedAuthenticationModes(props.AllowedAuthenticationModes)); err != nil {
				return fmt.Errorf("setting `allowed_authentication_modes`: %+v", err)
			}

			d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == batchaccount.PublicNetworkAccessTypeEnabled)

			if err := d.Set("network_profile", flattenBatchAccountNetworkProfile(props.NetworkProfile)); err != nil {
				return fmt.Errorf("setting `network_profile`: %+v", err)
			}

			if autoStorage := props.AutoStorage; autoStorage != nil {
				d.Set("storage_account_id", autoStorage.StorageAccountId)
				d.Set("storage_account_authentication_mode", string(pointer.From(autoStorage.AuthenticationMode)))

				if autoStorage.NodeIdentityReference != nil {
					d.Set("storage_account_node_identity", pointer.From(autoStorage.NodeIdentityReference.ResourceId))
				}
			} else {
				d.Set("storage_account_authentication_mode", "")
				d.Set("storage_account_id", "")
			}

			d.Set("pool_allocation_mode", string(pointer.From(props.PoolAllocationMode)))
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
				d.Set("key_vault_reference", []interface{}{})
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

func dataSourceBatchAccountEndpointAccessProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ip_rule": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"ip_range": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"action": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}
