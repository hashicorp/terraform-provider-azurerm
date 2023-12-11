// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceStorageManagementPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageManagementPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"rule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"filters": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"prefix_match": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},
									"blob_types": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"match_blob_index_tag": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},

												"operation": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},

												"value": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"actions": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"base_blob": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"tier_to_cool_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_cool_after_days_since_last_access_time_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"auto_tier_to_hot_from_cool_enabled": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
												},
												"tier_to_cool_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_last_access_time_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_last_tier_change_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_cold_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_cold_after_days_since_last_access_time_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_cold_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_last_access_time_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"snapshot": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"change_tier_to_archive_after_days_since_creation": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_last_tier_change_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"change_tier_to_cool_after_days_since_creation": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_cold_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"version": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"change_tier_to_archive_after_days_since_creation": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_last_tier_change_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"change_tier_to_cool_after_days_since_creation": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"tier_to_cold_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_creation": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceStorageManagementPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ManagementPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewStorageAccountManagementPolicyID(storageAccountId.SubscriptionId, storageAccountId.ResourceGroupName, storageAccountId.StorageAccountName, "default")
	resp, err := client.Get(ctx, *storageAccountId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("rule", flattenStorageManagementPolicyRules(props.Policy.Rules)); err != nil {
				return fmt.Errorf("flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}
