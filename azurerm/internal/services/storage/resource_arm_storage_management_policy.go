package storage

import (
	"fmt"
	"regexp"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceArmStorageManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageManagementPolicyCreateOrUpdate,
		Read:   resourceArmStorageManagementPolicyRead,
		Update: resourceArmStorageManagementPolicyCreateOrUpdate,
		Delete: resourceArmStorageManagementPolicyDelete,
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
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"blockBlob"}, false),
										},
										Set: schema.HashString,
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
													Default:      nil,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      nil,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"delete_after_days_since_modification_greater_than": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      nil,
													ValidateFunc: validation.IntAtLeast(0),
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
													ValidateFunc: validation.IntAtLeast(0),
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

func resourceArmStorageManagementPolicyCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
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

	parameters := storage.ManagementPolicy{
		Name: &name,
	}

	armRules, err := expandStorageManagementPolicyRules(d)
	if err != nil {
		return fmt.Errorf("Error expanding Azure Storage Management Policy Rules %q: %+v", storageAccountId, err)
	}

	parameters.ManagementPolicyProperties = &storage.ManagementPolicyProperties{
		Policy: &storage.ManagementPolicySchema{
			Rules: armRules,
		},
	}

	result, err := client.CreateOrUpdate(ctx, resourceGroupName, storageAccountName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Azure Storage Management Policy %q: %+v", storageAccountId, err)
	}

	result, err = client.Get(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return fmt.Errorf("Error getting created Azure Storage Management Policy %q: %+v", storageAccountId, err)
	}

	d.SetId(*result.ID)

	return resourceArmStorageManagementPolicyRead(d, meta)
}

func resourceArmStorageManagementPolicyRead(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("Error flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmStorageManagementPolicyDelete(d *schema.ResourceData, meta interface{}) error {
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

	_, err = client.Delete(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	return nil
}

// nolint unparam
func expandStorageManagementPolicyRules(d *schema.ResourceData) (*[]storage.ManagementPolicyRule, error) {
	var result []storage.ManagementPolicyRule

	rules := d.Get("rule").([]interface{})

	for k, v := range rules {
		if v != nil {
			result = append(result, expandStorageManagementPolicyRule(d, k))
		}
	}
	return &result, nil
}

func expandStorageManagementPolicyRule(d *schema.ResourceData, ruleIndex int) storage.ManagementPolicyRule {
	name := d.Get(fmt.Sprintf("rule.%d.name", ruleIndex)).(string)
	enabled := d.Get(fmt.Sprintf("rule.%d.enabled", ruleIndex)).(bool)
	typeVal := "Lifecycle"

	definition := storage.ManagementPolicyDefinition{
		Filters: &storage.ManagementPolicyFilter{},
		Actions: &storage.ManagementPolicyAction{},
	}
	filtersRef := d.Get(fmt.Sprintf("rule.%d.filters", ruleIndex)).([]interface{})
	if len(filtersRef) == 1 {
		if filtersRef[0] != nil {
			filterRef := filtersRef[0].(map[string]interface{})

			prefixMatches := []string{}
			prefixMatchesRef := filterRef["prefix_match"].(*schema.Set)
			if prefixMatchesRef != nil {
				for _, prefixMatchRef := range prefixMatchesRef.List() {
					prefixMatches = append(prefixMatches, prefixMatchRef.(string))
				}
			}
			definition.Filters.PrefixMatch = &prefixMatches

			blobTypes := []string{}
			blobTypesRef := filterRef["blob_types"].(*schema.Set)
			if blobTypesRef != nil {
				for _, blobTypeRef := range blobTypesRef.List() {
					blobTypes = append(blobTypes, blobTypeRef.(string))
				}
			}
			definition.Filters.BlobTypes = &blobTypes
		}
	}
	if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions", ruleIndex)); ok {
		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob", ruleIndex)); ok {
			baseBlob := &storage.ManagementPolicyBaseBlob{}
			if v, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", ruleIndex)); ok {
				if v != nil {
					baseBlob.TierToCool = &storage.DateAfterModification{
						DaysAfterModificationGreaterThan: utils.Float(float64(v.(int))),
					}
				}
			}
			if v, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", ruleIndex)); ok {
				if v != nil {
					baseBlob.TierToArchive = &storage.DateAfterModification{
						DaysAfterModificationGreaterThan: utils.Float(float64(v.(int))),
					}
				}
			}
			if v, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", ruleIndex)); ok {
				if v != nil {
					baseBlob.Delete = &storage.DateAfterModification{
						DaysAfterModificationGreaterThan: utils.Float(float64(v.(int))),
					}
				}
			}
			definition.Actions.BaseBlob = baseBlob
		}

		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.snapshot", ruleIndex)); ok {
			snapshot := &storage.ManagementPolicySnapShot{}
			if v, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", ruleIndex)); ok {
				v2 := float64(v.(int))
				snapshot.Delete = &storage.DateAfterCreation{DaysAfterCreationGreaterThan: &v2}
			}
			definition.Actions.Snapshot = snapshot
		}
	}

	rule := storage.ManagementPolicyRule{
		Name:       &name,
		Enabled:    &enabled,
		Type:       &typeVal,
		Definition: &definition,
	}
	return rule
}

func flattenStorageManagementPolicyRules(armRules *[]storage.ManagementPolicyRule) []interface{} {
	rules := make([]interface{}, 0)
	if armRules == nil {
		return rules
	}
	for _, armRule := range *armRules {
		rule := make(map[string]interface{})

		if armRule.Name != nil {
			rule["name"] = *armRule.Name
		}
		if armRule.Enabled != nil {
			rule["enabled"] = *armRule.Enabled
		}

		armDefinition := armRule.Definition
		if armDefinition != nil {
			armFilter := armDefinition.Filters
			if armFilter != nil {
				filter := make(map[string]interface{})
				if armFilter.PrefixMatch != nil {
					prefixMatches := make([]interface{}, 0)
					for _, armPrefixMatch := range *armFilter.PrefixMatch {
						prefixMatches = append(prefixMatches, armPrefixMatch)
					}
					filter["prefix_match"] = prefixMatches
				}
				if armFilter.BlobTypes != nil {
					blobTypes := make([]interface{}, 0)
					for _, armBlobType := range *armFilter.BlobTypes {
						blobTypes = append(blobTypes, armBlobType)
					}
					filter["blob_types"] = blobTypes
				}
				rule["filters"] = [1]interface{}{filter}
			}

			armAction := armDefinition.Actions
			if armAction != nil {
				action := make(map[string]interface{})
				armActionBaseBlob := armAction.BaseBlob
				if armActionBaseBlob != nil {
					baseBlob := make(map[string]interface{})
					if armActionBaseBlob.TierToCool != nil && armActionBaseBlob.TierToCool.DaysAfterModificationGreaterThan != nil {
						intTemp := int(*armActionBaseBlob.TierToCool.DaysAfterModificationGreaterThan)
						baseBlob["tier_to_cool_after_days_since_modification_greater_than"] = intTemp
					}
					if armActionBaseBlob.TierToArchive != nil && armActionBaseBlob.TierToArchive.DaysAfterModificationGreaterThan != nil {
						intTemp := int(*armActionBaseBlob.TierToArchive.DaysAfterModificationGreaterThan)
						baseBlob["tier_to_archive_after_days_since_modification_greater_than"] = intTemp
					}
					if armActionBaseBlob.Delete != nil && armActionBaseBlob.Delete.DaysAfterModificationGreaterThan != nil {
						intTemp := int(*armActionBaseBlob.Delete.DaysAfterModificationGreaterThan)
						baseBlob["delete_after_days_since_modification_greater_than"] = intTemp
					}
					action["base_blob"] = [1]interface{}{baseBlob}
				}

				armActionSnaphost := armAction.Snapshot
				if armActionSnaphost != nil {
					snapshot := make(map[string]interface{})
					if armActionSnaphost.Delete != nil && armActionSnaphost.Delete.DaysAfterCreationGreaterThan != nil {
						intTemp := int(*armActionSnaphost.Delete.DaysAfterCreationGreaterThan)
						snapshot["delete_after_days_since_creation_greater_than"] = intTemp
					}
					action["snapshot"] = [1]interface{}{snapshot}
				}

				rule["actions"] = [1]interface{}{action}
			}
		}

		rules = append(rules, rule)
	}

	return rules
}
