package storage

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceStorageManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageManagementPolicyCreateOrUpdate,
		Read:   resourceStorageManagementPolicyRead,
		Update: resourceStorageManagementPolicyCreateOrUpdate,
		Delete: resourceStorageManagementPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[a-zA-Z0-9]*$`),
								"A rule name can contain any combination of alpha numeric characters.",
							),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"filters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix_match": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									"blob_types": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"blockBlob",
												"appendBlob",
											}, false),
										},
										Set: schema.HashString,
									},

									"blob_index_match_tag": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.StorageBlobIndexTagName,
												},

												"operation": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"==",
													}, false),
													Default: "==",
												},

												"value": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.StorageBlobIndexTagValue,
												},
											},
										},
									},
								},
							},
						},
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"base_blob": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tier_to_cool_after_days_since_modification_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"delete_after_days_since_modification_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
											},
										},
									},
									"snapshot": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"delete_after_days_since_creation_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_archive_after_days_since_creation_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_cool_after_days_since_creation_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
											},
										},
									},
									"version": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"delete_after_days_since_creation_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_archive_after_days_since_creation_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_cool_after_days_since_creation_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
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

func resourceStorageManagementPolicyCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountId := d.Get("storage_account_id").(string)

	rid, err := azure.ParseAzureResourceID(storageAccountId)
	if err != nil {
		return err
	}
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	name := "default" // The name of the Storage Account Management Policy. It should always be 'default' (from https://docs.microsoft.com/en-us/rest/api/storagerp/managementpolicies/createorupdate)

	rules, err := expandStorageManagementPolicyRules(d.Get("rule").([]interface{}))
	if err != nil {
		return err
	}

	parameters := storage.ManagementPolicy{
		Name: &name,
		ManagementPolicyProperties: &storage.ManagementPolicyProperties{
			Policy: &storage.ManagementPolicySchema{
				Rules: rules,
			},
		},
	}

	result, err := client.CreateOrUpdate(ctx, resourceGroupName, storageAccountName, parameters)
	if err != nil {
		return fmt.Errorf("creating Azure Storage Management Policy %q: %+v", storageAccountId, err)
	}

	result, err = client.Get(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return fmt.Errorf("getting created Azure Storage Management Policy %q: %+v", storageAccountId, err)
	}

	d.SetId(*result.ID)

	return resourceStorageManagementPolicyRead(d, meta)
}

func resourceStorageManagementPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()

	rid, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return err
	}
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	result, err := client.Get(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}

	// TODO: switch this to look up the account and use that, rather than building this up
	storageAccountID := "/subscriptions/" + rid.SubscriptionID + "/resourceGroups/" + rid.ResourceGroup + "/providers/" + rid.Provider + "/storageAccounts/" + storageAccountName
	d.Set("storage_account_id", storageAccountID)

	if policy := result.Policy; policy != nil {
		policy := result.Policy
		if rules := policy.Rules; rules != nil {
			if err := d.Set("rule", flattenStorageManagementPolicyRules(policy.Rules)); err != nil {
				return fmt.Errorf("flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}

func resourceStorageManagementPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()

	rid, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return err
	}
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	if _, err = client.Delete(ctx, resourceGroupName, storageAccountName); err != nil {
		return err
	}
	return nil
}

// nolint unparam
func expandStorageManagementPolicyRules(inputs []interface{}) (*[]storage.ManagementPolicyRule, error) {
	result := make([]storage.ManagementPolicyRule, 0)
	if len(inputs) == 0 {
		return &result, nil
	}

	for _, input := range inputs {
		v := input.(map[string]interface{})
		rule := storage.ManagementPolicyRule{
			Name:    utils.String(v["name"].(string)),
			Enabled: utils.Bool(v["enabled"].(bool)),
			Type:    utils.String("Lifecycle"),
			Definition: &storage.ManagementPolicyDefinition{
				Actions: expandStorageManagementPolicyActions(v["actions"].([]interface{})),
				Filters: expandStorageManagementPolicyFilters(v["filters"].([]interface{})),
			},
		}
		if (rule.Definition.Actions.Version != nil || rule.Definition.Actions.Snapshot != nil) && rule.Definition.Filters.BlobIndexMatch != nil {
			return nil, fmt.Errorf("`blob_index_match_tag` is not supported as a filter for versions and snapshots")
		}
		result = append(result, rule)
	}
	return &result, nil
}

func expandStorageManagementPolicyFilters(inputs []interface{}) *storage.ManagementPolicyFilter {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	return &storage.ManagementPolicyFilter{
		PrefixMatch:    utils.ExpandStringSlice(input["prefix_match"].(*schema.Set).List()),
		BlobTypes:      utils.ExpandStringSlice(input["blob_types"].(*schema.Set).List()),
		BlobIndexMatch: expandAzureRmStorageBlobIndexMatch(input["blob_index_match_tag"].(*schema.Set).List()),
	}
}

func expandStorageManagementPolicyActions(inputs []interface{}) *storage.ManagementPolicyAction {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	return &storage.ManagementPolicyAction{
		BaseBlob: expandStorageManagementPolicyActionsBaseBlob(input["base_blob"].([]interface{})),
		Snapshot: expandStorageManagementPolicyActionsSnapshot(input["snapshot"].([]interface{})),
		Version:  expandStorageManagementPolicyActionsVersion(input["version"].([]interface{})),
	}
}

func expandStorageManagementPolicyActionsBaseBlob(inputs []interface{}) *storage.ManagementPolicyBaseBlob {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})
	result := storage.ManagementPolicyBaseBlob{}

	if v, ok := input["delete_after_days_since_modification_greater_than"].(int); ok && v != -1 {
		result.Delete = &storage.DateAfterModification{
			DaysAfterModificationGreaterThan: utils.Float(float64(v)),
		}
	}
	if v, ok := input["tier_to_archive_after_days_since_modification_greater_than"].(int); ok && v != -1 {
		result.TierToArchive = &storage.DateAfterModification{
			DaysAfterModificationGreaterThan: utils.Float(float64(v)),
		}
	}
	if v, ok := input["tier_to_cool_after_days_since_modification_greater_than"].(int); ok && v != -1 {
		result.TierToCool = &storage.DateAfterModification{
			DaysAfterModificationGreaterThan: utils.Float(float64(v)),
		}
	}

	return &result
}

func expandStorageManagementPolicyActionsSnapshot(inputs []interface{}) *storage.ManagementPolicySnapShot {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	result := storage.ManagementPolicySnapShot{}

	if v, ok := input["delete_after_days_since_creation_greater_than"].(int); ok && v != -1 {
		result.Delete = &storage.DateAfterCreation{
			DaysAfterCreationGreaterThan: utils.Float(float64(v)),
		}
	}
	if v, ok := input["tier_to_archive_after_days_since_creation_greater_than"].(int); ok && v != -1 {
		result.TierToArchive = &storage.DateAfterCreation{
			DaysAfterCreationGreaterThan: utils.Float(float64(v)),
		}
	}
	if v, ok := input["tier_to_cool_after_days_since_creation_greater_than"].(int); ok && v != -1 {
		result.TierToCool = &storage.DateAfterCreation{
			DaysAfterCreationGreaterThan: utils.Float(float64(v)),
		}
	}

	return &result
}

func expandStorageManagementPolicyActionsVersion(inputs []interface{}) *storage.ManagementPolicyVersion {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	result := storage.ManagementPolicyVersion{}

	if v, ok := input["delete_after_days_since_creation_greater_than"].(int); ok && v != -1 {
		result.Delete = &storage.DateAfterCreation{
			DaysAfterCreationGreaterThan: utils.Float(float64(v)),
		}
	}
	if v, ok := input["tier_to_archive_after_days_since_creation_greater_than"].(int); ok && v != -1 {
		result.TierToArchive = &storage.DateAfterCreation{
			DaysAfterCreationGreaterThan: utils.Float(float64(v)),
		}
	}
	if v, ok := input["tier_to_cool_after_days_since_creation_greater_than"].(int); ok && v != -1 {
		result.TierToCool = &storage.DateAfterCreation{
			DaysAfterCreationGreaterThan: utils.Float(float64(v)),
		}
	}

	return &result
}

func flattenStorageManagementPolicyRules(armRules *[]storage.ManagementPolicyRule) []interface{} {
	rules := make([]interface{}, 0)
	if armRules == nil || len(*armRules) == 0 {
		return rules
	}

	for _, armRule := range *armRules {
		var name string
		if armRule.Name != nil {
			name = *armRule.Name
		}
		var enabled bool
		if armRule.Enabled != nil {
			enabled = *armRule.Enabled
		}

		var filters, actions []interface{}
		if armRule.Definition != nil {
			filters = flattenStorageManagementPolicyFilters(armRule.Definition.Filters)
			actions = flattenStorageManagementPolicyActions(armRule.Definition.Actions)
		}

		rules = append(rules, map[string]interface{}{
			"name":    name,
			"enabled": enabled,
			"filters": filters,
			"actions": actions,
		})
	}
	return rules
}

func flattenStorageManagementPolicyFilters(filters *storage.ManagementPolicyFilter) []interface{} {
	if filters == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"prefix_match":         utils.FlattenStringSlice(filters.PrefixMatch),
			"blob_types":           utils.FlattenStringSlice(filters.BlobTypes),
			"blob_index_match_tag": flattenAzureRmStorageBlobIndexMatch(filters.BlobIndexMatch),
		},
	}
}

func flattenStorageManagementPolicyActions(actions *storage.ManagementPolicyAction) []interface{} {
	if actions == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"base_blob": flattenStorageManagementPolicyActionsBaseBlob(actions.BaseBlob),
			"snapshot":  flattenStorageManagementPolicyActionsSnapshot(actions.Snapshot),
			"version":   flattenStorageManagementPolicyActionsVersion(actions.Version),
		},
	}
}

func flattenStorageManagementPolicyActionsBaseBlob(baseBlob *storage.ManagementPolicyBaseBlob) []interface{} {
	if baseBlob == nil {
		return []interface{}{}
	}

	deleteModification, archiveModification, coolModification := -1, -1, -1
	if v := baseBlob.Delete; v != nil {
		if v.DaysAfterModificationGreaterThan != nil {
			deleteModification = int(*v.DaysAfterModificationGreaterThan)
		}
	}
	if v := baseBlob.TierToArchive; v != nil {
		if v.DaysAfterModificationGreaterThan != nil {
			archiveModification = int(*v.DaysAfterModificationGreaterThan)
		}
	}
	if v := baseBlob.TierToCool; v != nil {
		if v.DaysAfterModificationGreaterThan != nil {
			coolModification = int(*v.DaysAfterModificationGreaterThan)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"delete_after_days_since_modification_greater_than":          deleteModification,
			"tier_to_archive_after_days_since_modification_greater_than": archiveModification,
			"tier_to_cool_after_days_since_modification_greater_than":    coolModification,
		},
	}
}

func flattenStorageManagementPolicyActionsSnapshot(snapshot *storage.ManagementPolicySnapShot) []interface{} {
	if snapshot == nil {
		return []interface{}{}
	}

	deleteCreation, archiveCreation, coolCreation := -1, -1, -1
	if v := snapshot.Delete; v != nil {
		if v.DaysAfterCreationGreaterThan != nil {
			deleteCreation = int(*v.DaysAfterCreationGreaterThan)
		}
	}
	if v := snapshot.TierToArchive; v != nil {
		if v.DaysAfterCreationGreaterThan != nil {
			archiveCreation = int(*v.DaysAfterCreationGreaterThan)
		}
	}
	if v := snapshot.TierToCool; v != nil {
		if v.DaysAfterCreationGreaterThan != nil {
			coolCreation = int(*v.DaysAfterCreationGreaterThan)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"delete_after_days_since_creation_greater_than":          deleteCreation,
			"tier_to_archive_after_days_since_creation_greater_than": archiveCreation,
			"tier_to_cool_after_days_since_creation_greater_than":    coolCreation,
		},
	}
}

func flattenStorageManagementPolicyActionsVersion(version *storage.ManagementPolicyVersion) []interface{} {
	if version == nil {
		return []interface{}{}
	}

	deleteCreation, archiveCreation, coolCreation := -1, -1, -1
	if v := version.Delete; v != nil {
		if v.DaysAfterCreationGreaterThan != nil {
			deleteCreation = int(*v.DaysAfterCreationGreaterThan)
		}
	}
	if v := version.TierToArchive; v != nil {
		if v.DaysAfterCreationGreaterThan != nil {
			archiveCreation = int(*v.DaysAfterCreationGreaterThan)
		}
	}
	if v := version.TierToCool; v != nil {
		if v.DaysAfterCreationGreaterThan != nil {
			coolCreation = int(*v.DaysAfterCreationGreaterThan)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"delete_after_days_since_creation_greater_than":          deleteCreation,
			"tier_to_archive_after_days_since_creation_greater_than": archiveCreation,
			"tier_to_cool_after_days_since_creation_greater_than":    coolCreation,
		},
	}
}

func expandAzureRmStorageBlobIndexMatch(blobIndexMatches []interface{}) *[]storage.TagFilter {
	if len(blobIndexMatches) == 0 {
		return nil
	}

	results := make([]storage.TagFilter, 0)
	for _, v := range blobIndexMatches {
		blobIndexMatch := v.(map[string]interface{})

		filter := storage.TagFilter{
			Name:  utils.String(blobIndexMatch["name"].(string)),
			Op:    utils.String(blobIndexMatch["operation"].(string)),
			Value: utils.String(blobIndexMatch["value"].(string)),
		}

		results = append(results, filter)
	}

	return &results
}

func flattenAzureRmStorageBlobIndexMatch(blobIndexMatches *[]storage.TagFilter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if blobIndexMatches == nil || len(*blobIndexMatches) == 0 {
		return result
	}

	for _, blobIndexMatch := range *blobIndexMatches {
		var name, op, value string
		if blobIndexMatch.Name != nil {
			name = *blobIndexMatch.Name
		}
		if blobIndexMatch.Op != nil {
			op = *blobIndexMatch.Op
		}
		if blobIndexMatch.Value != nil {
			value = *blobIndexMatch.Value
		}
		result = append(result, map[string]interface{}{
			"name":      name,
			"operation": op,
			"value":     value,
		})
	}
	return result
}
