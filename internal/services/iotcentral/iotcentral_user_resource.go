// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	dataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

type IotCentralUserResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotCentralUserResource{}
)

type IotCentralUserModel struct {
	IotCentralApplicationId string `tfschema:"iotcentral_application_id"`
	UserId                  string `tfschema:"user_id"`
	Type                    string `tfschema:"type"`
	Email                   string `tfschema:"email"`
	Role                    []Role `tfschema:"role"`
}

func (r IotCentralUserResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"iotcentral_application_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apps.ValidateIotAppID,
		},
		"user_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.UserUserID,
		},
		"email": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"role": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"role_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"organization_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.OrganizationOrganizationID,
					},
				},
			},
		},
	}
}

func (r IotCentralUserResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r IotCentralUserResource) ResourceType() string {
	return "azurerm_iotcentral_user"
}

func (r IotCentralUserResource) ModelObject() interface{} {
	return &IotCentralUserModel{}
}

func (IotCentralUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.UserID
}

func (u IotCentralUserModel) AsEmailUser() dataplane.EmailUser {
	return dataplane.EmailUser{
		Email: &u.Email,
		ID:    &u.UserId,
		Type:  dataplane.TypeBasicUserTypeEmail,
		Roles: ConvertToRoleAssignments(u.Role),
	}
}

func TryValidateUserExistence(user dataplane.BasicUser) (string, bool) {
	if userValue, ok := user.AsEmailUser(); ok {
		return *userValue.ID, true
	}
	return "", false
}

func (r IotCentralUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralUserModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			appId, err := apps.ParseIotAppID(state.IotCentralApplicationId)
			if err != nil {
				return err
			}

			app, err := client.AppsClient.Get(ctx, *appId)
			if err != nil || app.Model == nil {
				return fmt.Errorf("checking for the presence of existing %q: %+v", appId, err)
			}

			userClient, err := client.UsersClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating user client: %+v", err)
			}

			userToCreate := state.AsEmailUser()

			user, err := userClient.Create(ctx, state.UserId, userToCreate)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.UserId, err)
			}

			_, isValid := TryValidateUserExistence(user.Value)
			if !isValid {
				return fmt.Errorf("unable to validate existence of user: id = %+v, type = Email after creating user: %+v", state.UserId, userToCreate)
			}

			id := parse.NewUserID(appId.SubscriptionId, appId.ResourceGroupName, appId.IotAppName, state.UserId)

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralUserResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			id, err := parse.UserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := apps.NewIotAppID(id.SubscriptionId, id.ResourceGroup, id.IotAppName)
			if err != nil {
				return err
			}

			app, err := client.AppsClient.Get(ctx, appId)
			if err != nil || app.Model == nil {
				return metadata.MarkAsGone(id)
			}

			userClient, err := client.UsersClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating user client: %+v", err)
			}

			user, err := userClient.Get(ctx, id.Name)
			_, isValid := TryValidateUserExistence(user.Value)
			if err != nil {
				if !isValid {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			emailUser, isValid := user.Value.AsEmailUser()
			if !isValid {
				return fmt.Errorf("unable to convert user to type EmailUser")
			}

			var state = IotCentralUserModel{
				IotCentralApplicationId: appId.ID(),
				UserId:                  id.Name,
				Type:                    "Email",
				Email:                   *emailUser.Email,
				Role:                    ConvertFromRoleAssignments(emailUser.Roles),
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r IotCentralUserResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralUserModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.UserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := apps.NewIotAppID(id.SubscriptionId, id.ResourceGroup, id.IotAppName)
			if err != nil {
				return err
			}

			app, err := client.AppsClient.Get(ctx, appId)
			if err != nil || app.Model == nil {
				return metadata.MarkAsGone(id)
			}

			userClient, err := client.UsersClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating user client: %+v", err)
			}

			existing, err := userClient.Get(ctx, id.Name)
			_, isValid := TryValidateUserExistence(existing.Value)
			if err != nil {
				if !isValid {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			emailUser, _ := existing.Value.AsEmailUser()

			if metadata.ResourceData.HasChange("role") {
				emailUser.Roles = ConvertToRoleAssignments(state.Role)
			}

			_, err = userClient.Update(ctx, *emailUser.ID, emailUser, "*")
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralUserModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.UserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := apps.NewIotAppID(id.SubscriptionId, id.ResourceGroup, id.IotAppName)
			if err != nil {
				return err
			}

			app, err := client.AppsClient.Get(ctx, appId)
			if err != nil || app.Model == nil {
				return metadata.MarkAsGone(id)
			}

			orgClient, err := client.UsersClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating user client: %+v", err)
			}

			_, err = orgClient.Remove(ctx, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
