// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_management_policy

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/managementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandStorageManagementPolicyRules(d *pluginsdk.ResourceData) ([]managementpolicies.ManagementPolicyRule, error) {
	var result []managementpolicies.ManagementPolicyRule

	rules := d.Get("rule").([]any)

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
				return nil, errors.New("`match_blob_index_tag` is not supported as a filter for versions and snapshots")
			}
			result = append(result, *rule)
		}
	}
	return result, nil
}

func expandStorageManagementPolicyRule(d *pluginsdk.ResourceData, ruleIndex int) (*managementpolicies.ManagementPolicyRule, error) {
	name := d.Get(fmt.Sprintf("rule.%d.name", ruleIndex)).(string)
	enabled := d.Get(fmt.Sprintf("rule.%d.enabled", ruleIndex)).(bool)
	typeVal := "Lifecycle"

	definition := managementpolicies.ManagementPolicyDefinition{
		Filters: &managementpolicies.ManagementPolicyFilter{},
		Actions: managementpolicies.ManagementPolicyAction{},
	}
	filtersRef := d.Get(fmt.Sprintf("rule.%d.filters", ruleIndex)).([]any)
	if len(filtersRef) == 1 {
		if filtersRef[0] != nil {
			filterRef := filtersRef[0].(map[string]any)

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
			definition.Filters.BlobTypes = blobTypes

			definition.Filters.BlobIndexMatch = expandAzureRmStorageBlobIndexMatch(filterRef["match_blob_index_tag"].(*pluginsdk.Set).List())
		}
	}
	if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions", ruleIndex)); ok {
		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.base_blob", ruleIndex)); ok {
			baseBlob := &managementpolicies.ManagementPolicyBaseBlob{}
			var (
				sinceMod, sinceAccess, sinceCreate       any
				sinceModOK, sinceAccessOK, sinceCreateOK bool
			)

			sinceMod = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", ruleIndex))
			sinceModOK = sinceMod != -1

			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cool_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1

			sinceCreate = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cool_after_days_since_creation_greater_than", ruleIndex))
			sinceCreateOK = sinceCreate != -1

			autoTierToHotOK := d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.auto_tier_to_hot_from_cool_enabled", ruleIndex)).(bool)
			if autoTierToHotOK && !sinceAccessOK {
				return nil, errors.New("`auto_tier_to_hot_from_cool_enabled` must be used together with `tier_to_cool_after_days_since_last_access_time_greater_than`")
			}

			var cnt int
			if sinceModOK {
				cnt++
			}
			if sinceAccessOK {
				cnt++
			}
			if sinceCreateOK {
				cnt++
			}
			if cnt > 1 {
				return nil, errors.New("only one of `tier_to_cool_after_days_since_modification_greater_than`, `tier_to_cool_after_days_since_last_access_time_greater_than`, `tier_to_cool_after_days_since_creation_greater_than` can be specified at the same time")
			}

			if sinceModOK || sinceAccessOK || sinceCreateOK {
				baseBlob.TierToCool = &managementpolicies.DateAfterModification{}
				if sinceModOK {
					baseBlob.TierToCool.DaysAfterModificationGreaterThan = pointer.To(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.TierToCool.DaysAfterLastAccessTimeGreaterThan = pointer.To(float64(sinceAccess.(int)))
				}
				if sinceCreateOK {
					baseBlob.TierToCool.DaysAfterCreationGreaterThan = pointer.To(float64(sinceCreate.(int)))
				}
				if autoTierToHotOK {
					baseBlob.EnableAutoTierToHotFromCool = pointer.To(autoTierToHotOK)
				}
			}

			sinceMod = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", ruleIndex))
			sinceModOK = sinceMod != -1
			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1
			sinceCreate = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_creation_greater_than", ruleIndex))
			sinceCreateOK = sinceCreate != -1

			cnt = 0
			if sinceModOK {
				cnt++
			}
			if sinceAccessOK {
				cnt++
			}
			if sinceCreateOK {
				cnt++
			}
			if cnt > 1 {
				return nil, errors.New("only one of `tier_to_archive_after_days_since_modification_greater_than`, `tier_to_archive_after_days_since_last_access_time_greater_than` and `tier_to_archive_after_days_since_creation_greater_than` can be specified at the same time")
			}

			if sinceModOK || sinceAccessOK || sinceCreateOK {
				baseBlob.TierToArchive = &managementpolicies.DateAfterModification{}
				if sinceModOK {
					baseBlob.TierToArchive.DaysAfterModificationGreaterThan = pointer.To(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.TierToArchive.DaysAfterLastAccessTimeGreaterThan = pointer.To(float64(sinceAccess.(int)))
				}
				if sinceCreateOK {
					baseBlob.TierToArchive.DaysAfterCreationGreaterThan = pointer.To(float64(sinceCreate.(int)))
				}
				if v := d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_archive_after_days_since_last_tier_change_greater_than", ruleIndex)); v != -1 {
					baseBlob.TierToArchive.DaysAfterLastTierChangeGreaterThan = pointer.To(float64(v.(int)))
				}
			}

			sinceMod = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", ruleIndex))
			sinceModOK = sinceMod != -1
			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.delete_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1
			sinceCreate = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.delete_after_days_since_creation_greater_than", ruleIndex))
			sinceCreateOK = sinceCreate != -1

			cnt = 0
			if sinceModOK {
				cnt++
			}
			if sinceAccessOK {
				cnt++
			}
			if sinceCreateOK {
				cnt++
			}
			if cnt > 1 {
				return nil, errors.New("only one of `delete_after_days_since_modification_greater_than`, `delete_after_days_since_last_access_time_greater_than` and `delete_after_days_since_creation_greater_than` can be specified at the same time")
			}
			if sinceModOK || sinceAccessOK || sinceCreateOK {
				baseBlob.Delete = &managementpolicies.DateAfterModification{}
				if sinceModOK {
					baseBlob.Delete.DaysAfterModificationGreaterThan = pointer.To(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.Delete.DaysAfterLastAccessTimeGreaterThan = pointer.To(float64(sinceAccess.(int)))
				}
				if sinceCreateOK {
					baseBlob.Delete.DaysAfterCreationGreaterThan = pointer.To(float64(sinceCreate.(int)))
				}
			}

			sinceMod = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cold_after_days_since_modification_greater_than", ruleIndex))
			sinceModOK = sinceMod != -1
			sinceAccess = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cold_after_days_since_last_access_time_greater_than", ruleIndex))
			sinceAccessOK = sinceAccess != -1
			sinceCreate = d.Get(fmt.Sprintf("rule.%d.actions.0.base_blob.0.tier_to_cold_after_days_since_creation_greater_than", ruleIndex))
			sinceCreateOK = sinceCreate != -1

			cnt = 0
			if sinceModOK {
				cnt++
			}
			if sinceAccessOK {
				cnt++
			}
			if sinceCreateOK {
				cnt++
			}
			if cnt > 1 {
				return nil, errors.New("only one of `tier_to_cold_after_days_since_modification_greater_than`, `tier_to_cold_after_days_since_last_access_time_greater_than` and `tier_to_cold_after_days_since_creation_greater_than` can be specified at the same time")
			}

			if sinceModOK || sinceAccessOK || sinceCreateOK {
				baseBlob.TierToCold = &managementpolicies.DateAfterModification{}
				if sinceModOK {
					baseBlob.TierToCold.DaysAfterModificationGreaterThan = pointer.To(float64(sinceMod.(int)))
				}
				if sinceAccessOK {
					baseBlob.TierToCold.DaysAfterLastAccessTimeGreaterThan = pointer.To(float64(sinceAccess.(int)))
				}
				if sinceCreateOK {
					baseBlob.TierToCold.DaysAfterCreationGreaterThan = pointer.To(float64(sinceCreate.(int)))
				}
			}

			definition.Actions.BaseBlob = baseBlob
		}

		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.snapshot", ruleIndex)); ok {
			snapshot := &managementpolicies.ManagementPolicySnapShot{}

			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", ruleIndex)); v != -1 {
				snapshot.Delete = &managementpolicies.DateAfterCreation{DaysAfterCreationGreaterThan: float64(v.(int))}
			}

			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.change_tier_to_archive_after_days_since_creation", ruleIndex)); v != -1 {
				snapshot.TierToArchive = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
				if vv := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.tier_to_archive_after_days_since_last_tier_change_greater_than", ruleIndex)); vv != -1 {
					snapshot.TierToArchive.DaysAfterLastTierChangeGreaterThan = pointer.To(float64(vv.(int)))
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.change_tier_to_cool_after_days_since_creation", ruleIndex)); v != -1 {
				snapshot.TierToCool = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.snapshot.0.tier_to_cold_after_days_since_creation_greater_than", ruleIndex)); v != -1 {
				snapshot.TierToCold = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
			}
			definition.Actions.Snapshot = snapshot
		}

		if _, ok := d.GetOk(fmt.Sprintf("rule.%d.actions.0.version", ruleIndex)); ok {
			version := &managementpolicies.ManagementPolicyVersion{}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.delete_after_days_since_creation", ruleIndex)); v != -1 {
				version.Delete = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.change_tier_to_archive_after_days_since_creation", ruleIndex)); v != -1 {
				version.TierToArchive = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
				if vv := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.tier_to_archive_after_days_since_last_tier_change_greater_than", ruleIndex)); vv != -1 {
					version.TierToArchive.DaysAfterLastTierChangeGreaterThan = pointer.To(float64(vv.(int)))
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.change_tier_to_cool_after_days_since_creation", ruleIndex)); v != -1 {
				version.TierToCool = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
			}
			if v := d.Get(fmt.Sprintf("rule.%d.actions.0.version.0.tier_to_cold_after_days_since_creation_greater_than", ruleIndex)); v != -1 {
				version.TierToCold = &managementpolicies.DateAfterCreation{
					DaysAfterCreationGreaterThan: float64(v.(int)),
				}
			}
			definition.Actions.Version = version
		}
	}

	return &managementpolicies.ManagementPolicyRule{
		Name:       name,
		Enabled:    &enabled,
		Type:       managementpolicies.RuleType(typeVal),
		Definition: definition,
	}, nil
}

func expandAzureRmStorageBlobIndexMatch(blobIndexMatches []any) *[]managementpolicies.TagFilter {
	if len(blobIndexMatches) == 0 {
		return nil
	}

	results := make([]managementpolicies.TagFilter, 0)
	for _, v := range blobIndexMatches {
		blobIndexMatch := v.(map[string]any)

		filter := managementpolicies.TagFilter{
			Name:  blobIndexMatch["name"].(string),
			Op:    blobIndexMatch["operation"].(string),
			Value: blobIndexMatch["value"].(string),
		}

		results = append(results, filter)
	}

	return &results
}
