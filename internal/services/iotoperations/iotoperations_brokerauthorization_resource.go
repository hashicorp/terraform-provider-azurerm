package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BrokerAuthorizationResource struct{}

var _ sdk.ResourceWithUpdate = BrokerAuthorizationResource{}

type BrokerAuthorizationModel struct {
	Name                   string                                 `tfschema:"name"`
	ResourceGroupName      string                                 `tfschema:"resource_group_name"`
	InstanceName           string                                 `tfschema:"instance_name"`
	BrokerName             string                                 `tfschema:"broker_name"`
	AuthorizationPolicies  []BrokerAuthorizationPolicyModel       `tfschema:"authorization_policies"`
	ExtendedLocation       *ExtendedLocationModel                 `tfschema:"extended_location"`
	Tags                   map[string]string                      `tfschema:"tags"`
	ProvisioningState      *string                                `tfschema:"provisioning_state"`
}

type BrokerAuthorizationPolicyModel struct {
	Cache *string                           `tfschema:"cache"`
	Rules []BrokerAuthorizationRuleModel    `tfschema:"rules"`
}

type BrokerAuthorizationRuleModel struct {
	BrokerResources      []BrokerAuthorizationBrokerResourceModel     `tfschema:"broker_resources"`
	Principals           []BrokerAuthorizationPrincipalModel          `tfschema:"principals"`
	StateStoreResources  []BrokerAuthorizationStateStoreResourceModel `tfschema:"state_store_resources"`
}

type BrokerAuthorizationBrokerResourceModel struct {
	Method    string   `tfschema:"method"`
	ClientIds []string `tfschema:"client_ids"`
	Topics    []string `tfschema:"topics"`
}

type BrokerAuthorizationPrincipalModel struct {
	Attributes []map[string]string `tfschema:"attributes"`
	ClientIds  []string            `tfschema:"client_ids"`
	Usernames  []string            `tfschema:"usernames"`
}

type BrokerAuthorizationStateStoreResourceModel struct {
	KeyType string   `tfschema:"key_type"`
	Keys    []string `tfschema:"keys"`
	Method  string   `tfschema:"method"`
}

func (r BrokerAuthorizationResource) ModelObject() interface{} {
	return &BrokerAuthorizationModel{}
}

func (r BrokerAuthorizationResource) ResourceType() string {
	return "azurerm_iotoperations_broker_authorization"
}

func (r BrokerAuthorizationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return brokerauthorization.ValidateAuthorizationID
}

func (r BrokerAuthorizationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
		},
		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"broker_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"authorization_policies": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cache": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Enabled",
							"Disabled",
						}, false),
					},
					"rules": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"broker_resources": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"method": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Connect",
													"Publish",
													"Subscribe",
												}, false),
											},
											"client_ids": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"topics": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
								"principals": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"attributes": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeMap,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
											},
											"client_ids": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"usernames": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
								"state_store_resources": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"key_type": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Pattern",
													"Binary",
												}, false),
											},
											"keys": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"method": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Read",
													"Write",
													"ReadWrite",
												}, false),
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
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r BrokerAuthorizationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
			Computed: true,
		},
	}
}

func (r BrokerAuthorizationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthorizationClient

			var model BrokerAuthorizationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := brokerauthorization.NewAuthorizationID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.BrokerName, model.Name)

			// Build payload
			payload := brokerauthorization.BrokerAuthorizationResource{
				Properties: expandBrokerAuthorizationProperties(model.AuthorizationPolicies),
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = expandExtendedLocation(model.ExtendedLocation)
			}

			if len(model.Tags) > 0 {
				payload.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BrokerAuthorizationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthorizationClient

			id, err := brokerauthorization.ParseAuthorizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := BrokerAuthorizationModel{
				Name:              id.AuthorizationName,
				ResourceGroupName: id.ResourceGroupName,
				InstanceName:      id.InstanceName,
				BrokerName:        id.BrokerName,
			}

			if respModel := resp.Model; respModel != nil {
				if respModel.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(respModel.ExtendedLocation)
				}

				if respModel.Tags != nil {
					model.Tags = *respModel.Tags
				}

				if respModel.Properties != nil {
					model.AuthorizationPolicies = flattenBrokerAuthorizationProperties(respModel.Properties)
					
					if respModel.Properties.ProvisioningState != nil {
						provisioningState := string(*respModel.Properties.ProvisioningState)
						model.ProvisioningState = &provisioningState
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r BrokerAuthorizationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthorizationClient

			id, err := brokerauthorization.ParseAuthorizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BrokerAuthorizationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !metadata.ResourceData.HasChange("authorization_policies") {
				return nil
			}

			payload := brokerauthorization.BrokerAuthorizationPatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include authorization policies if they changed
			if metadata.ResourceData.HasChange("authorization_policies") {
				payload.Properties = &brokerauthorization.BrokerAuthorizerPropertiesPatch{
					AuthorizationPolicies: expandBrokerAuthorizationPolicies(model.AuthorizationPolicies),
				}
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BrokerAuthorizationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthorizationClient

			id, err := brokerauthorization.ParseAuthorizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

// Helper functions for expand/flatten operations
func expandBrokerAuthorizationProperties(policies []BrokerAuthorizationPolicyModel) *brokerauthorization.BrokerAuthorizerProperties {
	if len(policies) == 0 {
		return &brokerauthorization.BrokerAuthorizerProperties{}
	}

	result := &brokerauthorization.BrokerAuthorizerProperties{
		AuthorizationPolicies: expandBrokerAuthorizationPolicies(policies),
	}

	return result
}

func expandBrokerAuthorizationPolicies(policies []BrokerAuthorizationPolicyModel) *[]brokerauthorization.BrokerAuthorizerConfig {
	if len(policies) == 0 {
		return nil
	}

	result := make([]brokerauthorization.BrokerAuthorizerConfig, 0, len(policies))

	for _, policy := range policies {
		authPolicy := brokerauthorization.BrokerAuthorizerConfig{}

		if policy.Cache != nil {
			authPolicy.Cache = policy.Cache
		}

		if len(policy.Rules) > 0 {
			authPolicy.Rules = expandBrokerAuthorizationRules(policy.Rules)
		}

		result = append(result, authPolicy)
	}

	return &result
}

func expandBrokerAuthorizationRules(rules []BrokerAuthorizationRuleModel) *[]brokerauthorization.BrokerAuthorizerRule {
	if len(rules) == 0 {
		return nil
	}

	result := make([]brokerauthorization.BrokerAuthorizerRule, 0, len(rules))

	for _, rule := range rules {
		authRule := brokerauthorization.BrokerAuthorizerRule{}

		if len(rule.BrokerResources) > 0 {
			authRule.BrokerResources = expandBrokerAuthorizationBrokerResources(rule.BrokerResources)
		}

		if len(rule.Principals) > 0 {
			authRule.Principals = expandBrokerAuthorizationPrincipals(rule.Principals)
		}

		if len(rule.StateStoreResources) > 0 {
			authRule.StateStoreResources = expandBrokerAuthorizationStateStoreResources(rule.StateStoreResources)
		}

		result = append(result, authRule)
	}

	return &result
}

func expandBrokerAuthorizationBrokerResources(resources []BrokerAuthorizationBrokerResourceModel) *[]brokerauthorization.BrokerResourceRule {
	if len(resources) == 0 {
		return nil
	}

	result := make([]brokerauthorization.BrokerResourceRule, 0, len(resources))

	for _, resource := range resources {
		brokerResource := brokerauthorization.BrokerResourceRule{
			Method: resource.Method,
		}

		if len(resource.ClientIds) > 0 {
			brokerResource.ClientIds = &resource.ClientIds
		}

		if len(resource.Topics) > 0 {
			brokerResource.Topics = &resource.Topics
		}

		result = append(result, brokerResource)
	}

	return &result
}

func expandBrokerAuthorizationPrincipals(principals []BrokerAuthorizationPrincipalModel) *[]brokerauthorization.PrincipalDefinition {
	if len(principals) == 0 {
		return nil
	}

	result := make([]brokerauthorization.PrincipalDefinition, 0, len(principals))

	for _, principal := range principals {
		principalDef := brokerauthorization.PrincipalDefinition{}

		if len(principal.ClientIds) > 0 {
			principalDef.ClientIds = &principal.ClientIds
		}

		if len(principal.Usernames) > 0 {
			principalDef.Usernames = &principal.Usernames
		}

		// Note: Attributes field expansion would go here when API supports it

		result = append(result, principalDef)
	}

	return &result
}

func expandBrokerAuthorizationStateStoreResources(resources []BrokerAuthorizationStateStoreResourceModel) *[]brokerauthorization.StateStoreResourceRule {
	if len(resources) == 0 {
		return nil
	}

	result := make([]brokerauthorization.StateStoreResourceRule, 0, len(resources))

	for _, resource := range resources {
		stateStoreResource := brokerauthorization.StateStoreResourceRule{
			KeyType: resource.KeyType,
			Keys:    resource.Keys,
			Method:  resource.Method,
		}

		result = append(result, stateStoreResource)
	}

	return &result
}

func flattenBrokerAuthorizationProperties(props *brokerauthorization.BrokerAuthorizerProperties) []BrokerAuthorizationPolicyModel {
	if props == nil || props.AuthorizationPolicies == nil {
		return []BrokerAuthorizationPolicyModel{}
	}

	result := make([]BrokerAuthorizationPolicyModel, 0, len(*props.AuthorizationPolicies))

	for _, policy := range *props.AuthorizationPolicies {
		authPolicy := BrokerAuthorizationPolicyModel{}

		if policy.Cache != nil {
			authPolicy.Cache = policy.Cache
		}

		if policy.Rules != nil {
			authPolicy.Rules = flattenBrokerAuthorizationRules(*policy.Rules)
		}

		result = append(result, authPolicy)
	}

	return result
}

func flattenBrokerAuthorizationRules(rules []brokerauthorization.BrokerAuthorizerRule) []BrokerAuthorizationRuleModel {
	if len(rules) == 0 {
		return []BrokerAuthorizationRuleModel{}
	}

	result := make([]BrokerAuthorizationRuleModel, 0, len(rules))

	for _, rule := range rules {
		authRule := BrokerAuthorizationRuleModel{}

		if rule.BrokerResources != nil {
			authRule.BrokerResources = flattenBrokerAuthorizationBrokerResources(*rule.BrokerResources)
		}

		if rule.Principals != nil {
			authRule.Principals = flattenBrokerAuthorizationPrincipals(*rule.Principals)
		}

		if rule.StateStoreResources != nil {
			authRule.StateStoreResources = flattenBrokerAuthorizationStateStoreResources(*rule.StateStoreResources)
		}

		result = append(result, authRule)
	}

	return result
}

func flattenBrokerAuthorizationBrokerResources(resources []brokerauthorization.BrokerResourceRule) []BrokerAuthorizationBrokerResourceModel {
	result := make([]BrokerAuthorizationBrokerResourceModel, 0, len(resources))

	for _, resource := range resources {
		brokerResource := BrokerAuthorizationBrokerResourceModel{
			Method: resource.Method,
		}

		if resource.ClientIds != nil {
			brokerResource.ClientIds = *resource.ClientIds
		}

		if resource.Topics != nil {
			brokerResource.Topics = *resource.Topics
		}

		result = append(result, brokerResource)
	}

	return result
}

func flattenBrokerAuthorizationPrincipals(principals []brokerauthorization.PrincipalDefinition) []BrokerAuthorizationPrincipalModel {
	result := make([]BrokerAuthorizationPrincipalModel, 0, len(principals))

	for _, principal := range principals {
		principalModel := BrokerAuthorizationPrincipalModel{}

		if principal.ClientIds != nil {
			principalModel.ClientIds = *principal.ClientIds
		}

		if principal.Usernames != nil {
			principalModel.Usernames = *principal.Usernames
		}

		// Note: Attributes field flattening would go here when API supports it

		result = append(result, principalModel)
	}

	return result
}

func flattenBrokerAuthorizationStateStoreResources(resources []brokerauthorization.StateStoreResourceRule) []BrokerAuthorizationStateStoreResourceModel {
	result := make([]BrokerAuthorizationStateStoreResourceModel, 0, len(resources))

	for _, resource := range resources {
		stateStoreResource := BrokerAuthorizationStateStoreResourceModel{
			KeyType: resource.KeyType,
			Keys:    resource.Keys,
			Method:  resource.Method,
		}

		result = append(result, stateStoreResource)
	}

	return result
}