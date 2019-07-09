package azurerm

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmStorageManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageManagementPolicyCreateOrUpdate,
		Read:          resourceArmStorageManagementPolicyRead,
		Update:        resourceArmStorageManagementPolicyCreateOrUpdate,
		Delete:        resourceArmStorageManagementPolicyDelete,
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							validation.StringMatch(
								regexp.MustCompile(`^[a-zA-Z0-9]*$`),
								"A rule name can contain any combination of alpha numeric characters.",
							),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Lifecycle"}, false),
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
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
										ValidateFunc: validation.StringInSlice([]string{"blockBlob"}, false),
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
													Type:     schema.TypeInt,
													Optional: true,
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"delete_after_days_since_modification_greater_than": {
													Type:     schema.TypeInt,
													Optional: true,
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
													Type:     schema.TypeInt,
													Optional: true,
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

func resourceArmStorageManagementPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storage.ManagementPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	storageAccountId := d.Get("storage_account_id").(string)

	rid, err := parseAzureResourceID(storageAccountId)
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	name := "default" // The name of the Storage Account Management Policy. It should always be 'default' (from https://docs.microsoft.com/en-us/rest/api/storagerp/managementpolicies/createorupdate)

	parameters := storage.ManagementPolicy{
		Name: &name,
	}

	rules := d.Get("rule").([]interface{})
	armRules, err := expandStorageManagementPolicyRules(rules)

	parameters.ManagementPolicyProperties = &storage.ManagementPolicyProperties{
		Policy: &storage.ManagementPolicySchema{
			Rules: armRules,
		},
	}

	result, err := client.CreateOrUpdate(ctx, resourceGroupName, storageAccountName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Azure Storage Management Policy %q: %+v", storageAccountId, err)
	}

	d.SetId(*result.ID)

	return resourceArmStorageManagementPolicyRead(d, meta)
}

func resourceArmStorageManagementPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storage.ManagementPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	storageAccountId := d.Get("storage_account_id").(string)

	rid, err := parseAzureResourceID(storageAccountId)
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	result, err := client.Get(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}

	if result.Policy != nil {
		policy := result.Policy
		if policy.Rules != nil {
			if err := d.Set("rule", flattenStorageManagementPolicyRules(policy.Rules)); err != nil {
				return fmt.Errorf("Error flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmStorageManagementPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storage.ManagementPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	storageAccountId := d.Get("storage_account_id").(string)

	rid, err := parseAzureResourceID(storageAccountId)
	resourceGroupName := rid.ResourceGroup
	storageAccountName := rid.Path["storageAccounts"]

	_, err = client.Delete(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	return nil
}

func expandStorageManagementPolicyRules(list []interface{}) (*[]storage.ManagementPolicyRule, error) {
	result := []storage.ManagementPolicyRule{}

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		policyRule, err := expandStorageManagementPolicyRule(item)
		if err != nil {
			return nil, err
		}
		result = append(result, *policyRule)
	}
	return &result, nil
}

func expandStorageManagementPolicyRule(ref map[string]interface{}) (*storage.ManagementPolicyRule, error) {
	if len(ref) == 0 {
		return nil, fmt.Errorf("Error: storage management policy rule should be defined")
	}

	name := ref["name"].(string)
	enabled := ref["enabled"].(bool)
	typeVal := ref["type"].(string)

	definition := storage.ManagementPolicyDefinition{
		Filters: &storage.ManagementPolicyFilter{},
		Actions: &storage.ManagementPolicyAction{},
	}
	filtersRef := ref["filters"].([]interface{})
	if len(filtersRef) == 1 {
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
	actionsRef := ref["actions"].([]interface{})
	if len(actionsRef) == 1 {
		actionRef := actionsRef[0].(map[string]interface{})

		baseBlobsRef := actionRef["base_blob"].([]interface{})
		if len(baseBlobsRef) == 1 {
			baseBlob := &storage.ManagementPolicyBaseBlob{}
			baseBlobRef := baseBlobsRef[0].(map[string]interface{})
			if v, ok := baseBlobRef["tier_to_cool_after_days_since_modification_greater_than"]; ok {
				v2 := float64(v.(int))
				baseBlob.TierToCool = &storage.DateAfterModification{DaysAfterModificationGreaterThan: &v2}
			}
			if v, ok := baseBlobRef["tier_to_archive_after_days_since_modification_greater_than"]; ok {
				v2 := float64(v.(int))
				baseBlob.TierToArchive = &storage.DateAfterModification{DaysAfterModificationGreaterThan: &v2}
			}
			if v, ok := baseBlobRef["delete_after_days_since_modification_greater_than"]; ok {
				v2 := float64(v.(int))
				baseBlob.Delete = &storage.DateAfterModification{DaysAfterModificationGreaterThan: &v2}
			}
			definition.Actions.BaseBlob = baseBlob
		}

		snapshotRef := actionRef["snapshot"].([]interface{})
		if len(snapshotRef) == 1 {
			snapshot := &storage.ManagementPolicySnapShot{}
			snapshotRef := snapshotRef[0].(map[string]interface{})
			if v, ok := snapshotRef["delete_after_days_since_creation_greater_than"]; ok {
				v2 := float64(v.(int))
				snapshot.Delete = &storage.DateAfterCreation{DaysAfterCreationGreaterThan: &v2}
			}
			definition.Actions.Snapshot = snapshot
		}
	}

	rule := &storage.ManagementPolicyRule{
		Name:       &name,
		Enabled:    &enabled,
		Type:       &typeVal,
		Definition: &definition,
	}
	return rule, nil
}

func flattenStorageManagementPolicyRules(armRules *[]storage.ManagementPolicyRule) []interface{} {
	rules := make([]interface{}, 0)

	for _, armRule := range *armRules {
		rule := make(map[string]interface{})

		if armRule.Name != nil {
			rule["name"] = *armRule.Name
		}
		if armRule.Type != nil {
			rule["type"] = *armRule.Type
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
