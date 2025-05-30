// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpacecraftResource struct{}

var _ sdk.ResourceWithUpdate = SpacecraftResource{}

var _ sdk.ResourceWithDeprecationAndNoReplacement = SpacecraftResource{}

func (r SpacecraftResource) DeprecationMessage() string {
	return "The `azurerm_orbital_spacecraft` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider."
}

type SpacecraftResourceModel struct {
	Name            string                `tfschema:"name"`
	ResourceGroup   string                `tfschema:"resource_group_name"`
	Location        string                `tfschema:"location"`
	NoradId         string                `tfschema:"norad_id"`
	Links           []SpacecraftLinkModel `tfschema:"links"`
	TitleLine       string                `tfschema:"title_line"`
	TwoLineElements []string              `tfschema:"two_line_elements"`
	Tags            map[string]string     `tfschema:"tags"`
}

func (r SpacecraftResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"norad_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(5, 5),
		},

		"links": SpacecraftLinkSchema(),

		"two_line_elements": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 2,
			MaxItems: 2,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringLenBetween(69, 69),
			},
		},

		"title_line": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": tags.Schema(),
	}
}

func (r SpacecraftResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SpacecraftResource) ModelObject() interface{} {
	return &SpacecraftResourceModel{}
}

func (r SpacecraftResource) ResourceType() string {
	return "azurerm_orbital_spacecraft"
}

func (r SpacecraftResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpacecraftResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Orbital.SpacecraftClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := spacecraft.NewSpacecraftID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			links, err := expandSpacecraftLinks(model.Links)
			if err != nil {
				return fmt.Errorf("expanding `links`: %+v", err)
			}
			spacecraftProperties := spacecraft.SpacecraftsProperties{
				Links:     links,
				NoradId:   utils.String(model.NoradId),
				TleLine1:  model.TwoLineElements[0],
				TleLine2:  model.TwoLineElements[1],
				TitleLine: model.TitleLine,
			}

			spacecraft := spacecraft.Spacecraft{
				Id:         utils.String(id.ID()),
				Location:   model.Location,
				Name:       utils.String(model.Name),
				Properties: spacecraftProperties,
				Tags:       &model.Tags,
			}
			if err = client.CreateOrUpdateThenPoll(ctx, id, spacecraft); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
	}
}

func (r SpacecraftResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.SpacecraftClient
			id, err := spacecraft.ParseSpacecraftID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				twoLineElements := []string{props.TleLine1, props.TleLine2}
				state := SpacecraftResourceModel{
					Name:            id.SpacecraftName,
					ResourceGroup:   id.ResourceGroupName,
					Location:        model.Location,
					NoradId:         *props.NoradId,
					TwoLineElements: twoLineElements,
					TitleLine:       props.TitleLine,
				}
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				spacecraftLinks, err := flattenSpacecraftLinks(props.Links)
				if err != nil {
					return err
				}
				state.Links = spacecraftLinks

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r SpacecraftResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.SpacecraftClient
			id, err := spacecraft.ParseSpacecraftID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SpacecraftResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return spacecraft.ValidateSpacecraftID
}

func (r SpacecraftResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.SpacecraftClient
			id, err := spacecraft.ParseSpacecraftID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state SpacecraftResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			spacecraftLinks, err := expandSpacecraftLinks(state.Links)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChangesExcept("name", "resource_group_name") {
				spacecraft := spacecraft.Spacecraft{
					Location: state.Location,
					Properties: spacecraft.SpacecraftsProperties{
						Links:     spacecraftLinks,
						NoradId:   utils.String(state.NoradId),
						TitleLine: state.TitleLine,
						TleLine1:  state.TwoLineElements[0],
						TleLine2:  state.TwoLineElements[1],
					},
					Tags: &state.Tags,
				}

				if err := client.CreateOrUpdateThenPoll(ctx, *id, spacecraft); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}
