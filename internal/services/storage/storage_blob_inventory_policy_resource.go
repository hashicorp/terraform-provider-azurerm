// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
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
			_, err := parse.BlobInventoryPolicyID(id)
			return err
		}),

		Schema: storageBlobInventoryPolicyResourceSchema(),
		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			rules := diff.Get("rules").(*pluginsdk.Set).List()
			for _, rule := range rules {
				v := rule.(map[string]interface{})
				if v["scope"] != string(storage.ObjectTypeBlob) && len(v["filter"].([]interface{})) != 0 {
					return fmt.Errorf("the `filter` can only be set when the `scope` is `%s`", storage.ObjectTypeBlob)
				}
			}

			return nil
		}),
	}
}

func storageBlobInventoryPolicyResourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

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
							string(storage.FormatCsv),
							string(storage.FormatParquet),
						}, false),
					},

					"schedule": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(storage.ScheduleDaily),
							string(storage.ScheduleWeekly),
						}, false),
					},

					"scope": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(storage.ObjectTypeBlob),
							string(storage.ObjectTypeContainer),
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
	}
}

func resourceStorageBlobInventoryPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Storage.BlobInventoryPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccount, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewBlobInventoryPolicyID(subscriptionId, storageAccount.ResourceGroupName, storageAccount.StorageAccountName, "Default")

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_storage_blob_inventory_policy", id.ID())
		}
	}

	props := storage.BlobInventoryPolicy{
		BlobInventoryPolicyProperties: &storage.BlobInventoryPolicyProperties{
			Policy: &storage.BlobInventoryPolicySchema{
				Enabled: utils.Bool(true),
				Type:    utils.String("Inventory"),
				Rules:   expandBlobInventoryPolicyRules(d.Get("rules").(*pluginsdk.Set).List()),
			},
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, props); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageBlobInventoryPolicyRead(d, meta)
}

func resourceStorageBlobInventoryPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Storage.BlobInventoryPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BlobInventoryPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] storage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("storage_account_id", commonids.NewStorageAccountID(subscriptionId, id.ResourceGroup, id.StorageAccountName).ID())
	if props := resp.BlobInventoryPolicyProperties; props != nil {
		if policy := props.Policy; policy != nil {
			if policy.Enabled == nil || !*policy.Enabled {
				log.Printf("[INFO] storage %q is not enabled - removing from state", d.Id())
				d.SetId("")
				return nil
			}

			d.Set("rules", flattenBlobInventoryPolicyRules(policy.Rules))
		}
	}
	return nil
}

func resourceStorageBlobInventoryPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.BlobInventoryPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BlobInventoryPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.StorageAccountName); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	return nil
}

func expandBlobInventoryPolicyRules(input []interface{}) *[]storage.BlobInventoryPolicyRule {
	results := make([]storage.BlobInventoryPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, storage.BlobInventoryPolicyRule{
			Enabled:     utils.Bool(true),
			Name:        utils.String(v["name"].(string)),
			Destination: utils.String(v["storage_container_name"].(string)),
			Definition: &storage.BlobInventoryPolicyDefinition{
				Format:       storage.Format(v["format"].(string)),
				Schedule:     storage.Schedule(v["schedule"].(string)),
				ObjectType:   storage.ObjectType(v["scope"].(string)),
				SchemaFields: utils.ExpandStringSlice(v["schema_fields"].([]interface{})),
				Filters:      expandBlobInventoryPolicyFilter(v["filter"].([]interface{})),
			},
		})
	}
	return &results
}

func expandBlobInventoryPolicyFilter(input []interface{}) *storage.BlobInventoryPolicyFilter {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storage.BlobInventoryPolicyFilter{
		PrefixMatch:         utils.ExpandStringSlice(v["prefix_match"].(*pluginsdk.Set).List()),
		ExcludePrefix:       utils.ExpandStringSlice(v["exclude_prefixes"].(*pluginsdk.Set).List()),
		BlobTypes:           utils.ExpandStringSlice(v["blob_types"].(*pluginsdk.Set).List()),
		IncludeBlobVersions: utils.Bool(v["include_blob_versions"].(bool)),
		IncludeDeleted:      utils.Bool(v["include_deleted"].(bool)),
		IncludeSnapshots:    utils.Bool(v["include_snapshots"].(bool)),
	}
}

func flattenBlobInventoryPolicyRules(input *[]storage.BlobInventoryPolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var destination string
		if item.Destination != nil {
			destination = *item.Destination
		}

		if item.Enabled == nil || !*item.Enabled || item.Definition == nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"name":                   name,
			"storage_container_name": destination,
			"format":                 string(item.Definition.Format),
			"schedule":               string(item.Definition.Schedule),
			"scope":                  string(item.Definition.ObjectType),
			"schema_fields":          utils.FlattenStringSlice(item.Definition.SchemaFields),
			"filter":                 flattenBlobInventoryPolicyFilter(item.Definition.Filters),
		})
	}
	return results
}

func flattenBlobInventoryPolicyFilter(input *storage.BlobInventoryPolicyFilter) []interface{} {
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
