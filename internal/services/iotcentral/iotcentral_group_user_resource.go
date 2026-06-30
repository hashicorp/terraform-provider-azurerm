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

type IotCentralGroupUserResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotCentralGroupUserResource{}
)

type IotCentralGroupUserModel struct {
	IotCentralApplicationId string `tfschema:"iotcentral_application_id"`
	UserId                  string `tfschema:"user_id"`
	Type                    string `tfschema:"type"`
	TenantId                string `tfschema:"tenant_id"`
	ObjectId                string `tfschema:"object_id"`
	Role                    []Role `tfschema:"role"`
}

func (r IotCentralGroupUserResource) Arguments() map[string]*pluginsdk.Schema {
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
		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"object_id": {
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

func (r IotCentralGroupUserResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r IotCentralGroupUserResource) ResourceType() string {
	return "azurerm_iotcentral_group_user"
}

func (r IotCentralGroupUserResource) ModelObject() interface{} {
	return &IotCentralGroupUserModel{}
}

func (IotCentralGroupUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.UserID
}

func (u IotCentralGroupUserModel) AsGroupUser() dataplane.ADGroupUser {
	return dataplane.ADGroupUser{
		TenantID: &u.TenantId,
		ObjectID: &u.ObjectId,
		ID:       &u.UserId,
		Type:     dataplane.TypeBasicUserTypeAdGroup,
		Roles:    ConvertToRoleAssignments(u.Role),
	}
}

func TryValidateGroupUserExistence(user dataplane.BasicUser) (string, bool) {
	if userValue, ok := user.AsADGroupUser(); ok {
		return *userValue.ID, true
	}
	return "", false
}

func (r IotCentralGroupUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralGroupUserModel
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

			userToCreate := state.AsGroupUser()

			user, err := userClient.Create(ctx, state.UserId, userToCreate)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.UserId, err)
			}

			_, isValid := TryValidateGroupUserExistence(user.Value)
			if !isValid {
				return fmt.Errorf("unable to validate existence of user: id = %+v, type = Group after creating user: %+v", state.UserId, userToCreate)
			}

			id := parse.NewUserID(appId.SubscriptionId, appId.ResourceGroupName, appId.IotAppName, state.UserId)

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralGroupUserResource) Read() sdk.ResourceFunc {
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
			_, isValid := TryValidateGroupUserExistence(user.Value)
			if err != nil {
				if !isValid {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			groupUser, isValid := user.Value.AsADGroupUser()
			if !isValid {
				return fmt.Errorf("unable to convert user to type GroupUser")
			}

			var state = IotCentralGroupUserModel{
				IotCentralApplicationId: appId.ID(),
				UserId:                  id.Name,
				Type:                    "Group",
				TenantId:                *groupUser.TenantID,
				ObjectId:                *groupUser.ObjectID,
				Role:                    ConvertFromRoleAssignments(groupUser.Roles),
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r IotCentralGroupUserResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralGroupUserModel
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
			_, isValid := TryValidateGroupUserExistence(existing.Value)
			if err != nil {
				if !isValid {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			groupUser, _ := existing.Value.AsADGroupUser()

			if metadata.ResourceData.HasChange("role") {
				groupUser.Roles = ConvertToRoleAssignments(state.Role)
			}

			_, err = userClient.Update(ctx, *groupUser.ID, groupUser, "*")
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralGroupUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralGroupUserModel
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
