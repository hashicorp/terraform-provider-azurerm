// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scenvironmentrecords"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confluent/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ConfluentEnvironmentResource struct{}

type ConfluentEnvironmentResourceModel struct {
	EnvironmentId     string                                     `tfschema:"environment_id"`
	OrganizationId    string                                     `tfschema:"organization_id"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	DisplayName       string                                     `tfschema:"display_name"`
	StreamGovernance  []ConfluentStreamGovernanceModel           `tfschema:"stream_governance"`

	// Computed
	Id               string                       `tfschema:"id"`
	Kind             string                       `tfschema:"kind"`
	Metadata         []ConfluentMetadataModel     `tfschema:"metadata"`
}

type ConfluentStreamGovernanceModel struct {
	Package string `tfschema:"package"`
}

type ConfluentMetadataModel struct {
	Self             string `tfschema:"self"`
	ResourceName     string `tfschema:"resource_name"`
	CreatedTimestamp string `tfschema:"created_timestamp"`
	UpdatedTimestamp string `tfschema:"updated_timestamp"`
	DeletedTimestamp string `tfschema:"deleted_timestamp"`
}

func (r ConfluentEnvironmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stream_governance": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"package": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(scenvironmentrecords.PossibleValuesForPackage(), false),
					},
				},
			},
		},
	}
}

func (r ConfluentEnvironmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
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

					"created_timestamp": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"updated_timestamp": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"deleted_timestamp": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r ConfluentEnvironmentResource) ModelObject() interface{} {
	return &ConfluentEnvironmentResourceModel{}
}

func (r ConfluentEnvironmentResource) ResourceType() string {
	return "azurerm_confluent_environment"
}

func (r ConfluentEnvironmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.EnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ConfluentEnvironmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := scenvironmentrecords.NewEnvironmentID(subscriptionId, model.ResourceGroupName, model.OrganizationId, model.EnvironmentId)

			existing, err := client.OrganizationGetEnvironmentById(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_confluent_environment", id.ID())
			}

			payload := scenvironmentrecords.SCEnvironmentRecord{
				Kind:       pointer.To("Environment"),
				Properties: expandConfluentEnvironmentProperties(model),
			}

			if model.DisplayName != "" {
				payload.Name = pointer.To(model.DisplayName)
			}

			if _, err := client.EnvironmentCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ConfluentEnvironmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.EnvironmentClient

			id, err := scenvironmentrecords.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.OrganizationGetEnvironmentById(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ConfluentEnvironmentResourceModel
			state.EnvironmentId = id.EnvironmentId
			state.OrganizationId = id.OrganizationName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Id = pointer.From(model.Id)
				state.Kind = pointer.From(model.Kind)
				state.DisplayName = pointer.From(model.Name)

				if props := model.Properties; props != nil {
					state.Metadata = flattenConfluentMetadata(props.Metadata)
					state.StreamGovernance = flattenConfluentStreamGovernance(props.StreamGovernanceConfig)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ConfluentEnvironmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.EnvironmentClient

			id, err := scenvironmentrecords.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.EnvironmentDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConfluentEnvironmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.EnvironmentID
}

func expandConfluentEnvironmentProperties(model ConfluentEnvironmentResourceModel) *scenvironmentrecords.EnvironmentProperties {
	props := &scenvironmentrecords.EnvironmentProperties{}

	if len(model.StreamGovernance) > 0 {
		props.StreamGovernanceConfig = &scenvironmentrecords.StreamGovernanceConfig{
			Package: pointer.To(scenvironmentrecords.Package(model.StreamGovernance[0].Package)),
		}
	}

	return props
}

func flattenConfluentMetadata(input *scenvironmentrecords.SCMetadataEntity) []ConfluentMetadataModel {
	if input == nil {
		return []ConfluentMetadataModel{}
	}

	return []ConfluentMetadataModel{
		{
			Self:             pointer.From(input.Self),
			ResourceName:     pointer.From(input.ResourceName),
			CreatedTimestamp: pointer.From(input.CreatedTimestamp),
			UpdatedTimestamp: pointer.From(input.UpdatedTimestamp),
			DeletedTimestamp: pointer.From(input.DeletedTimestamp),
		},
	}
}

func flattenConfluentStreamGovernance(input *scenvironmentrecords.StreamGovernanceConfig) []ConfluentStreamGovernanceModel {
	if input == nil || input.Package == nil {
		return []ConfluentStreamGovernanceModel{}
	}

	return []ConfluentStreamGovernanceModel{
		{
			Package: string(pointer.From(input.Package)),
		},
	}
}
