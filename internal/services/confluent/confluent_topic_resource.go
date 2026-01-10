// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/topicrecords"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confluent/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ConfluentTopicResource struct{}

type ConfluentTopicResourceModel struct {
	TopicName         string                       `tfschema:"topic_name"`
	ClusterId         string                       `tfschema:"cluster_id"`
	EnvironmentId     string                       `tfschema:"environment_id"`
	OrganizationId    string                       `tfschema:"organization_id"`
	ResourceGroupName string                       `tfschema:"resource_group_name"`
	PartitionsCount   string                       `tfschema:"partitions_count"`
	ReplicationFactor string                       `tfschema:"replication_factor"`
	Configs           []ConfluentTopicConfigModel  `tfschema:"configs"`

	// Computed
	Id       string                         `tfschema:"id"`
	TopicId  string                         `tfschema:"topic_id"`
	Kind     string                         `tfschema:"kind"`
	Metadata []ConfluentTopicMetadataModel  `tfschema:"metadata"`
}

type ConfluentTopicConfigModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type ConfluentTopicMetadataModel struct {
	Self         string `tfschema:"self"`
	ResourceName string `tfschema:"resource_name"`
}

func (r ConfluentTopicResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"topic_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"organization_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"partitions_count": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"replication_factor": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"configs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r ConfluentTopicResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"topic_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"self": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"resource_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r ConfluentTopicResource) ModelObject() interface{} {
	return &ConfluentTopicResourceModel{}
}

func (r ConfluentTopicResource) ResourceType() string {
	return "azurerm_confluent_topic"
}

func (r ConfluentTopicResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.TopicClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ConfluentTopicResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := topicrecords.NewTopicID(subscriptionId, model.ResourceGroupName, model.OrganizationId, model.EnvironmentId, model.ClusterId, model.TopicName)

			existing, err := client.TopicsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_confluent_topic", id.ID())
			}

			payload := topicrecords.TopicRecord{
				Name:       pointer.To(model.TopicName),
				Properties: expandConfluentTopicProperties(model),
			}

			if _, err := client.TopicsCreate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ConfluentTopicResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.TopicClient

			id, err := topicrecords.ParseTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.TopicsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ConfluentTopicResourceModel
			state.TopicName = id.TopicName
			state.ClusterId = id.ClusterId
			state.EnvironmentId = id.EnvironmentId
			state.OrganizationId = id.OrganizationName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Id = pointer.From(model.Id)

				if props := model.Properties; props != nil {
					state.TopicId = pointer.From(props.TopicId)
					state.Kind = pointer.From(props.Kind)
					state.PartitionsCount = pointer.From(props.PartitionsCount)
					state.ReplicationFactor = pointer.From(props.ReplicationFactor)
					state.Metadata = flattenConfluentTopicMetadata(props.Metadata)
					state.Configs = flattenConfluentTopicConfigs(props.InputConfigs)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ConfluentTopicResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.TopicClient

			id, err := topicrecords.ParseTopicID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.TopicsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConfluentTopicResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.TopicID
}

func expandConfluentTopicProperties(model ConfluentTopicResourceModel) *topicrecords.TopicProperties {
	props := &topicrecords.TopicProperties{}

	if model.PartitionsCount != "" {
		props.PartitionsCount = pointer.To(model.PartitionsCount)
	}

	if model.ReplicationFactor != "" {
		props.ReplicationFactor = pointer.To(model.ReplicationFactor)
	}

	if len(model.Configs) > 0 {
		configs := make([]topicrecords.TopicsInputConfig, 0)
		for _, cfg := range model.Configs {
			configs = append(configs, topicrecords.TopicsInputConfig{
				Name:  pointer.To(cfg.Name),
				Value: pointer.To(cfg.Value),
			})
		}
		props.InputConfigs = &configs
	}

	return props
}

func flattenConfluentTopicMetadata(input *topicrecords.TopicMetadataEntity) []ConfluentTopicMetadataModel {
	if input == nil {
		return []ConfluentTopicMetadataModel{}
	}

	return []ConfluentTopicMetadataModel{
		{
			Self:         pointer.From(input.Self),
			ResourceName: pointer.From(input.ResourceName),
		},
	}
}

func flattenConfluentTopicConfigs(input *[]topicrecords.TopicsInputConfig) []ConfluentTopicConfigModel {
	if input == nil || len(*input) == 0 {
		return []ConfluentTopicConfigModel{}
	}

	result := make([]ConfluentTopicConfigModel, 0)
	for _, cfg := range *input {
		result = append(result, ConfluentTopicConfigModel{
			Name:  pointer.From(cfg.Name),
			Value: pointer.From(cfg.Value),
		})
	}

	return result
}
