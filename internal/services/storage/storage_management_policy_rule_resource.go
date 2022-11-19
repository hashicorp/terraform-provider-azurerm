package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageAccountManagementPolicyRuleResource struct{}

var _ sdk.ResourceWithUpdate = StorageAccountManagementPolicyRuleResource{}

type StorageAccountManagementPolicyRuleModel struct {
	Name               string                                          `tfschema:"name"`
	ManagementPolicyId string                                          `tfschema:"management_policy_id"`
	Enabled            bool                                            `tfschema:"enabled"`
	Actions            []StorageAccountManagementPolicyRuleActionModel `tfschema:"action"`
	Filters            []StorageAccountManagementPolicyRuleFilterModel `tfschema:"filter"`
}

type StorageAccountManagementPolicyRuleActionBaseBlobModel struct {
	DeleteAfterDaysSinceLastAccessTimeGreaterThan        int `tfschema:"delete_after_days_since_last_access_time_greater_than"`
	DeleteAfterDaysSinceModificationGreaterThan          int `tfschema:"delete_after_days_since_modification_greater_than"`
	TierToArchiveAfterDaysSinceLastAccessTimeGreaterThan int `tfschema:"tier_to_archive_after_days_since_last_access_time_greater_than"`
	TierToArchiveAfterDaysSinceLastTierChangeGreaterThan int `tfschema:"tier_to_archive_after_days_since_last_tier_change_greater_than"`
	TierToArchiveAfterDaysSinceModificationGreaterThan   int `tfschema:"tier_to_archive_after_days_since_modification_greater_than"`
	TierToCoolAfterDaysSinceLastAccessTimeGreaterThan    int `tfschema:"tier_to_cool_after_days_since_last_access_time_greater_than"`
	TierToCoolAfterDaysSinceModificationGreaterThan      int `tfschema:"tier_to_cool_after_days_since_modification_greater_than"`
}
type StorageAccountManagementPolicyRuleActionSnapshotModel struct {
	ChangeTierToArchiveAfterDaysSinceCreation            int `tfschema:"change_tier_to_archive_after_days_since_creation"`
	ChangeTierToCoolAfterDaysSinceCreation               int `tfschema:"change_tier_to_cool_after_days_since_creation"`
	DeleteAfterDaysSinceCreationGreaterThan              int `tfschema:"delete_after_days_since_creation_greater_than"`
	TierToArchiveAfterDaysSinceLastTierChangeGreaterThan int `tfschema:"tier_to_archive_after_days_since_last_tier_change_greater_than"`
}
type StorageAccountManagementPolicyRuleActionVersionModel struct {
	ChangeTierToArchiveAfterDaysSinceCreation            int `tfschema:"change_tier_to_archive_after_days_since_creation"`
	ChangeTierToCoolAfterDaysSinceCreation               int `tfschema:"change_tier_to_cool_after_days_since_creation"`
	DeleteAfterDaysSinceCreation                         int `tfschema:"delete_after_days_since_creation"`
	TierToArchiveAfterDaysSinceLastTierChangeGreaterThan int `tfschema:"tier_to_archive_after_days_since_last_tier_change_greater_than"`
}
type StorageAccountManagementPolicyRuleActionModel struct {
	BaseBlob []StorageAccountManagementPolicyRuleActionBaseBlobModel `tfschema:"base_blob"`
	Snapshot []StorageAccountManagementPolicyRuleActionSnapshotModel `tfschema:"snapshot"`
	Version  []StorageAccountManagementPolicyRuleActionVersionModel  `tfschema:"version"`
}
type StorageAccountManagementPolicyRuleFilterMatchBlobIndexTagModel struct {
	Name      string `tfschema:"name"`
	Operation string `tfschema:"operation"`
	Value     string `tfschema:"value"`
}
type StorageAccountManagementPolicyRuleFilterModel struct {
	BlobTypes         []string                                                         `tfschema:"blob_types"`
	MatchBlobIndexTag []StorageAccountManagementPolicyRuleFilterMatchBlobIndexTagModel `tfschema:"match_blob_index_tag"`
	PrefixMatch       []string                                                         `tfschema:"prefix_match"`
}

func (r StorageAccountManagementPolicyRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"management_policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageAccountManagementPolicyID,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"action": {
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
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
									ValidateFunc: validation.IntBetween(0, 99999),
								},
								"tier_to_cool_after_days_since_last_access_time_greater_than": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
									ValidateFunc: validation.IntBetween(0, 99999),
								},
								"tier_to_archive_after_days_since_modification_greater_than": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
									ValidateFunc: validation.IntBetween(0, 99999),
								},
								"tier_to_archive_after_days_since_last_access_time_greater_than": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
									ValidateFunc: validation.IntBetween(0, 99999),
								},
								"tier_to_archive_after_days_since_last_tier_change_greater_than": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
									ValidateFunc: validation.IntBetween(0, 99999),
								},
								"delete_after_days_since_modification_greater_than": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
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
								"tier_to_archive_after_days_since_last_tier_change_greater_than": {
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
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      -1,
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
								"tier_to_archive_after_days_since_last_tier_change_greater_than": {
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
		"filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_types": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"blockBlob",
								"appendBlob",
							}, false),
						},
					},
					"prefix_match": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"match_blob_index_tag": {
						Type:     pluginsdk.TypeList,
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
	}
}

func (r StorageAccountManagementPolicyRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageAccountManagementPolicyRuleResource) ResourceType() string {
	return "azurerm_storage_management_policy_rule"
}

func (r StorageAccountManagementPolicyRuleResource) ModelObject() interface{} {
	return &StorageAccountManagementPolicyRuleModel{}
}

func (r StorageAccountManagementPolicyRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StorageAccountManagementPolicyRuleID
}

func (r StorageAccountManagementPolicyRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ManagementPoliciesClient

			var plan StorageAccountManagementPolicyRuleModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			policyId, err := parse.StorageAccountManagementPolicyID(plan.ManagementPolicyId)
			if err != nil {
				return err
			}

			id := parse.NewStorageAccountManagementPolicyRuleID(policyId.SubscriptionId, policyId.ResourceGroup, policyId.StorageAccountName, policyId.ManagementPolicyName, plan.Name)

			locks.ByName(id.StorageAccountName, StorageAccountManagementPolicyResourceName)
			defer locks.UnlockByName(id.StorageAccountName, StorageAccountManagementPolicyResourceName)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", policyId, err)
				}
			}
			if utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if r.GetRuleInPolicy(existing, plan.Name) != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			rule := &storage.ManagementPolicyRule{
				Name:    &id.RuleName,
				Enabled: &plan.Enabled,
				Type:    utils.String("Lifecycle"),
				Definition: &storage.ManagementPolicyDefinition{
					Filters: r.expandFilter(plan.Filters),
					Actions: r.expandAction(plan.Actions),
				},
			}

			r.AddRuleInPolicy(&existing, *rule)

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, existing); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r StorageAccountManagementPolicyRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ManagementPoliciesClient
			id, err := parse.StorageAccountManagementPolicyRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			policyId := parse.NewStorageAccountManagementPolicyID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ManagementPolicyName)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", policyId, err)
			}

			rule := r.GetRuleInPolicy(existing, id.RuleName)
			if rule == nil {
				return metadata.MarkAsGone(id)
			}

			enabled := false
			if rule.Enabled != nil {
				enabled = *rule.Enabled
			}

			model := StorageAccountManagementPolicyRuleModel{
				Name:               id.RuleName,
				ManagementPolicyId: policyId.ID(),
				Enabled:            enabled,
			}

			if rule.Definition != nil {
				model.Actions = r.flattenAction(rule.Definition.Actions)
				model.Filters = r.flattenFilter(rule.Definition.Filters)
			}

			return metadata.Encode(&model)
		},
	}
}

func (r StorageAccountManagementPolicyRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.StorageAccountManagementPolicyRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.StorageAccountName, StorageAccountManagementPolicyResourceName)
			defer locks.UnlockByName(id.StorageAccountName, StorageAccountManagementPolicyResourceName)

			var plan StorageAccountManagementPolicyRuleModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := metadata.Client.Storage.ManagementPoliciesClient

			policyId := parse.NewStorageAccountManagementPolicyID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ManagementPolicyName)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", policyId, err)
			}

			rule := r.GetRuleInPolicy(existing, id.RuleName)
			if rule == nil {
				return fmt.Errorf("retrieving %s: no such rule found", id)
			}

			if metadata.ResourceData.HasChange("enabled") {
				rule.Enabled = utils.Bool(plan.Enabled)
			}
			if metadata.ResourceData.HasChange("action") || metadata.ResourceData.HasChange("filter") {
				if rule.Definition == nil {
					rule.Definition = &storage.ManagementPolicyDefinition{}
				}
				if metadata.ResourceData.HasChange("action") {
					rule.Definition.Actions = r.expandAction(plan.Actions)
				}
				if metadata.ResourceData.HasChange("filter") {
					rule.Definition.Filters = r.expandFilter(plan.Filters)
				}
			}

			if err := r.UpdateRuleInPolicy(&existing, id.RuleName, *rule); err != nil {
				return fmt.Errorf("updating %s within the management policy: %v", id, err)
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r StorageAccountManagementPolicyRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ManagementPoliciesClient

			id, err := parse.StorageAccountManagementPolicyRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.StorageAccountName, StorageAccountManagementPolicyResourceName)
			defer locks.UnlockByName(id.StorageAccountName, StorageAccountManagementPolicyResourceName)

			policyId := parse.NewStorageAccountManagementPolicyID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.ManagementPolicyName)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", policyId, err)
			}

			r.DeleteRuleInPolicy(existing, id.RuleName)

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r StorageAccountManagementPolicyRuleResource) AddRuleInPolicy(policy *storage.ManagementPolicy, rule storage.ManagementPolicyRule) {
	if policy == nil {
		return
	}
	if policy.ManagementPolicyProperties == nil {
		policy.ManagementPolicyProperties = &storage.ManagementPolicyProperties{}
	}
	if policy.Policy == nil {
		policy.Policy = &storage.ManagementPolicySchema{}
	}
	if policy.Policy.Rules == nil {
		policy.Policy.Rules = &[]storage.ManagementPolicyRule{}
	}
	*policy.Policy.Rules = append(*policy.Policy.Rules, rule)
	return
}

func (r StorageAccountManagementPolicyRuleResource) GetRuleInPolicy(policy storage.ManagementPolicy, name string) *storage.ManagementPolicyRule {
	if policy.ManagementPolicyProperties == nil {
		return nil
	}
	if policy.ManagementPolicyProperties.Policy == nil {
		return nil
	}
	if policy.ManagementPolicyProperties.Policy.Rules == nil {
		return nil
	}
	for _, rule := range *policy.ManagementPolicyProperties.Policy.Rules {
		if rule.Name == nil {
			continue
		}
		if name == *rule.Name {
			return &rule
		}
	}
	return nil
}

func (r StorageAccountManagementPolicyRuleResource) UpdateRuleInPolicy(policy *storage.ManagementPolicy, name string, rule storage.ManagementPolicyRule) error {
	if policy.ManagementPolicyProperties == nil {
		return fmt.Errorf("nil .properties of the policy")
	}
	if policy.ManagementPolicyProperties.Policy == nil {
		return fmt.Errorf("nil properties.policy of the policy")
	}
	if policy.ManagementPolicyProperties.Policy.Rules == nil {
		return fmt.Errorf("nil properties.policy.rules of the policy")
	}

	var (
		updated  bool
		newRules []storage.ManagementPolicyRule
	)

	for _, r := range *policy.ManagementPolicyProperties.Policy.Rules {
		if r.Name != nil && name == *r.Name {
			newRules = append(newRules, rule)
			updated = true
			continue
		}
		newRules = append(newRules, r)
	}

	*policy.ManagementPolicyProperties.Policy.Rules = newRules

	if !updated {
		return fmt.Errorf("no rule named %s", name)
	}

	return nil
}

func (r StorageAccountManagementPolicyRuleResource) DeleteRuleInPolicy(policy storage.ManagementPolicy, name string) {
	if policy.ManagementPolicyProperties == nil {
		return
	}
	if policy.ManagementPolicyProperties.Policy == nil {
		return
	}
	if policy.ManagementPolicyProperties.Policy.Rules == nil {
		return
	}

	var newRules []storage.ManagementPolicyRule
	for _, r := range *policy.ManagementPolicyProperties.Policy.Rules {
		if r.Name != nil && name == *r.Name {
			continue
		}
		newRules = append(newRules, r)
	}

	*policy.ManagementPolicyProperties.Policy.Rules = newRules
}

func (r StorageAccountManagementPolicyRuleResource) expandFilter(input []StorageAccountManagementPolicyRuleFilterModel) *storage.ManagementPolicyFilter {
	if len(input) == 0 {
		return nil
	}

	filter := input[0]

	out := &storage.ManagementPolicyFilter{
		BlobTypes: &filter.BlobTypes,
	}

	if len(filter.PrefixMatch) != 0 {
		var matches []string
		for _, match := range filter.PrefixMatch {
			matches = append(matches, match)
		}
		out.PrefixMatch = &matches
	}

	if len(filter.MatchBlobIndexTag) != 0 {
		var tags []storage.TagFilter
		for _, tag := range filter.MatchBlobIndexTag {
			tags = append(tags, storage.TagFilter{
				Name:  &tag.Name,
				Op:    &tag.Operation,
				Value: &tag.Value,
			})
		}
		out.BlobIndexMatch = &tags
	}

	return out
}

func (r StorageAccountManagementPolicyRuleResource) flattenFilter(input *storage.ManagementPolicyFilter) []StorageAccountManagementPolicyRuleFilterModel {
	if input == nil {
		return nil
	}

	filterModel := StorageAccountManagementPolicyRuleFilterModel{}

	if input.BlobTypes != nil {
		filterModel.BlobTypes = *input.BlobTypes
	}
	if input.BlobIndexMatch != nil {
		tags := []StorageAccountManagementPolicyRuleFilterMatchBlobIndexTagModel{}
		for _, match := range *input.BlobIndexMatch {
			tag := StorageAccountManagementPolicyRuleFilterMatchBlobIndexTagModel{}
			if match.Name != nil {
				tag.Name = *match.Name
			}
			if match.Op != nil {
				tag.Operation = *match.Op
			}
			if match.Value != nil {
				tag.Value = *match.Value
			}
			tags = append(tags, tag)
		}
		filterModel.MatchBlobIndexTag = tags
	}
	if input.PrefixMatch != nil {
		filterModel.PrefixMatch = *input.PrefixMatch
	}

	return []StorageAccountManagementPolicyRuleFilterModel{filterModel}
}

func (r StorageAccountManagementPolicyRuleResource) expandAction(input []StorageAccountManagementPolicyRuleActionModel) *storage.ManagementPolicyAction {
	if len(input) == 0 {
		return nil
	}
	actionModel := input[0]

	action := &storage.ManagementPolicyAction{}

	if len(actionModel.BaseBlob) != 0 {
		baseBlobModel := actionModel.BaseBlob[0]

		baseBlob := &storage.ManagementPolicyBaseBlob{}

		// Tier to cool
		sinceMod := baseBlobModel.TierToCoolAfterDaysSinceModificationGreaterThan
		sinceModOK := sinceMod != -1
		sinceAccess := baseBlobModel.TierToCoolAfterDaysSinceLastAccessTimeGreaterThan
		sinceAccessOK := sinceAccess != -1
		if sinceModOK || sinceAccessOK {
			baseBlob.TierToCool = &storage.DateAfterModification{}
			if sinceModOK {
				baseBlob.TierToCool.DaysAfterModificationGreaterThan = utils.Float(float64(sinceMod))
			}
			if sinceAccessOK {
				baseBlob.TierToCool.DaysAfterLastAccessTimeGreaterThan = utils.Float(float64(sinceAccess))
			}
		}

		// Tier to archive
		sinceMod = baseBlobModel.TierToArchiveAfterDaysSinceModificationGreaterThan
		sinceModOK = sinceMod != -1
		sinceAccess = baseBlobModel.TierToArchiveAfterDaysSinceLastAccessTimeGreaterThan
		sinceAccessOK = sinceAccess != -1
		if sinceModOK || sinceAccessOK {
			baseBlob.TierToArchive = &storage.DateAfterModification{}
			if sinceModOK {
				baseBlob.TierToArchive.DaysAfterModificationGreaterThan = utils.Float(float64(sinceMod))
			}
			if sinceAccessOK {
				baseBlob.TierToArchive.DaysAfterLastAccessTimeGreaterThan = utils.Float(float64(sinceAccess))
			}
			if v := baseBlobModel.TierToArchiveAfterDaysSinceLastTierChangeGreaterThan; v != -1 {
				baseBlob.TierToArchive.DaysAfterLastTierChangeGreaterThan = utils.Float(float64(v))
			}
		}

		// Delete
		sinceMod = baseBlobModel.DeleteAfterDaysSinceModificationGreaterThan
		sinceModOK = sinceMod != -1
		sinceAccess = baseBlobModel.DeleteAfterDaysSinceLastAccessTimeGreaterThan
		sinceAccessOK = sinceAccess != -1
		if sinceModOK || sinceAccessOK {
			baseBlob.Delete = &storage.DateAfterModification{}
			if sinceModOK {
				baseBlob.Delete.DaysAfterModificationGreaterThan = utils.Float(float64(sinceMod))
			}
			if sinceAccessOK {
				baseBlob.Delete.DaysAfterLastAccessTimeGreaterThan = utils.Float(float64(sinceAccess))
			}
		}

		action.BaseBlob = baseBlob
	}

	if len(actionModel.Snapshot) != 0 {
		snapshotModel := actionModel.Snapshot[0]

		snapshot := &storage.ManagementPolicySnapShot{}

		if v := snapshotModel.DeleteAfterDaysSinceCreationGreaterThan; v != -1 {
			snapshot.Delete = &storage.DateAfterCreation{DaysAfterCreationGreaterThan: utils.Float(float64(v))}
		}
		if v := snapshotModel.ChangeTierToArchiveAfterDaysSinceCreation; v != -1 {
			snapshot.TierToArchive = &storage.DateAfterCreation{
				DaysAfterCreationGreaterThan: utils.Float(float64(v)),
			}
			if v := snapshotModel.TierToArchiveAfterDaysSinceLastTierChangeGreaterThan; v != -1 {
				snapshot.TierToArchive.DaysAfterLastTierChangeGreaterThan = utils.Float(float64(v))
			}
		}
		if v := snapshotModel.ChangeTierToCoolAfterDaysSinceCreation; v != -1 {
			snapshot.TierToCool = &storage.DateAfterCreation{
				DaysAfterCreationGreaterThan: utils.Float(float64(v)),
			}
		}

		action.Snapshot = snapshot
	}

	if len(actionModel.Version) != 0 {
		versionModel := actionModel.Version[0]

		version := &storage.ManagementPolicyVersion{}

		if v := versionModel.DeleteAfterDaysSinceCreation; v != -1 {
			version.Delete = &storage.DateAfterCreation{
				DaysAfterCreationGreaterThan: utils.Float(float64(v)),
			}
		}
		if v := versionModel.ChangeTierToArchiveAfterDaysSinceCreation; v != -1 {
			version.TierToArchive = &storage.DateAfterCreation{
				DaysAfterCreationGreaterThan: utils.Float(float64(v)),
			}
			if v := versionModel.TierToArchiveAfterDaysSinceLastTierChangeGreaterThan; v != -1 {
				version.TierToArchive.DaysAfterLastTierChangeGreaterThan = utils.Float(float64(v))
			}
		}
		if v := versionModel.ChangeTierToCoolAfterDaysSinceCreation; v != -1 {
			version.TierToCool = &storage.DateAfterCreation{
				DaysAfterCreationGreaterThan: utils.Float(float64(v)),
			}
		}

		action.Version = version
	}

	return action
}

func (r StorageAccountManagementPolicyRuleResource) flattenAction(input *storage.ManagementPolicyAction) []StorageAccountManagementPolicyRuleActionModel {
	if input == nil {
		return nil
	}

	actionModel := StorageAccountManagementPolicyRuleActionModel{}

	if input.BaseBlob != nil {
		var (
			tierToCoolSinceMod               = -1
			tierToCoolSinceAccess            = -1
			tierToArchiveSinceMod            = -1
			tierToArchiveSinceAccess         = -1
			tierToArchiveSinceLastTierChange = -1
			deleteSinceMod                   = -1
			deleteSinceAccess                = -1
		)

		if props := input.BaseBlob.TierToCool; props != nil {
			if props.DaysAfterModificationGreaterThan != nil {
				tierToCoolSinceMod = int(*props.DaysAfterModificationGreaterThan)
			}
			if props.DaysAfterLastAccessTimeGreaterThan != nil {
				tierToCoolSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
			}
		}
		if props := input.BaseBlob.TierToArchive; props != nil {
			if props.DaysAfterModificationGreaterThan != nil {
				tierToArchiveSinceMod = int(*props.DaysAfterModificationGreaterThan)
			}
			if props.DaysAfterLastAccessTimeGreaterThan != nil {
				tierToArchiveSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
			}
			if props.DaysAfterLastTierChangeGreaterThan != nil {
				tierToArchiveSinceLastTierChange = int(*props.DaysAfterLastTierChangeGreaterThan)
			}
		}
		if props := input.BaseBlob.Delete; props != nil {
			if props.DaysAfterModificationGreaterThan != nil {
				deleteSinceMod = int(*props.DaysAfterModificationGreaterThan)
			}
			if props.DaysAfterLastAccessTimeGreaterThan != nil {
				deleteSinceAccess = int(*props.DaysAfterLastAccessTimeGreaterThan)
			}
		}

		actionModel.BaseBlob = []StorageAccountManagementPolicyRuleActionBaseBlobModel{
			{
				TierToCoolAfterDaysSinceLastAccessTimeGreaterThan:    tierToCoolSinceAccess,
				TierToCoolAfterDaysSinceModificationGreaterThan:      tierToCoolSinceMod,
				TierToArchiveAfterDaysSinceLastAccessTimeGreaterThan: tierToArchiveSinceAccess,
				TierToArchiveAfterDaysSinceLastTierChangeGreaterThan: tierToArchiveSinceLastTierChange,
				TierToArchiveAfterDaysSinceModificationGreaterThan:   tierToArchiveSinceMod,
				DeleteAfterDaysSinceLastAccessTimeGreaterThan:        deleteSinceAccess,
				DeleteAfterDaysSinceModificationGreaterThan:          deleteSinceMod,
			},
		}
	}

	if input.Snapshot != nil {
		deleteAfterCreation, archiveAfterCreation, archiveAfterLastTierChange, coolAfterCreation := -1, -1, -1, -1
		if input.Snapshot.Delete != nil && input.Snapshot.Delete.DaysAfterCreationGreaterThan != nil {
			deleteAfterCreation = int(*input.Snapshot.Delete.DaysAfterCreationGreaterThan)
		}
		if input.Snapshot.TierToArchive != nil && input.Snapshot.TierToArchive.DaysAfterCreationGreaterThan != nil {
			archiveAfterCreation = int(*input.Snapshot.TierToArchive.DaysAfterCreationGreaterThan)

			if v := input.Snapshot.TierToArchive.DaysAfterLastTierChangeGreaterThan; v != nil {
				archiveAfterLastTierChange = int(*v)
			}
		}
		if input.Snapshot.TierToCool != nil && input.Snapshot.TierToCool.DaysAfterCreationGreaterThan != nil {
			coolAfterCreation = int(*input.Snapshot.TierToCool.DaysAfterCreationGreaterThan)
		}
		actionModel.Snapshot = []StorageAccountManagementPolicyRuleActionSnapshotModel{
			{
				ChangeTierToArchiveAfterDaysSinceCreation:            archiveAfterCreation,
				ChangeTierToCoolAfterDaysSinceCreation:               coolAfterCreation,
				DeleteAfterDaysSinceCreationGreaterThan:              deleteAfterCreation,
				TierToArchiveAfterDaysSinceLastTierChangeGreaterThan: archiveAfterLastTierChange,
			},
		}
	}

	if input.Version != nil {
		deleteAfterCreation, archiveAfterCreation, archiveAfterLastTierChange, coolAfterCreation := -1, -1, -1, -1
		if input.Version.Delete != nil && input.Version.Delete.DaysAfterCreationGreaterThan != nil {
			deleteAfterCreation = int(*input.Version.Delete.DaysAfterCreationGreaterThan)
		}
		if input.Version.TierToArchive != nil && input.Version.TierToArchive.DaysAfterCreationGreaterThan != nil {
			archiveAfterCreation = int(*input.Version.TierToArchive.DaysAfterCreationGreaterThan)

			if v := input.Version.TierToArchive.DaysAfterLastTierChangeGreaterThan; v != nil {
				archiveAfterLastTierChange = int(*v)
			}
		}
		if input.Version.TierToCool != nil && input.Version.TierToCool.DaysAfterCreationGreaterThan != nil {
			coolAfterCreation = int(*input.Version.TierToCool.DaysAfterCreationGreaterThan)
		}
		actionModel.Version = []StorageAccountManagementPolicyRuleActionVersionModel{
			{
				ChangeTierToArchiveAfterDaysSinceCreation:            archiveAfterCreation,
				ChangeTierToCoolAfterDaysSinceCreation:               coolAfterCreation,
				DeleteAfterDaysSinceCreation:                         deleteAfterCreation,
				TierToArchiveAfterDaysSinceLastTierChangeGreaterThan: archiveAfterLastTierChange,
			},
		}
	}
	return []StorageAccountManagementPolicyRuleActionModel{actionModel}
}
