// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate        = EventGridPartnerNamespaceChannelResource{}
	_ sdk.ResourceWithCustomizeDiff = EventGridPartnerNamespaceChannelResource{}
)

type EventGridPartnerNamespaceChannelResource struct{}

type EventGridPartnerNamespaceChannelResourceModel struct {
	ChannelName                       string              `tfschema:"name"`
	PartnerNamespaceId                string              `tfschema:"partner_namespace_id"`
	ChannelType                       string              `tfschema:"channel_type"`
	ExpirationTimeIfNotActivatedInUtc string              `tfschema:"expiration_time_if_not_activated_in_utc"`
	PartnerTopic                      []PartnerTopicModel `tfschema:"partner_topic"`
	ReadinessState                    string              `tfschema:"readiness_state"`
}

type PartnerTopicModel struct {
	SubscriptionId       string                `tfschema:"subscription_id"`
	ResourceGroupName    string                `tfschema:"resource_group_name"`
	Name                 string                `tfschema:"name"`
	Source               string                `tfschema:"source"`
	EventTypeDefinitions []EventTypeDefinition `tfschema:"event_type_definitions"`
}

type EventTypeDefinition struct {
	InlineEventTypes []InlineEventTypeModel `tfschema:"inline_event_type"`
	Kind             string                 `tfschema:"kind"`
}

type InlineEventTypeModel struct {
	Name             string `tfschema:"name"`
	DataSchemaURL    string `tfschema:"data_schema_url"`
	Description      string `tfschema:"description"`
	DisplayName      string `tfschema:"display_name"`
	DocumentationURL string `tfschema:"documentation_url"`
}

// MessageForActivation is a problematic field as the API generates a custom default message that can be longer than the allowed length if not included.
// As such it has been excluded and left to the server to set the default.
// https://github.com/Azure/azure-rest-api-specs/issues/37280
func (EventGridPartnerNamespaceChannelResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
				"`name` must be between 3 and 50 characters. It can contain only letters, numbers and hyphens (-).",
			),
		},
		"partner_namespace_id": commonschema.ResourceIDReferenceRequiredForceNew(&partnernamespaces.PartnerNamespaceId{}),
		"channel_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(channels.ChannelTypePartnerTopic),
			ValidateFunc: validation.StringInSlice([]string{
				string(channels.ChannelTypePartnerTopic),
			}, false),
		},
		"expiration_time_if_not_activated_in_utc": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C due to api making a number of changes
			// - default is set to 7 days from creation
			// - if activated, this field is removed from the response
			Computed: true,
			ValidateFunc: validation.All(validation.IsRFC3339Time,
				validate.ExpirationTimeIfNotActivated(),
			),
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				if value, ok := d.GetOk("readiness_state"); ok && value.(string) == string(channels.ReadinessStateActivated) {
					return true
				}
				return false
			},
		},
		"partner_topic": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
							"`name` must be between 3 and 50 characters. It can contain only letters, numbers and hyphens (-).",
						),
					},
					"subscription_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsUUID,
					},
					"resource_group_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_type_definitions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"inline_event_type": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 128),
											},
											"display_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 128),
											},
											"data_schema_url": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.IsURLWithHTTPorHTTPS,
											},
											"description": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(0, 256),
											},
											"documentation_url": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.IsURLWithHTTPorHTTPS,
											},
										},
									},
								},
								"kind": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(channels.PossibleValuesForEventDefinitionKind(), false),
									Default:      channels.EventDefinitionKindInline,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (EventGridPartnerNamespaceChannelResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"readiness_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r EventGridPartnerNamespaceChannelResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config EventGridPartnerNamespaceChannelResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if config.ChannelType == string(channels.ChannelTypePartnerTopic) && config.PartnerTopic == nil {
				return fmt.Errorf("`partner_topic` is required when `channel_type` is `PartnerTopic`")
			}

			return nil
		},
	}
}

func (r EventGridPartnerNamespaceChannelResource) ModelObject() interface{} {
	return &EventGridPartnerNamespaceChannelResourceModel{}
}

func (EventGridPartnerNamespaceChannelResource) ResourceType() string {
	return "azurerm_eventgrid_partner_namespace_channel"
}

func (r EventGridPartnerNamespaceChannelResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config EventGridPartnerNamespaceChannelResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			partnerNamespaceId, err := partnernamespaces.ParsePartnerNamespaceID(config.PartnerNamespaceId)
			if err != nil {
				return err
			}

			id := channels.NewChannelID(subscriptionId, partnerNamespaceId.ResourceGroupName, partnerNamespaceId.PartnerNamespaceName, config.ChannelName)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := channels.Channel{
				Name: pointer.To(config.ChannelName),
				Properties: &channels.ChannelProperties{
					ChannelType:                     pointer.To(channels.ChannelType(config.ChannelType)),
					ExpirationTimeIfNotActivatedUtc: pointer.To(config.ExpirationTimeIfNotActivatedInUtc),
					PartnerTopicInfo:                expandPartnerNamespaceChannelPartnerTopic(config.PartnerTopic),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r EventGridPartnerNamespaceChannelResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels

			id, err := channels.ParseChannelID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EventGridPartnerNamespaceChannelResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := existing.Model

			// API generated default message can be longer than the allowed length so we have to clear it from update payload
			payload.Properties.MessageForActivation = nil

			if metadata.ResourceData.HasChange("expiration_time_if_not_activated_in_utc") && config.ReadinessState != string(channels.ReadinessStateActivated) {
				payload.Properties.ExpirationTimeIfNotActivatedUtc = pointer.To(config.ExpirationTimeIfNotActivatedInUtc)
			}

			if metadata.ResourceData.HasChange("partner_topic") {
				payload.Properties.PartnerTopicInfo = expandPartnerNamespaceChannelPartnerTopic(config.PartnerTopic)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r EventGridPartnerNamespaceChannelResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels

			id, err := channels.ParseChannelID(metadata.ResourceData.Id())
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

func (r EventGridPartnerNamespaceChannelResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return channels.ValidateChannelID
}

func (r EventGridPartnerNamespaceChannelResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels

			id, err := channels.ParseChannelID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			partnerNamespaceId := partnernamespaces.NewPartnerNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.PartnerNamespaceName)

			state := EventGridPartnerNamespaceChannelResourceModel{
				PartnerNamespaceId: partnerNamespaceId.ID(),
				ChannelName:        id.ChannelName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.ChannelType = pointer.FromEnum(props.ChannelType)
					state.ExpirationTimeIfNotActivatedInUtc = pointer.From(props.ExpirationTimeIfNotActivatedUtc)
					state.ReadinessState = pointer.FromEnum(props.ReadinessState)
					state.PartnerTopic = flattenPartnerNamespaceChannelPartnerTopic(props.PartnerTopicInfo)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func expandPartnerNamespaceChannelPartnerTopic(input []PartnerTopicModel) *channels.PartnerTopicInfo {
	if len(input) == 0 {
		return nil
	}

	partnerTopic := input[0]
	var eventTypeInfo *channels.EventTypeInfo

	if len(partnerTopic.EventTypeDefinitions) > 0 {
		eventTypeInfo = &channels.EventTypeInfo{
			InlineEventTypes: pointer.To(expandInlineEventTypes(partnerTopic.EventTypeDefinitions[0].InlineEventTypes)),
			Kind:             pointer.ToEnum[channels.EventDefinitionKind](partnerTopic.EventTypeDefinitions[0].Kind),
		}
	}

	return &channels.PartnerTopicInfo{
		AzureSubscriptionId: pointer.To(partnerTopic.SubscriptionId),
		EventTypeInfo:       eventTypeInfo,
		ResourceGroupName:   pointer.To(partnerTopic.ResourceGroupName),
		Name:                pointer.To(partnerTopic.Name),
		Source:              pointer.To(partnerTopic.Source),
	}
}

func flattenPartnerNamespaceChannelPartnerTopic(input *channels.PartnerTopicInfo) []PartnerTopicModel {
	if input == nil {
		return []PartnerTopicModel{}
	}

	var eventTypeDefinitions []EventTypeDefinition

	if input.EventTypeInfo != nil {
		eventTypeDefinitions = []EventTypeDefinition{
			{
				InlineEventTypes: flattenInlineEventTypes(pointer.From(input.EventTypeInfo.InlineEventTypes)),
				Kind:             pointer.FromEnum(input.EventTypeInfo.Kind),
			},
		}
	}

	return []PartnerTopicModel{
		{
			SubscriptionId:       pointer.From(input.AzureSubscriptionId),
			ResourceGroupName:    pointer.From(input.ResourceGroupName),
			Name:                 pointer.From(input.Name),
			Source:               pointer.From(input.Source),
			EventTypeDefinitions: eventTypeDefinitions,
		},
	}
}

func expandInlineEventTypes(inlineEvents []InlineEventTypeModel) map[string]channels.InlineEventProperties {
	if len(inlineEvents) == 0 {
		return nil
	}

	inlineEventsMap := make(map[string]channels.InlineEventProperties)

	for _, eventType := range inlineEvents {
		inlineEventsMap[eventType.Name] = channels.InlineEventProperties{
			DataSchemaURL:    pointer.To(eventType.DataSchemaURL),
			Description:      pointer.To(eventType.Description),
			DisplayName:      pointer.To(eventType.DisplayName),
			DocumentationURL: pointer.To(eventType.DocumentationURL),
		}
	}

	return inlineEventsMap
}

func flattenInlineEventTypes(inlineEventsMap map[string]channels.InlineEventProperties) []InlineEventTypeModel {
	if inlineEventsMap == nil {
		return []InlineEventTypeModel{}
	}

	inlineEventTypes := make([]InlineEventTypeModel, 0, len(inlineEventsMap))

	for name, properties := range inlineEventsMap {
		inlineEventTypes = append(inlineEventTypes, InlineEventTypeModel{
			Name:             name,
			DataSchemaURL:    pointer.From(properties.DataSchemaURL),
			Description:      pointer.From(properties.Description),
			DisplayName:      pointer.From(properties.DisplayName),
			DocumentationURL: pointer.From(properties.DocumentationURL),
		})
	}

	return inlineEventTypes
}
