package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"storage_container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageContainerName,
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

						"filter": {
							Type:     pluginsdk.TypeList,
							Required: true,
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

									"include_snapshots": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"prefix_match": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
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
	client := meta.(*clients.Client).Storage.BlobInventoryPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccount, err := parse.StorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewBlobInventoryPolicyID(subscriptionId, storageAccount.ResourceGroup, storageAccount.Name, "Default")

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
				Enabled:     utils.Bool(true),
				Destination: utils.String(d.Get("storage_container_name").(string)),
				Type:        utils.String("Inventory"),
				Rules:       expandBlobInventoryPolicyRules(d.Get("rules").(*pluginsdk.Set).List()),
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
	d.Set("storage_account_id", parse.NewStorageAccountID(subscriptionId, id.ResourceGroup, id.StorageAccountName).ID())
	if props := resp.BlobInventoryPolicyProperties; props != nil {
		if policy := props.Policy; policy != nil {
			if policy.Enabled == nil || !*policy.Enabled {
				log.Printf("[INFO] storage %q is not enabled - removing from state", d.Id())
				d.SetId("")
				return nil
			}
			d.Set("storage_container_name", policy.Destination)
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
			Enabled: utils.Bool(true),
			Name:    utils.String(v["name"].(string)),
			Definition: &storage.BlobInventoryPolicyDefinition{
				Filters: expandBlobInventoryPolicyFilter(v["filter"].([]interface{})),
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
		BlobTypes:           utils.ExpandStringSlice(v["blob_types"].(*pluginsdk.Set).List()),
		IncludeBlobVersions: utils.Bool(v["include_blob_versions"].(bool)),
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
		if item.Enabled == nil || !*item.Enabled || item.Definition == nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"name":   name,
			"filter": flattenBlobInventoryPolicyFilter(item.Definition.Filters),
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
	var includeSnapshots bool
	if input.IncludeSnapshots != nil {
		includeSnapshots = *input.IncludeSnapshots
	}
	return []interface{}{
		map[string]interface{}{
			"blob_types":            utils.FlattenStringSlice(input.BlobTypes),
			"include_blob_versions": includeBlobVersions,
			"include_snapshots":     includeSnapshots,
			"prefix_match":          utils.FlattenStringSlice(input.PrefixMatch),
		},
	}
}
