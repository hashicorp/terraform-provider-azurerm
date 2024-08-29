// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobinventorypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageBlobInventoryPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageBlobInventoryPolicyCreateUpdate,
		Read:   resourceStorageBlobInventoryPolicyRead,
		Update: resourceStorageBlobInventoryPolicyCreateUpdate,
		Delete: resourceStorageBlobInventoryPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.BlobInventoryPolicyV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.StorageAccountId{}),

			"rules": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"storage_container_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.StorageContainerName,
						},

						"format": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(blobinventorypolicies.FormatCsv),
								string(blobinventorypolicies.FormatParquet),
							}, false),
						},

						"schedule": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(blobinventorypolicies.ScheduleDaily),
								string(blobinventorypolicies.ScheduleWeekly),
							}, false),
						},

						"scope": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(blobinventorypolicies.ObjectTypeBlob),
								string(blobinventorypolicies.ObjectTypeContainer),
							}, false),
						},

						"schema_fields": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"filter": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"blob_types": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"blockBlob",
												"appendBlob",
												"pageBlob",
											}, false),
										},
									},

									"include_blob_versions": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"include_deleted": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"include_snapshots": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"prefix_match": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										MaxItems: 10,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"exclude_prefixes": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										MaxItems: 10,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
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

func resourceStorageBlobInventoryPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Storage.ResourceManager.BlobInventoryPolicies
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// This Resource is 1:1 with a Storage Account, therefore we use it's Resource ID for this Resource
	// however we want to ensure it's in the same subscription, so we'll build this up here
	id := commonids.NewStorageAccountID(subscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_blob_inventory_policy", id.ID())
		}
	}

	rules, err := expandBlobInventoryPolicyRules(d.Get("rules").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}

	payload := blobinventorypolicies.BlobInventoryPolicy{
		Properties: &blobinventorypolicies.BlobInventoryPolicyProperties{
			Policy: blobinventorypolicies.BlobInventoryPolicySchema{
				Enabled: true,
				Type:    blobinventorypolicies.InventoryRuleTypeInventory,
				Rules:   rules,
			},
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageBlobInventoryPolicyRead(d, meta)
}

func resourceStorageBlobInventoryPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.BlobInventoryPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Blob Inventory Policy for %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("storage_account_id", id.ID())
		if props := model.Properties; props != nil {
			if !props.Policy.Enabled {
				log.Printf("[INFO] Blob Inventory Policy is not enabled for %s - removing from state", *id)
				d.SetId("")
				return nil
			}

			d.Set("rules", flattenBlobInventoryPolicyRules(props.Policy.Rules))
		}
	}

	return nil
}

func resourceStorageBlobInventoryPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.BlobInventoryPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	return nil
}

func expandBlobInventoryPolicyRules(input []interface{}) ([]blobinventorypolicies.BlobInventoryPolicyRule, error) {
	results := make([]blobinventorypolicies.BlobInventoryPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		filters, err := expandBlobInventoryPolicyFilter(v["filter"].([]interface{}), v["scope"].(string))
		if err != nil {
			return nil, fmt.Errorf("%s rule is invalid: %+v", v["name"].(string), err)
		}

		results = append(results, blobinventorypolicies.BlobInventoryPolicyRule{
			Enabled:     true,
			Name:        v["name"].(string),
			Destination: v["storage_container_name"].(string),
			Definition: blobinventorypolicies.BlobInventoryPolicyDefinition{
				Format:       blobinventorypolicies.Format(v["format"].(string)),
				Schedule:     blobinventorypolicies.Schedule(v["schedule"].(string)),
				ObjectType:   blobinventorypolicies.ObjectType(v["scope"].(string)),
				SchemaFields: *utils.ExpandStringSlice(v["schema_fields"].([]interface{})),
				Filters:      filters,
			},
		})
	}
	return results, nil
}

func expandBlobInventoryPolicyFilter(input []interface{}, objectType string) (*blobinventorypolicies.BlobInventoryPolicyFilter, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})
	policyFilter := &blobinventorypolicies.BlobInventoryPolicyFilter{
		PrefixMatch:         utils.ExpandStringSlice(v["prefix_match"].(*pluginsdk.Set).List()),
		ExcludePrefix:       utils.ExpandStringSlice(v["exclude_prefixes"].(*pluginsdk.Set).List()),
		BlobTypes:           utils.ExpandStringSlice(v["blob_types"].(*pluginsdk.Set).List()),
		IncludeBlobVersions: utils.Bool(v["include_blob_versions"].(bool)),
		IncludeDeleted:      utils.Bool(v["include_deleted"].(bool)),
		IncludeSnapshots:    utils.Bool(v["include_snapshots"].(bool)),
	}

	// If the objectType is Container, the following values must be nil when passed to the API
	if objectType == string(blobinventorypolicies.ObjectTypeContainer) {
		if len(*policyFilter.BlobTypes) > 0 || *policyFilter.IncludeBlobVersions || *policyFilter.IncludeSnapshots {
			return nil, fmt.Errorf("`blobTypes`, `includeBlobVersions`, `includeSnapshots` cannot be used with objectType `Container`")
		}
		policyFilter.BlobTypes = nil
		policyFilter.IncludeBlobVersions = nil
		policyFilter.IncludeSnapshots = nil
	}

	return policyFilter, nil
}

func flattenBlobInventoryPolicyRules(input []blobinventorypolicies.BlobInventoryPolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		if !item.Enabled {
			continue
		}

		results = append(results, map[string]interface{}{
			"name":                   item.Name,
			"storage_container_name": item.Destination,
			"format":                 string(item.Definition.Format),
			"schedule":               string(item.Definition.Schedule),
			"scope":                  string(item.Definition.ObjectType),
			"schema_fields":          item.Definition.SchemaFields,
			"filter":                 flattenBlobInventoryPolicyFilter(item.Definition.Filters),
		})
	}
	return results
}

func flattenBlobInventoryPolicyFilter(input *blobinventorypolicies.BlobInventoryPolicyFilter) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var includeBlobVersions bool
	if input.IncludeBlobVersions != nil {
		includeBlobVersions = *input.IncludeBlobVersions
	}
	var includeDeleted bool
	if input.IncludeDeleted != nil {
		includeDeleted = *input.IncludeDeleted
	}
	var includeSnapshots bool
	if input.IncludeSnapshots != nil {
		includeSnapshots = *input.IncludeSnapshots
	}
	return []interface{}{
		map[string]interface{}{
			"blob_types":            utils.FlattenStringSlice(input.BlobTypes),
			"include_blob_versions": includeBlobVersions,
			"include_deleted":       includeDeleted,
			"include_snapshots":     includeSnapshots,
			"prefix_match":          utils.FlattenStringSlice(input.PrefixMatch),
			"exclude_prefixes":      utils.FlattenStringSlice(input.ExcludePrefix),
		},
	}
}
