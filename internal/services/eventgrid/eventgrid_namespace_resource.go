// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EventGridNamespaceResource{}

type EventGridNamespaceResource struct{}

type EventGridNamespaceResourceModel struct {
	Name                     string                                     `tfschema:"name"`
	Location                 string                                     `tfschema:"location"`
	ResourceGroup            string                                     `tfschema:"resource_group_name"`
	Capacity                 int64                                      `tfschema:"capacity"`
	InboundIpRules           []InboundIpRuleModel                       `tfschema:"inbound_ip_rule"`
	Identity                 []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccess      string                                     `tfschema:"public_network_access"`
	Sku                      string                                     `tfschema:"sku"`
	TopicSpacesConfiguration []TopicSpacesConfigurationModel            `tfschema:"topic_spaces_configuration"`
	Tags                     map[string]string                          `tfschema:"tags"`
}

type InboundIpRuleModel struct {
	IpMask string `tfschema:"ip_mask"`
	Action string `tfschema:"action"`
}

type TopicSpacesConfigurationModel struct {
	AlternativeAuthenticationNameSources       []string                 `tfschema:"alternative_authentication_name_source"`
	MaximumClientSessionsPerAuthenticationName int64                    `tfschema:"maximum_client_sessions_per_authentication_name"`
	MaximumSessionExpiryInHours                int64                    `tfschema:"maximum_session_expiry_in_hours"`
	RouteTopicResourceId                       string                   `tfschema:"route_topic_id"`
	DynamicRoutingEnrichment                   []RoutingEnrichmentModel `tfschema:"dynamic_routing_enrichment"`
	StaticRoutingEnrichment                    []RoutingEnrichmentModel `tfschema:"static_routing_enrichment"`
}

type RoutingEnrichmentModel struct {
	Key   string `tfschema:"key"`
	Value string `tfschema:"value"`
}

func (r EventGridNamespaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 50),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"capacity": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 40),
			Default:      1,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"inbound_ip_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 128,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_mask": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.CIDR,
					},
					"action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(namespaces.IPActionTypeAllow),
						ValidateFunc: validation.StringInSlice([]string{
							string(namespaces.IPActionTypeAllow),
						}, false),
					},
				},
			},
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.PublicNetworkAccessEnabled),
				string(namespaces.PublicNetworkAccessDisabled),
			}, false),
			Default: string(namespaces.PublicNetworkAccessEnabled),
		},

		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
			Default:      namespaces.SkuNameStandard,
		},

		"topic_spaces_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"alternative_authentication_name_source": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForAlternativeAuthenticationNameSource(), false),
						},
					},

					"maximum_client_sessions_per_authentication_name": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(1, 100),
					},

					"maximum_session_expiry_in_hours": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(1, 8),
					},

					"route_topic_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: topics.ValidateTopicID,
					},

					"dynamic_routing_enrichment": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 20),
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 128),
								},
							},
						},
					},

					"static_routing_enrichment": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 20),
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 128),
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r EventGridNamespaceResource) ModelObject() interface{} {
	return &EventGridNamespaceResourceModel{}
}

func (r EventGridNamespaceResource) ResourceType() string {
	return "azurerm_eventgrid_namespace"
}

func (r EventGridNamespaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r EventGridNamespaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model EventGridNamespaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := namespaces.NewNamespaceID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			namespace := namespaces.Namespace{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Name:     pointer.To(model.Name),
				Properties: &namespaces.NamespaceProperties{
					InboundIPRules:      expandInboundIPRules(model.InboundIpRules),
					PublicNetworkAccess: pointer.To(namespaces.PublicNetworkAccess(model.PublicNetworkAccess)),
				},
				Sku: &namespaces.NamespaceSku{
					Capacity: pointer.To(model.Capacity),
					Name:     pointer.To(namespaces.SkuName(model.Sku)),
				},
				Tags: pointer.To(model.Tags),
			}

			if len(model.TopicSpacesConfiguration) > 0 {
				namespace.Properties.TopicSpacesConfiguration = expandTopicSpacesConfiguration(model.TopicSpacesConfiguration)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, namespace); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EventGridNamespaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient

			var model EventGridNamespaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			payload := namespaces.NamespaceUpdateParameters{
				Properties: &namespaces.NamespaceUpdateParameterProperties{},
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = identity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("capacity") {
				payload.Sku = &namespaces.NamespaceSku{
					Capacity: pointer.To(model.Capacity),
				}
			}

			if metadata.ResourceData.HasChange("inbound_ip_rule") {
				payload.Properties.InboundIPRules = expandInboundIPRules(model.InboundIpRules)
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.To(namespaces.PublicNetworkAccess(model.PublicNetworkAccess))
			}

			if metadata.ResourceData.HasChange("topic_spaces_configuration") {
				payload.Properties.TopicSpacesConfiguration = expandTopicSpacesConfigurationUpdate(model.TopicSpacesConfiguration)
			}

			if err = client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", *id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EventGridNamespaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := EventGridNamespaceResourceModel{
				Name:          id.NamespaceName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if model.Sku != nil {
					state.Sku = string(pointer.From(model.Sku.Name))
					state.Capacity = pointer.From(model.Sku.Capacity)
				}
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(flattenedIdentity)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					topicSpacesConfig, err := flattenTopicSpacesConfiguration(props.TopicSpacesConfiguration)
					if err != nil {
						return fmt.Errorf("flattening `topic_spaces_configuration`: %v", err)
					}
					state.TopicSpacesConfiguration = topicSpacesConfig
					state.InboundIpRules = flattenInboundIPRules(props.InboundIPRules)
					state.PublicNetworkAccess = string(pointer.From(props.PublicNetworkAccess))
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r EventGridNamespaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}

			return nil
		},
	}
}

func (r EventGridNamespaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespaces.ValidateNamespaceID
}

func expandInboundIPRules(input []InboundIpRuleModel) *[]namespaces.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	ipRules := make([]namespaces.InboundIPRule, 0)
	for _, v := range input {
		ipRules = append(ipRules, namespaces.InboundIPRule{
			Action: pointer.To(namespaces.IPActionType(v.Action)),
			IPMask: pointer.To(v.IpMask),
		})
	}
	return &ipRules
}

func flattenInboundIPRules(ipRules *[]namespaces.InboundIPRule) []InboundIpRuleModel {
	output := make([]InboundIpRuleModel, 0)

	if ipRules == nil || len(*ipRules) == 0 {
		return output
	}

	for _, v := range *ipRules {
		output = append(output, InboundIpRuleModel{
			IpMask: pointer.From(v.IPMask),
			Action: string(pointer.From(v.Action)),
		})
	}
	return output
}

func expandTopicSpacesConfiguration(input []TopicSpacesConfigurationModel) *namespaces.TopicSpacesConfiguration {
	if input == nil {
		return nil
	}
	topicSpacesConfig := namespaces.TopicSpacesConfiguration{
		State: pointer.To(namespaces.TopicSpacesConfigurationStateEnabled),
		ClientAuthentication: &namespaces.ClientAuthenticationSettings{
			AlternativeAuthenticationNameSources: expandAlternativeAuthenticationNameSources(input[0].AlternativeAuthenticationNameSources),
		},
		MaximumClientSessionsPerAuthenticationName: pointer.To(input[0].MaximumClientSessionsPerAuthenticationName),
		MaximumSessionExpiryInHours:                pointer.To(input[0].MaximumSessionExpiryInHours),
		RouteTopicResourceId:                       pointer.To(input[0].RouteTopicResourceId),
		RoutingEnrichments: &namespaces.RoutingEnrichments{
			Dynamic: expandDynamicRoutingEnrichments(input[0].DynamicRoutingEnrichment),
			Static:  expandStaticRoutingEnrichments(input[0].StaticRoutingEnrichment),
		},
	}
	return &topicSpacesConfig
}

func expandTopicSpacesConfigurationUpdate(input []TopicSpacesConfigurationModel) *namespaces.UpdateTopicSpacesConfigurationInfo {
	if input == nil {
		return nil
	}

	topicSpacesConfig := namespaces.UpdateTopicSpacesConfigurationInfo{
		State: pointer.To(namespaces.TopicSpacesConfigurationStateEnabled),
		ClientAuthentication: &namespaces.ClientAuthenticationSettings{
			AlternativeAuthenticationNameSources: expandAlternativeAuthenticationNameSources(input[0].AlternativeAuthenticationNameSources),
		},
		MaximumClientSessionsPerAuthenticationName: pointer.To(input[0].MaximumClientSessionsPerAuthenticationName),
		MaximumSessionExpiryInHours:                pointer.To(input[0].MaximumSessionExpiryInHours),
		RouteTopicResourceId:                       pointer.To(input[0].RouteTopicResourceId),
		RoutingEnrichments: &namespaces.RoutingEnrichments{
			Dynamic: expandDynamicRoutingEnrichments(input[0].DynamicRoutingEnrichment),
			Static:  expandStaticRoutingEnrichments(input[0].StaticRoutingEnrichment),
		},
	}
	return &topicSpacesConfig
}

func expandAlternativeAuthenticationNameSources(input []string) *[]namespaces.AlternativeAuthenticationNameSource {
	if len(input) == 0 {
		return nil
	}

	nameSources := make([]namespaces.AlternativeAuthenticationNameSource, 0)
	for _, v := range input {
		nameSources = append(nameSources, namespaces.AlternativeAuthenticationNameSource(v))
	}
	return &nameSources
}

func expandDynamicRoutingEnrichments(input []RoutingEnrichmentModel) *[]namespaces.DynamicRoutingEnrichment {
	if len(input) == 0 {
		return nil
	}
	dynamicRoutingEnrichments := make([]namespaces.DynamicRoutingEnrichment, 0)
	for _, v := range input {
		dynamicRoutingEnrichments = append(dynamicRoutingEnrichments, namespaces.DynamicRoutingEnrichment{
			Value: pointer.To(v.Value),
			Key:   pointer.To(v.Key),
		})
	}
	return &dynamicRoutingEnrichments
}

func expandStaticRoutingEnrichments(input []RoutingEnrichmentModel) *[]namespaces.StaticRoutingEnrichment {
	if len(input) == 0 {
		return nil
	}

	staticRoutingEnrichments := make([]namespaces.StaticRoutingEnrichment, 0)
	for _, v := range input {
		staticRoutingEnrichments = append(staticRoutingEnrichments, namespaces.StaticStringRoutingEnrichment{
			Value:     pointer.To(v.Value),
			Key:       pointer.To(v.Key),
			ValueType: namespaces.StaticRoutingEnrichmentTypeString,
		})
	}
	return &staticRoutingEnrichments
}

func flattenTopicSpacesConfiguration(topicSpacesConfig *namespaces.TopicSpacesConfiguration) ([]TopicSpacesConfigurationModel, error) {
	var output TopicSpacesConfigurationModel
	if topicSpacesConfig == nil {
		return nil, nil
	}

	output.MaximumSessionExpiryInHours = pointer.From(topicSpacesConfig.MaximumSessionExpiryInHours)
	output.MaximumClientSessionsPerAuthenticationName = pointer.From(topicSpacesConfig.MaximumClientSessionsPerAuthenticationName)
	var routeId string
	if topicSpacesConfig.RouteTopicResourceId != nil && *topicSpacesConfig.RouteTopicResourceId != "" {
		id, err := topics.ParseTopicID(*topicSpacesConfig.RouteTopicResourceId)
		if err != nil {
			return nil, err
		}
		routeId = id.ID()
	}
	output.RouteTopicResourceId = routeId

	if topicSpacesConfig.ClientAuthentication != nil {
		output.AlternativeAuthenticationNameSources = flattenAlternativeAuthenticationNameSources(topicSpacesConfig.ClientAuthentication.AlternativeAuthenticationNameSources)
	}
	if topicSpacesConfig.RoutingEnrichments != nil {
		output.DynamicRoutingEnrichment = flattenDynamicRoutingEnrichments(topicSpacesConfig.RoutingEnrichments.Dynamic)
		output.StaticRoutingEnrichment = flattenStaticRoutingEnrichments(topicSpacesConfig.RoutingEnrichments.Static)
	}

	return []TopicSpacesConfigurationModel{output}, nil
}

func flattenAlternativeAuthenticationNameSources(nameSources *[]namespaces.AlternativeAuthenticationNameSource) []string {
	output := make([]string, 0)

	if nameSources == nil || len(*nameSources) == 0 {
		return output
	}

	for _, v := range *nameSources {
		output = append(output, string(v))
	}
	return output
}

func flattenDynamicRoutingEnrichments(dynamicRoutingEnrichments *[]namespaces.DynamicRoutingEnrichment) []RoutingEnrichmentModel {
	output := make([]RoutingEnrichmentModel, 0)
	if dynamicRoutingEnrichments == nil || len(*dynamicRoutingEnrichments) == 0 {
		return output
	}

	for _, v := range *dynamicRoutingEnrichments {
		output = append(output, RoutingEnrichmentModel{
			Value: pointer.From(v.Value),
			Key:   pointer.From(v.Key),
		})
	}
	return output
}

func flattenStaticRoutingEnrichments(staticRoutingEnrichments *[]namespaces.StaticRoutingEnrichment) []RoutingEnrichmentModel {
	output := make([]RoutingEnrichmentModel, 0)
	if staticRoutingEnrichments == nil || len(*staticRoutingEnrichments) == 0 {
		return output
	}

	for _, v := range *staticRoutingEnrichments {
		output = append(output, RoutingEnrichmentModel{
			Value: pointer.From(v.(namespaces.StaticStringRoutingEnrichment).Value),
			Key:   pointer.From(v.(namespaces.StaticStringRoutingEnrichment).Key),
		})
	}
	return output
}
