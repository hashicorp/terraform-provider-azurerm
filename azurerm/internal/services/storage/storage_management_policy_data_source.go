package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceStorageManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStorageManagementPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"storage_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix_match": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									"blob_types": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"match_blob_index_tag": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},

												"operation": {
													Type:     schema.TypeString,
													Computed: true,
												},

												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"base_blob": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tier_to_cool_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"snapshot": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"change_tier_to_archive_after_days_since_creation": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"change_tier_to_cool_after_days_since_creation": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"delete_after_days_since_creation_greater_than": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"version": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"change_tier_to_archive_after_days_since_creation": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"change_tier_to_cool_after_days_since_creation": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"delete_after_days_since_creation": {
													Type:     schema.TypeInt,
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

func dataSourceStorageManagementPolicyRead(d *schema.ResourceData, meta interface{}) error {
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
