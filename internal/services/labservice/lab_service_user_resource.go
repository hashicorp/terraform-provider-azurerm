// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServiceUserModel struct {
	Name                 string `tfschema:"name"`
	LabId                string `tfschema:"lab_id"`
	Email                string `tfschema:"email"`
	AdditionalUsageQuota string `tfschema:"additional_usage_quota"`
}

type LabServiceUserResource struct{}

var _ sdk.ResourceWithUpdate = LabServiceUserResource{}

func (r LabServiceUserResource) ResourceType() string {
	return "azurerm_lab_service_user"
}

func (r LabServiceUserResource) ModelObject() interface{} {
	return &LabServiceUserModel{}
}

func (r LabServiceUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return user.ValidateUserID
}

func (r LabServiceUserResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"lab_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: lab.ValidateLabID,
		},

		"email": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Email,
		},

		"additional_usage_quota": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT0S",
			ValidateFunc: azValidate.ISO8601Duration,
		},
	}
}

func (r LabServiceUserResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServiceUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServiceUserModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LabService.UserClient
			labId, err := lab.ParseLabID(model.LabId)
			if err != nil {
				return err
			}

			id := user.NewUserID(labId.SubscriptionId, labId.ResourceGroupName, labId.LabName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &user.User{
				Properties: user.UserProperties{
					Email: model.Email,
				},
			}

			if model.AdditionalUsageQuota != "" {
				properties.Properties.AdditionalUsageQuota = &model.AdditionalUsageQuota
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LabServiceUserResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.UserClient

			id, err := user.ParseUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServiceUserModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			props := resp.Model
			if props == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("additional_usage_quota") {
				props.Properties.AdditionalUsageQuota = utils.String(model.AdditionalUsageQuota)
			}

			if metadata.ResourceData.HasChange("email") {
				props.Properties.Email = model.Email
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *props); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LabServiceUserResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.UserClient

			id, err := user.ParseUserID(metadata.ResourceData.Id())
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
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := LabServiceUserModel{
				Name:  id.UserName,
				LabId: lab.NewLabID(id.SubscriptionId, id.ResourceGroupName, id.LabName).ID(),
			}

			properties := &model.Properties
			state.Email = properties.Email

			if properties.AdditionalUsageQuota != nil {
				state.AdditionalUsageQuota = *properties.AdditionalUsageQuota
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServiceUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.UserClient

			id, err := user.ParseUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
