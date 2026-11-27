// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_management_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/managementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageManagementPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ManagementPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rid, err := parse.StorageAccountManagementPolicyID(d.Id())
	if err != nil {
		return err
	}

	accountId := commonids.NewStorageAccountID(rid.SubscriptionId, rid.ResourceGroup, rid.StorageAccountName)

	result, err := client.Get(ctx, accountId)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", rid)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", rid, err)
	}

	d.Set("storage_account_id", accountId.ID())

	if model := result.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("rule", flattenStorageManagementPolicyRules(props.Policy.Rules)); err != nil {
				return fmt.Errorf("flattening `rule`: %+v", err)
			}
		}
	}

	return nil
}

func flattenStorageManagementPolicyRules(armRules []managementpolicies.ManagementPolicyRule) []any {
	rules := make([]any, 0)
	if armRules == nil {
		return rules
	}
	for _, armRule := range armRules {
		rule := make(map[string]any)

		rule["name"] = armRule.Name
		rule["enabled"] = armRule.Enabled

		armDefinition := armRule.Definition
		armFilter := armDefinition.Filters
		if armFilter != nil {
			filter := make(map[string]any)
			if armFilter.PrefixMatch != nil {
				prefixMatches := make([]any, 0)
				for _, armPrefixMatch := range *armFilter.PrefixMatch {
					prefixMatches = append(prefixMatches, armPrefixMatch)
				}
				filter["prefix_match"] = prefixMatches
			}
			if armFilter.BlobTypes != nil {
				blobTypes := make([]any, 0)
				for _, armBlobType := range armFilter.BlobTypes {
					blobTypes = append(blobTypes, armBlobType)
				}
				filter["blob_types"] = blobTypes
			}

			filter["match_blob_index_tag"] = flattenAzureRmStorageBlobIndexMatch(armFilter.BlobIndexMatch)

			rule["filters"] = []any{filter}
		}

		armAction := armDefinition.Actions
		action := make(map[string]any)
		armActionBaseBlob := armAction.BaseBlob
		if armActionBaseBlob != nil {
			var (
				tierToCoolSinceMod               = -1
				tierToCoolSinceAccess            = -1
				tierToCoolSinceCreate            = -1
				autoTierToHotOK                  = false
				tierToArchiveSinceMod            = -1
				tierToArchiveSinceAccess         = -1
				tierToArchiveSinceCreate         = -1
				tierToArchiveSinceLastTierChange = -1
				tierToColdSinceMod               = -1
				tierToColdSinceAccess            = -1
				tierToColdSinceCreate            = -1
				deleteSinceMod                   = -1
				deleteSinceAccess                = -1
				deleteSinceCreate                = -1
			)

			if v := armActionBaseBlob.EnableAutoTierToHotFromCool; v != nil {
				autoTierToHotOK = *v
			}
			if props := armActionBaseBlob.TierToCool; props != nil {
				if props.DaysAfterModificationGreaterThan != nil {
					tierToCoolSinceMod = int(*props.DaysAfterModificationGreaterThan)
				}
				if props.DaysAfterLastAccessTimeGreaterThan != nil {
					tierToCoolSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
				}
				if props.DaysAfterCreationGreaterThan != nil {
					tierToCoolSinceCreate = int(*props.DaysAfterCreationGreaterThan)
				}
			}
			if props := armActionBaseBlob.TierToArchive; props != nil {
				if props.DaysAfterModificationGreaterThan != nil {
					tierToArchiveSinceMod = int(*props.DaysAfterModificationGreaterThan)
				}
				if props.DaysAfterLastAccessTimeGreaterThan != nil {
					tierToArchiveSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
				}
				if props.DaysAfterLastTierChangeGreaterThan != nil {
					tierToArchiveSinceLastTierChange = int(*props.DaysAfterLastTierChangeGreaterThan)
				}
				if props.DaysAfterCreationGreaterThan != nil {
					tierToArchiveSinceCreate = int(*props.DaysAfterCreationGreaterThan)
				}
			}
			if props := armActionBaseBlob.TierToCold; props != nil {
				if props.DaysAfterModificationGreaterThan != nil {
					tierToColdSinceMod = int(*props.DaysAfterModificationGreaterThan)
				}
				if props.DaysAfterLastAccessTimeGreaterThan != nil {
					tierToColdSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
				}
				if props.DaysAfterCreationGreaterThan != nil {
					tierToColdSinceCreate = int(*props.DaysAfterCreationGreaterThan)
				}
			}
			if props := armActionBaseBlob.Delete; props != nil {
				if props.DaysAfterModificationGreaterThan != nil {
					deleteSinceMod = int(*props.DaysAfterModificationGreaterThan)
				}
				if props.DaysAfterLastAccessTimeGreaterThan != nil {
					deleteSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
				}
				if props.DaysAfterCreationGreaterThan != nil {
					deleteSinceCreate = int(*props.DaysAfterCreationGreaterThan)
				}
			}
			action["base_blob"] = []any{
				map[string]any{
					"auto_tier_to_hot_from_cool_enabled":                             autoTierToHotOK,
					"tier_to_cool_after_days_since_modification_greater_than":        tierToCoolSinceMod,
					"tier_to_cool_after_days_since_last_access_time_greater_than":    tierToCoolSinceAccess,
					"tier_to_cool_after_days_since_creation_greater_than":            tierToCoolSinceCreate,
					"tier_to_archive_after_days_since_modification_greater_than":     tierToArchiveSinceMod,
					"tier_to_archive_after_days_since_last_access_time_greater_than": tierToArchiveSinceAccess,
					"tier_to_archive_after_days_since_last_tier_change_greater_than": tierToArchiveSinceLastTierChange,
					"tier_to_archive_after_days_since_creation_greater_than":         tierToArchiveSinceCreate,
					"tier_to_cold_after_days_since_modification_greater_than":        tierToColdSinceMod,
					"tier_to_cold_after_days_since_last_access_time_greater_than":    tierToColdSinceAccess,
					"tier_to_cold_after_days_since_creation_greater_than":            tierToColdSinceCreate,
					"delete_after_days_since_modification_greater_than":              deleteSinceMod,
					"delete_after_days_since_last_access_time_greater_than":          deleteSinceAccess,
					"delete_after_days_since_creation_greater_than":                  deleteSinceCreate,
				},
			}
		}

		armActionSnapshot := armAction.Snapshot
		if armActionSnapshot != nil {
			var (
				deleteAfterCreation        = -1
				archiveAfterCreation       = -1
				archiveAfterLastTierChange = -1
				coolAfterCreation          = -1
				tierToColdSinceCreate      = -1
			)
			if armActionSnapshot.Delete != nil {
				deleteAfterCreation = int(armActionSnapshot.Delete.DaysAfterCreationGreaterThan)
			}
			if armActionSnapshot.TierToArchive != nil {
				archiveAfterCreation = int(armActionSnapshot.TierToArchive.DaysAfterCreationGreaterThan)

				if v := armActionSnapshot.TierToArchive.DaysAfterLastTierChangeGreaterThan; v != nil {
					archiveAfterLastTierChange = int(*v)
				}
			}
			if armActionSnapshot.TierToCold != nil {
				tierToColdSinceCreate = int(armActionSnapshot.TierToCold.DaysAfterCreationGreaterThan)
			}
			if armActionSnapshot.TierToCool != nil {
				coolAfterCreation = int(armActionSnapshot.TierToCool.DaysAfterCreationGreaterThan)
			}
			action["snapshot"] = []any{map[string]any{
				"delete_after_days_since_creation_greater_than":                  deleteAfterCreation,
				"change_tier_to_archive_after_days_since_creation":               archiveAfterCreation,
				"tier_to_archive_after_days_since_last_tier_change_greater_than": archiveAfterLastTierChange,
				"tier_to_cold_after_days_since_creation_greater_than":            tierToColdSinceCreate,
				"change_tier_to_cool_after_days_since_creation":                  coolAfterCreation,
			}}
		}

		if armActionVersion := armAction.Version; armActionVersion != nil {
			var (
				deleteAfterCreation        = -1
				archiveAfterCreation       = -1
				archiveAfterLastTierChange = -1
				coolAfterCreation          = -1
				tierToColdSinceCreate      = -1
			)
			if armActionVersion.Delete != nil {
				deleteAfterCreation = int(armActionVersion.Delete.DaysAfterCreationGreaterThan)
			}
			if armActionVersion.TierToArchive != nil {
				archiveAfterCreation = int(armActionVersion.TierToArchive.DaysAfterCreationGreaterThan)

				if v := armActionVersion.TierToArchive.DaysAfterLastTierChangeGreaterThan; v != nil {
					archiveAfterLastTierChange = int(*v)
				}
			}
			if armActionVersion.TierToCold != nil {
				tierToColdSinceCreate = int(armActionVersion.TierToCold.DaysAfterCreationGreaterThan)
			}
			if armActionVersion.TierToCool != nil {
				coolAfterCreation = int(armActionVersion.TierToCool.DaysAfterCreationGreaterThan)
			}
			action["version"] = []any{map[string]any{
				"delete_after_days_since_creation":                               deleteAfterCreation,
				"change_tier_to_archive_after_days_since_creation":               archiveAfterCreation,
				"tier_to_archive_after_days_since_last_tier_change_greater_than": archiveAfterLastTierChange,
				"tier_to_cold_after_days_since_creation_greater_than":            tierToColdSinceCreate,
				"change_tier_to_cool_after_days_since_creation":                  coolAfterCreation,
			}}
		}

		rule["actions"] = []any{action}

		rules = append(rules, rule)
	}

	return rules
}

func flattenAzureRmStorageBlobIndexMatch(blobIndexMatches *[]managementpolicies.TagFilter) []map[string]any {
	result := make([]map[string]any, 0)

	if blobIndexMatches == nil || len(*blobIndexMatches) == 0 {
		return result
	}

	for _, blobIndexMatch := range *blobIndexMatches {
		result = append(result, map[string]any{
			"name":      blobIndexMatch.Name,
			"operation": blobIndexMatch.Op,
			"value":     blobIndexMatch.Value,
		})
	}
	return result
}
