package storage

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_modification_greater_than": {
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
												"change_tier_to_cool_after_days_since_creation": {
													Type:     pluginsdk.TypeInt,
													Optional: true,
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
												"change_tier_to_cool_after_days_since_creation": {
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
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountId := d.Get("storage_account_id").(string)

	rid, err := azure.ParseAzureResourceID(storageAccountId)
	if err != nil {
		return err
	}
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	result, err := client.Get(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	d.SetId(*result.ID)

	if props := result.ManagementPolicyProperties; props != nil {
		if policy := props.Policy; policy != nil {
			if err := d.Set("rule", flattenStorageManagementPolicyRules(policy.Rules)); err != nil {
				return fmt.Errorf("flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}
