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
	Name                  string                         `tfschema:"name"`
	ResourceGroupName     string                         `tfschema:"resource_group_name"`
	InstanceName          string                         `tfschema:"instance_name"`
	BrokerName            string                         `tfschema:"broker_name"`
	ExtendedLocationName  *string                        `tfschema:"extended_location_name"`
	ExtendedLocationType  *string                        `tfschema:"extended_location_type"`
	AuthorizationPolicies BrokerAuthorizationConfigModel `tfschema:"authorization_policies"`
	ProvisioningState     *string                        `tfschema:"provisioning_state"`
}

type BrokerAuthorizationConfigModel struct {
	Cache *string                        `tfschema:"cache"`
	Rules []BrokerAuthorizationRuleModel `tfschema:"rules"`
}

type BrokerAuthorizationRuleModel struct {
	BrokerResources     []BrokerAuthorizationBrokerResourceModel     `tfschema:"broker_resources"`
	Principals          BrokerAuthorizationPrincipalModel            `tfschema:"principals"`
	StateStoreResources []BrokerAuthorizationStateStoreResourceModel `tfschema:"state_store_resources"`
}

type BrokerAuthorizationBrokerResourceModel struct {
	Method  string   `tfschema:"method"`
	Clients []string `tfschema:"clients"`
	Topics  []string `tfschema:"topics"`
}

type BrokerAuthorizationPrincipalModel struct {
	Attributes []map[string]string `tfschema:"attributes"`
	Clients    []string            `tfschema:"clients"`
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
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "CustomLocation",
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
		"authorization_policies": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cache": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Enabled",
						ValidateFunc: validation.StringInSlice([]string{
							"Enabled",
							"Disabled",
						}, false),
					},
					"rules": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"broker_resources": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
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
											"clients": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
											"topics": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},
								"principals": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"clients": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
											"usernames": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
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
													"Binary",
													"Pattern",
													"String",
												}, false),
											},
											"keys": {
												Type:     pluginsdk.TypeList,
												Required: true,
												MinItems: 1,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
											"method": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Read",
													"ReadWrite",
													"Write",
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
	}
}

func (r BrokerAuthorizationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
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

			// Check if resource already exists
			existing, err := client.Get(ctx, id)
			if err == nil && existing.Model != nil {
				return fmt.Errorf("IoT Operations Broker Authorization %q already exists", id.AuthorizationName)
			}

			// Build payload
			payload := brokerauthorization.BrokerAuthorizationResource{
				ExtendedLocation: brokerauthorization.ExtendedLocation{
					Name: *model.ExtendedLocationName,
					Type: brokerauthorization.ExtendedLocationType(*model.ExtendedLocationType),
				},
				Properties: expandBrokerAuthorizationProperties(model.AuthorizationPolicies),
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
				model.ExtendedLocationName = &respModel.ExtendedLocation.Name
				extendedLocationType := string(respModel.ExtendedLocation.Type)
				model.ExtendedLocationType = &extendedLocationType

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

			// Since there's no separate Update method, use CreateOrUpdate
			payload := brokerauthorization.BrokerAuthorizationResource{
				ExtendedLocation: brokerauthorization.ExtendedLocation{
					Name: *model.ExtendedLocationName,
					Type: brokerauthorization.ExtendedLocationType(*model.ExtendedLocationType),
				},
				Properties: expandBrokerAuthorizationProperties(model.AuthorizationPolicies),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
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
func expandBrokerAuthorizationProperties(config BrokerAuthorizationConfigModel) *brokerauthorization.BrokerAuthorizationProperties {
	authConfig := brokerauthorization.AuthorizationConfig{
		Rules: expandBrokerAuthorizationRules(config.Rules),
	}

	if config.Cache != nil {
		cache := brokerauthorization.OperationalMode(*config.Cache)
		authConfig.Cache = &cache
	}

	return &brokerauthorization.BrokerAuthorizationProperties{
		AuthorizationPolicies: authConfig,
	}
}

func expandBrokerAuthorizationRules(rules []BrokerAuthorizationRuleModel) *[]brokerauthorization.AuthorizationRule {
	if len(rules) == 0 {
		return nil
	}

	result := make([]brokerauthorization.AuthorizationRule, 0, len(rules))

	for _, rule := range rules {
		authRule := brokerauthorization.AuthorizationRule{
			BrokerResources: expandBrokerAuthorizationBrokerResources(rule.BrokerResources),
			Principals:      expandBrokerAuthorizationPrincipals(rule.Principals),
		}

		if len(rule.StateStoreResources) > 0 {
			authRule.StateStoreResources = expandBrokerAuthorizationStateStoreResources(rule.StateStoreResources)
		}

		result = append(result, authRule)
	}

	return &result
}

func expandBrokerAuthorizationBrokerResources(resources []BrokerAuthorizationBrokerResourceModel) []brokerauthorization.BrokerResourceRule {
	result := make([]brokerauthorization.BrokerResourceRule, 0, len(resources))

	for _, resource := range resources {
		brokerResource := brokerauthorization.BrokerResourceRule{
			Method: brokerauthorization.BrokerResourceDefinitionMethods(resource.Method),
		}

		if len(resource.Clients) > 0 {
			brokerResource.ClientIds = &resource.Clients
		}

		if len(resource.Topics) > 0 {
			brokerResource.Topics = &resource.Topics
		}

		result = append(result, brokerResource)
	}

	return result
}

func expandBrokerAuthorizationPrincipals(principal BrokerAuthorizationPrincipalModel) brokerauthorization.PrincipalDefinition {
	result := brokerauthorization.PrincipalDefinition{}

	if len(principal.Clients) > 0 {
		result.ClientIds = &principal.Clients
	}

	if len(principal.Usernames) > 0 {
		result.Usernames = &principal.Usernames
	}

	if len(principal.Attributes) > 0 {
		result.Attributes = &principal.Attributes
	}

	return result
}

func expandBrokerAuthorizationStateStoreResources(resources []BrokerAuthorizationStateStoreResourceModel) *[]brokerauthorization.StateStoreResourceRule {
	result := make([]brokerauthorization.StateStoreResourceRule, 0, len(resources))

	for _, resource := range resources {
		stateStoreResource := brokerauthorization.StateStoreResourceRule{
			KeyType: brokerauthorization.StateStoreResourceKeyTypes(resource.KeyType),
			Keys:    resource.Keys,
			Method:  brokerauthorization.StateStoreResourceDefinitionMethods(resource.Method),
		}

		result = append(result, stateStoreResource)
	}

	return &result
}

func flattenBrokerAuthorizationProperties(props *brokerauthorization.BrokerAuthorizationProperties) BrokerAuthorizationConfigModel {
	result := BrokerAuthorizationConfigModel{}

	if props.AuthorizationPolicies.Cache != nil {
		cache := string(*props.AuthorizationPolicies.Cache)
		result.Cache = &cache
	}

	if props.AuthorizationPolicies.Rules != nil {
		result.Rules = flattenBrokerAuthorizationRules(*props.AuthorizationPolicies.Rules)
	}

	return result
}

func flattenBrokerAuthorizationRules(rules []brokerauthorization.AuthorizationRule) []BrokerAuthorizationRuleModel {
	result := make([]BrokerAuthorizationRuleModel, 0, len(rules))

	for _, rule := range rules {
		ruleModel := BrokerAuthorizationRuleModel{
			BrokerResources: flattenBrokerAuthorizationBrokerResources(rule.BrokerResources),
			Principals:      flattenBrokerAuthorizationPrincipals(rule.Principals),
		}

		if rule.StateStoreResources != nil {
			ruleModel.StateStoreResources = flattenBrokerAuthorizationStateStoreResources(*rule.StateStoreResources)
		}

		result = append(result, ruleModel)
	}

	return result
}

func flattenBrokerAuthorizationBrokerResources(resources []brokerauthorization.BrokerResourceRule) []BrokerAuthorizationBrokerResourceModel {
	result := make([]BrokerAuthorizationBrokerResourceModel, 0, len(resources))

	for _, resource := range resources {
		brokerResource := BrokerAuthorizationBrokerResourceModel{
			Method: string(resource.Method),
		}

		if resource.ClientIds != nil {
			brokerResource.Clients = *resource.ClientIds
		}

		if resource.Topics != nil {
			brokerResource.Topics = *resource.Topics
		}

		result = append(result, brokerResource)
	}

	return result
}

func flattenBrokerAuthorizationPrincipals(principal brokerauthorization.PrincipalDefinition) BrokerAuthorizationPrincipalModel {
	result := BrokerAuthorizationPrincipalModel{}

	if principal.ClientIds != nil {
		result.Clients = *principal.ClientIds
	}

	if principal.Usernames != nil {
		result.Usernames = *principal.Usernames
	}

	if principal.Attributes != nil {
		result.Attributes = *principal.Attributes
	}

	return result
}

func flattenBrokerAuthorizationStateStoreResources(resources []brokerauthorization.StateStoreResourceRule) []BrokerAuthorizationStateStoreResourceModel {
	result := make([]BrokerAuthorizationStateStoreResourceModel, 0, len(resources))

	for _, resource := range resources {
		stateStoreResource := BrokerAuthorizationStateStoreResourceModel{
			KeyType: string(resource.KeyType),
			Keys:    resource.Keys,
			Method:  string(resource.Method),
		}

		result = append(result, stateStoreResource)
	}

	return result
}
