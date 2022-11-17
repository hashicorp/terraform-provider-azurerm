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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageAccountManagementPolicyRuleResource struct{}

var _ sdk.ResourceWithUpdate = StorageAccountManagementPolicyRuleResource{}

type StorageAccountManagementPolicyRuleModel struct {
	Name               string `tfschema:"name"`
	ManagementPolicyId string `tfschema:"management_policy_id"`
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
		"filters": {
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
	}
}

func (r StorageAccountManagementPolicyRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageAccountManagementPolicyRuleResource) ResourceType() string {
	return "azurerm_storage_account_management_policy_rule"
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
			if r.getRuleByName(existing, plan.Name) != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			rule, err := r.expandRule(plan)
			if err != nil {
				return fmt.Errorf("expanding rule: %v", err)
			}

			r.appendRuleToPolicy(&existing, *rule)

			// TODO: construct params from model

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.StorageAccountName, params); err != nil {
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

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			// TODO: construct the params
			model := StorageAccountManagementPolicyRuleModel{}

			model.Tags = tags.Flatten(existing.Tags)

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

			params, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			// TODO: update the params
			// if props := params.Properties; props != nil {
			// 	if metadata.ResourceData.HasChange("xxx") {
			// 		props.Xxx = plan.Xxx
			// 	}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, params)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
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

			future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r StorageAccountManagementPolicyRuleResource) getRuleByName(policy storage.ManagementPolicy, name string) *storage.ManagementPolicyRule {
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

func (r StorageAccountManagementPolicyRuleResource) appendRuleToPolicy(policy *storage.ManagementPolicy, rule storage.ManagementPolicyRule) {
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

func (r StorageAccountManagementPolicyRuleResource) expandRule(input StorageAccountManagementPolicyRuleModel) (*storage.ManagementPolicyRule, error) {

}
