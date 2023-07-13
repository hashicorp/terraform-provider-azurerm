// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package voiceservices

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/communicationsgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/testlines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CommunicationsGatewayTestLineResourceModel struct {
	Name                                 string                    `tfschema:"name"`
	Location                             string                    `tfschema:"location"`
	VoiceServicesCommunicationsGatewayId string                    `tfschema:"voice_services_communications_gateway_id"`
	PhoneNumber                          string                    `tfschema:"phone_number"`
	Purpose                              testlines.TestLinePurpose `tfschema:"purpose"`
	Tags                                 map[string]string         `tfschema:"tags"`
}

type CommunicationsGatewayTestLineResource struct{}

var _ sdk.ResourceWithUpdate = CommunicationsGatewayTestLineResource{}

func (r CommunicationsGatewayTestLineResource) ResourceType() string {
	return "azurerm_voice_services_communications_gateway_test_line"
}

func (r CommunicationsGatewayTestLineResource) ModelObject() interface{} {
	return &CommunicationsGatewayTestLineResourceModel{}
}

func (r CommunicationsGatewayTestLineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return testlines.ValidateTestLineID
}

func (r CommunicationsGatewayTestLineResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9-]{3,24}$"),
				"The name can only contain letters, numbers and dashes, the name length must be from 3 to 24 characters.",
			),
		},

		"location": commonschema.Location(),

		"voice_services_communications_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: communicationsgateways.ValidateCommunicationsGatewayID,
		},

		"phone_number": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"purpose": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(testlines.TestLinePurposeManual),
				string(testlines.TestLinePurposeAutomated),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r CommunicationsGatewayTestLineResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CommunicationsGatewayTestLineResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CommunicationsGatewayTestLineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.VoiceServices.TestLinesClient
			communicationsGatewayId, err := communicationsgateways.ParseCommunicationsGatewayID(model.VoiceServicesCommunicationsGatewayId)
			if err != nil {
				return err
			}

			id := testlines.NewTestLineID(communicationsGatewayId.SubscriptionId, communicationsGatewayId.ResourceGroupName, communicationsGatewayId.CommunicationsGatewayName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := testlines.TestLine{
				Location: location.Normalize(model.Location),
				Properties: &testlines.TestLineProperties{
					PhoneNumber: model.PhoneNumber,
					Purpose:     model.Purpose,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CommunicationsGatewayTestLineResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VoiceServices.TestLinesClient

			id, err := testlines.ParseTestLineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CommunicationsGatewayTestLineResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("phone_number") {
				properties.Properties.PhoneNumber = model.PhoneNumber
			}

			if metadata.ResourceData.HasChange("purpose") {
				properties.Properties.Purpose = model.Purpose
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CommunicationsGatewayTestLineResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VoiceServices.TestLinesClient

			id, err := testlines.ParseTestLineID(metadata.ResourceData.Id())
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state := CommunicationsGatewayTestLineResourceModel{
				Name:                                 id.TestLineName,
				VoiceServicesCommunicationsGatewayId: communicationsgateways.NewCommunicationsGatewayID(id.SubscriptionId, id.ResourceGroupName, id.CommunicationsGatewayName).ID(),
				Location:                             location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				state.PhoneNumber = properties.PhoneNumber
				state.Purpose = properties.Purpose
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CommunicationsGatewayTestLineResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VoiceServices.TestLinesClient

			id, err := testlines.ParseTestLineID(metadata.ResourceData.Id())
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
