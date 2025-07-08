// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespacetopics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EventGridNamespaceTopicResource{}

type EventGridNamespaceTopicResource struct{}

type EventGridNamespaceTopicResourceModel struct {
	Name                 string `tfschema:"name"`
	NamespaceId          string `tfschema:"namespace_id"`
	EventRetentionInDays int64  `tfschema:"event_retention_in_days"`
	InputSchema          string `tfschema:"input_schema"`
}

func (r EventGridNamespaceTopicResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 50), // ^[a-zA-Z0-9-]*$
		},

		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespacetopics.ValidateNamespaceID,
		},

		"event_retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 7),
		},

		"input_schema": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "CloudEventSchemaV1_0",
			ValidateFunc: validation.StringInSlice([]string{
				"CloudEventSchemaV1_0",
			}, false),
		},
	}
}

func (r EventGridNamespaceTopicResource) ModelObject() interface{} {
	return &EventGridNamespaceTopicResourceModel{}
}

func (r EventGridNamespaceTopicResource) ResourceType() string {
	return "azurerm_eventgrid_namespace_topic"
}

func (r EventGridNamespaceTopicResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r EventGridNamespaceTopicResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient

			var model EventGridNamespaceTopicResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			subscriptionId := ""
			namespaceName := ""
			resourceGroupName := ""
			if v := model.NamespaceId; v != "" {
				namespaceId, err := namespaces.ParseNamespaceID(v)
				if err != nil {
					return err
				}
				subscriptionId = namespaceId.SubscriptionId
				namespaceName = namespaceId.NamespaceName
				resourceGroupName = namespaceId.ResourceGroupName
			}

			id := namespacetopics.NewNamespaceTopicID(subscriptionId, resourceGroupName, namespaceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			namespaceTopic := namespacetopics.NamespaceTopic{
				Name: pointer.To(model.Name),
				Properties: &namespacetopics.NamespaceTopicProperties{
					EventRetentionInDays: pointer.To(model.EventRetentionInDays),
					InputSchema:          pointer.To(namespacetopics.EventInputSchema(model.InputSchema)),
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, namespaceTopic); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EventGridNamespaceTopicResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient

			var model EventGridNamespaceTopicResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := namespacetopics.ParseNamespaceTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			payload := namespacetopics.NamespaceTopicUpdateParameters{
				Properties: &namespacetopics.NamespaceTopicUpdateParameterProperties{},
			}

			if metadata.ResourceData.HasChange("event_retention_in_days") {
				payload.Properties.EventRetentionInDays = pointer.To(model.EventRetentionInDays)
			}

			if err = client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", *id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EventGridNamespaceTopicResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient

			id, err := namespacetopics.ParseNamespaceTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := EventGridNamespaceTopicResourceModel{
				Name:        id.TopicName,
				NamespaceId: namespacetopics.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.EventRetentionInDays = int64(pointer.From(props.EventRetentionInDays))
					state.InputSchema = string(pointer.From(props.InputSchema))
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r EventGridNamespaceTopicResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.NamespaceTopicsClient

			id, err := namespacetopics.ParseNamespaceTopicID(metadata.ResourceData.Id())
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

func (r EventGridNamespaceTopicResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespacetopics.ValidateNamespaceTopicID
}
