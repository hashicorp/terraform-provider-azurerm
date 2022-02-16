package storage

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageManagementPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageManagementPolicyCreateOrUpdate,
		Read:   resourceStorageManagementPolicyRead,
		Update: resourceStorageManagementPolicyCreateOrUpdate,
		Delete: resourceStorageManagementPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageAccountManagementPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"rule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[a-zA-Z0-9-]*$`),
								"A rule name can contain any combination of alpha numeric characters.",
							),
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						//lintignore:XS003
						"filters": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"prefix_match": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},
									"blob_types": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"blockBlob",
												"appendBlob",
											}, false),
										},
										Set: pluginsdk.HashString,
									},

									"match_blob_index_tag": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validate.StorageBlobIndexTagName,
												},

												"operation": {
													Type:     pluginsdk.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"==",
													}, false),
													Default: "==",
												},

												"value": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validate.StorageBlobIndexTagValue,
												},
											},
										},
									},
								},
							},
						},
						//lintignore:XS003
						"actions": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									//lintignore:XS003
									"base_blob": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"tier_to_cool_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Optional: true,
													Default:  nil,
													// todo: default change to -1 to allow value 0 in 3.0
													// for issue https://github.com/hashicorp/terraform-provider-azurerm/issues/6158
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_cool_after_days_since_last_access_time_greater_than": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_archive_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Optional: true,
													Default:  nil,
													// todo: default change to -1 to allow value 0 in 3.0
													// for issue https://github.com/hashicorp/terraform-provider-azurerm/issues/6158
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"tier_to_archive_after_days_since_last_access_time_greater_than": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"delete_after_days_since_modification_greater_than": {
													Type:     pluginsdk.TypeInt,
													Optional: true,
													Default:  nil,
													// todo: default change to -1 to allow value 0 in 3.0
													// for issue https://github.com/hashicorp/terraform-provider-azurerm/issues/6158
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"delete_after_days_since_last_access_time_greater_than": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
											},
										},
									},
									//lintignore:XS003
									"snapshot": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"change_tier_to_archive_after_days_since_creation": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"change_tier_to_cool_after_days_since_creation": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"delete_after_days_since_creation_greater_than": {
													Type:     pluginsdk.TypeInt,
													Optional: true,
													// todo: default change to -1 to allow value 0 in 3.0
													// for issue https://github.com/hashicorp/terraform-provider-azurerm/issues/6158
													ValidateFunc: validation.IntBetween(0, 99999),
												},
											},
										},
									},
									"version": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"change_tier_to_archive_after_days_since_creation": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"change_tier_to_cool_after_days_since_creation": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 99999),
												},
												"delete_after_days_since_creation": {
													Type:         pluginsdk.TypeInt,
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

func resourceStorageManagementPolicyCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rid, err := parse.StorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// The name of the Storage Account Management Policy. It should always be 'default' (from https://docs.microsoft.com/en-us/rest/api/storagerp/managementpolicies/createorupdate)
	mgmtPolicyId := parse.NewStorageAccountManagementPolicyID(rid.SubscriptionId, rid.ResourceGroup, rid.Name, "default")

	parameters := storage.ManagementPolicy{
		Name: &mgmtPolicyId.ManagementPolicyName,
	}

	armRules, err := expandStorageManagementPolicyRules(d)
	if err != nil {
		return fmt.Errorf("expanding %s: %+v", mgmtPolicyId, err)
	}

	parameters.ManagementPolicyProperties = &storage.ManagementPolicyProperties{
		Policy: &storage.ManagementPolicySchema{
			Rules: armRules,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, rid.ResourceGroup, rid.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", mgmtPolicyId, err)
	}

	d.SetId(mgmtPolicyId.ID())

	return resourceStorageManagementPolicyRead(d, meta)
}

func resourceStorageManagementPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rid, err := parse.StorageAccountManagementPolicyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.Get(ctx, rid.ResourceGroup, rid.StorageAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", rid)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", rid, err)
	}

	storageAccountID := parse.NewStorageAccountID(rid.SubscriptionId, rid.ResourceGroup, rid.StorageAccountName)
	d.Set("storage_account_id", storageAccountID.ID())

	if policy := result.Policy; policy != nil {
		policy := result.Policy
		if rules := policy.Rules; rules != nil {
			if err := d.Set("rule", flattenStorageManagementPolicyRules(rules)); err != nil {
				return fmt.Errorf("flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}

func resourceStorageManagementPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ManagementPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rid, err := parse.StorageAccountManagementPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, rid.ResourceGroup, rid.StorageAccountName); err != nil {
		return fmt.Errorf("deleting %s: %+v", rid, err)
	}
	return nil
}

// nolint unparam
func expandStorageManagementPolicyRules(d *pluginsdk.ResourceData) (*[]storage.ManagementPolicyRule, error) {
	var result []storage.ManagementPolicyRule

	rules := d.Get("rule").([]interface{})

	for k, v := range rules {
		if v != nil {
			rule, err := expandStorageManagementPolicyRule(d, k)
			if err != nil {
				return nil, fmt.Errorf("expanding the %dth rule: %+v", k, err)
			}
			_, blobIndexExist := d.GetOk(fmt.Sprintf("rule.%d.filters.0.match_blob_index_tag", k))
			_, snapshotExist := d.GetOk(fmt.Sprintf("rule.%d.actions.0.snapshot", k))
			_, versionExist := d.GetOk(fmt.Sprintf("rule.%d.actions.0.version", k))
			if blobIndexExist && (snapshotExist || versionExist) {
				return nil, fmt.Errorf("`match_blob_index_tag` is not supported as a filter for versions and snapshots")
			}
			result = append(result, *rule)
		}
	}
	return &result, nil
}

func expandStorageManagementPolicyRule(d *pluginsdk.ResourceData, ruleIndex int) (*storage.ManagementPolicyRule, error) {
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
			prefixMatchesRef := filterRef["prefix_match"].(*pluginsdk.Set)
			if prefixMatchesRef != nil {
				for _, prefixMatchRef := range prefixMatchesRef.List() {
					prefixMatches = append(prefixMatches, prefixMatchRef.(string))
				}
			}
			definition.Filters.PrefixMatch = &prefixMatches

			blobTypes := []string{}
			blobTypesRef := filterRef["blob_types"].(*pluginsdk.Set)
			if blobTypesRef != nil {
				for _, blobTypeRef := range blobTypesRef.List() {
					blobTypes = append(blobTypes, blobTypeRef.(string))
				}
			}
			definition.Filters.BlobTypes = &blobTypes

			definition.Filters.BlobIndexMatch = expandAzureRmStorageBlobIndexMatch(filterRef["match_blob_index_tag"].(*pluginsdk.Set).List())
		}
	}
	if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions", ruleIndex)); ok {
		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob", ruleIndex)); ok {
			baseBlob := &storage.ManagementPolicyBaseBlob{}
			var (
				sinceMod, sinceAccess     interface{}
				sinceModOK, sinceAccessOK bool
			)

			sinceMod, sinceModOK = d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", ruleIndex))
			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cool_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1
			if sinceModOK && sinceAccessOK {
				return nil, fmt.Errorf("can't specify `tier_to_cool_after_days_since_modification_greater_than` and `tier_to_cool_after_days_since_last_access_time_greater_than` at the same time")
			}
			if sinceModOK || sinceAccessOK {
				baseBlob.TierToCool = &storage.DateAfterModification{}
				if sinceModOK {
					baseBlob.TierToCool.DaysAfterModificationGreaterThan = utils.Float(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.TierToCool.DaysAfterLastAccessTimeGreaterThan = utils.Float(float64(sinceAccess.(int)))
				}
			}

			sinceMod, sinceModOK = d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", ruleIndex))
			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1
			if sinceModOK && sinceAccessOK {
				return nil, fmt.Errorf("can't specify `tier_to_archive_after_days_since_modification_greater_than` and `tier_to_archive_after_days_since_last_access_time_greater_than` at the same time")
			}
			if sinceModOK || sinceAccessOK {
				baseBlob.TierToArchive = &storage.DateAfterModification{}
				if sinceModOK {
					baseBlob.TierToArchive.DaysAfterModificationGreaterThan = utils.Float(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.TierToArchive.DaysAfterLastAccessTimeGreaterThan = utils.Float(float64(sinceAccess.(int)))
				}
			}

			sinceMod, sinceModOK = d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", ruleIndex))
			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.delete_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1
			if sinceModOK && sinceAccessOK {
				return nil, fmt.Errorf("can't specify `delete_after_days_since_modification_greater_than` and `delete_after_days_since_last_access_time_greater_than` at the same time")
			}
			if sinceModOK || sinceAccessOK {
				baseBlob.Delete = &storage.DateAfterModification{}
				if sinceModOK {
					baseBlob.Delete.DaysAfterModificationGreaterThan = utils.Float(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.Delete.DaysAfterLastAccessTimeGreaterThan = utils.Float(float64(sinceAccess.(int)))
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
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.change_tier_to_archive_after_days_since_creation", ruleIndex)); v != -1 {
				snapshot.TierToArchive = &storage.DateAfterCreation{
					DaysAfterCreationGreaterThan: utils.Float(float64(v.(int))),
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.change_tier_to_cool_after_days_since_creation", ruleIndex)); v != -1 {
				snapshot.TierToCool = &storage.DateAfterCreation{
					DaysAfterCreationGreaterThan: utils.Float(float64(v.(int))),
				}
			}
			definition.Actions.Snapshot = snapshot
		}

		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.version", ruleIndex)); ok {
			version := &storage.ManagementPolicyVersion{}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.delete_after_days_since_creation", ruleIndex)); v != -1 {
				version.Delete = &storage.DateAfterCreation{
					DaysAfterCreationGreaterThan: utils.Float(float64(v.(int))),
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.change_tier_to_archive_after_days_since_creation", ruleIndex)); v != -1 {
				version.TierToArchive = &storage.DateAfterCreation{
					DaysAfterCreationGreaterThan: utils.Float(float64(v.(int))),
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.change_tier_to_cool_after_days_since_creation", ruleIndex)); v != -1 {
				version.TierToCool = &storage.DateAfterCreation{
					DaysAfterCreationGreaterThan: utils.Float(float64(v.(int))),
				}
			}
			definition.Actions.Version = version
		}
	}

	return &storage.ManagementPolicyRule{
		Name:       &name,
		Enabled:    &enabled,
		Type:       &typeVal,
		Definition: &definition,
	}, nil
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

				filter["match_blob_index_tag"] = flattenAzureRmStorageBlobIndexMatch(armFilter.BlobIndexMatch)

				rule["filters"] = []interface{}{filter}
			}

			armAction := armDefinition.Actions
			if armAction != nil {
				action := make(map[string]interface{})
				armActionBaseBlob := armAction.BaseBlob
				if armActionBaseBlob != nil {
					var (
						tierToCoolSinceMod       = 0
						tierToCoolSinceAccess    = -1
						tierToArchiveSinceMod    = 0
						tierToArchiveSinceAccess = -1
						deleteSinceMod           = 0
						deleteSinceAccess        = -1
					)
					if props := armActionBaseBlob.TierToCool; props != nil {
						if props.DaysAfterModificationGreaterThan != nil {
							tierToCoolSinceMod = int(*props.DaysAfterModificationGreaterThan)
						}
						if props.DaysAfterLastAccessTimeGreaterThan != nil {
							tierToCoolSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
						}
					}
					if props := armActionBaseBlob.TierToArchive; props != nil {
						if props.DaysAfterModificationGreaterThan != nil {
							tierToArchiveSinceMod = int(*props.DaysAfterModificationGreaterThan)
						}
						if props.DaysAfterLastAccessTimeGreaterThan != nil {
							tierToArchiveSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
						}
					}
					if props := armActionBaseBlob.Delete; props != nil {
						if props.DaysAfterModificationGreaterThan != nil {
							deleteSinceMod = int(*props.DaysAfterModificationGreaterThan)
						}
						if props.DaysAfterLastAccessTimeGreaterThan != nil {
							deleteSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
						}
					}
					action["base_blob"] = []interface{}{
						map[string]interface{}{
							"tier_to_cool_after_days_since_modification_greater_than":        tierToCoolSinceMod,
							"tier_to_cool_after_days_since_last_access_time_greater_than":    tierToCoolSinceAccess,
							"tier_to_archive_after_days_since_modification_greater_than":     tierToArchiveSinceMod,
							"tier_to_archive_after_days_since_last_access_time_greater_than": tierToArchiveSinceAccess,
							"delete_after_days_since_modification_greater_than":              deleteSinceMod,
							"delete_after_days_since_last_access_time_greater_than":          deleteSinceAccess,
						},
					}
				}

				armActionSnaphost := armAction.Snapshot
				if armActionSnaphost != nil {
					deleteAfterCreation, archiveAfterCreation, coolAfterCreation := 0, -1, -1
					if armActionSnaphost.Delete != nil && armActionSnaphost.Delete.DaysAfterCreationGreaterThan != nil {
						deleteAfterCreation = int(*armActionSnaphost.Delete.DaysAfterCreationGreaterThan)
					}
					if armActionSnaphost.TierToArchive != nil && armActionSnaphost.TierToArchive.DaysAfterCreationGreaterThan != nil {
						archiveAfterCreation = int(*armActionSnaphost.TierToArchive.DaysAfterCreationGreaterThan)
					}
					if armActionSnaphost.TierToCool != nil && armActionSnaphost.TierToCool.DaysAfterCreationGreaterThan != nil {
						coolAfterCreation = int(*armActionSnaphost.TierToCool.DaysAfterCreationGreaterThan)
					}
					action["snapshot"] = []interface{}{map[string]interface{}{
						"delete_after_days_since_creation_greater_than":    deleteAfterCreation,
						"change_tier_to_archive_after_days_since_creation": archiveAfterCreation,
						"change_tier_to_cool_after_days_since_creation":    coolAfterCreation,
					}}
				}

				if armActionVersion := armAction.Version; armActionVersion != nil {
					deleteAfterCreation, archiveAfterCreation, coolAfterCreation := -1, -1, -1
					if armActionVersion.Delete != nil && armActionVersion.Delete.DaysAfterCreationGreaterThan != nil {
						deleteAfterCreation = int(*armActionVersion.Delete.DaysAfterCreationGreaterThan)
					}
					if armActionVersion.TierToArchive != nil && armActionVersion.TierToArchive.DaysAfterCreationGreaterThan != nil {
						archiveAfterCreation = int(*armActionVersion.TierToArchive.DaysAfterCreationGreaterThan)
					}
					if armActionVersion.TierToCool != nil && armActionVersion.TierToCool.DaysAfterCreationGreaterThan != nil {
						coolAfterCreation = int(*armActionVersion.TierToCool.DaysAfterCreationGreaterThan)
					}
					action["version"] = []interface{}{map[string]interface{}{
						"delete_after_days_since_creation":                 deleteAfterCreation,
						"change_tier_to_archive_after_days_since_creation": archiveAfterCreation,
						"change_tier_to_cool_after_days_since_creation":    coolAfterCreation,
					}}
				}

				rule["actions"] = []interface{}{action}
			}
		}

		rules = append(rules, rule)
	}

	return rules
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
