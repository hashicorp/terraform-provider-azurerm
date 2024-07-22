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

type IotCentralServicePrincipalUserResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotCentralServicePrincipalUserResource{}
)

type IotCentralServicePrincipalUserModel struct {
	IotCentralApplicationId string `tfschema:"iotcentral_application_id"`
	UserId                  string `tfschema:"user_id"`
	Type                    string `tfschema:"type"`
	TenantId                string `tfschema:"tenant_id"`
	ObjectId                string `tfschema:"object_id"`
	Role                    []Role `tfschema:"role"`
}

func (r IotCentralServicePrincipalUserResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r IotCentralServicePrincipalUserResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r IotCentralServicePrincipalUserResource) ResourceType() string {
	return "azurerm_iotcentral_service_principal_user"
}

func (r IotCentralServicePrincipalUserResource) ModelObject() interface{} {
	return &IotCentralServicePrincipalUserModel{}
}

func (IotCentralServicePrincipalUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.UserID
}

func (u IotCentralServicePrincipalUserModel) AsServicePrincipalUser() dataplane.ServicePrincipalUser {
	return dataplane.ServicePrincipalUser{
		TenantID: &u.TenantId,
		ObjectID: &u.ObjectId,
		ID:       &u.UserId,
		Type:     dataplane.TypeBasicUserTypeServicePrincipal,
		Roles:    ConvertToRoleAssignments(u.Role),
	}
}

func TryValidateServicePrincipalUserExistence(user dataplane.BasicUser) (string, bool) {
	if userValue, ok := user.AsServicePrincipalUser(); ok {
		return *userValue.ID, true
	}
	return "", false
}

func (r IotCentralServicePrincipalUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralServicePrincipalUserModel
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

			userToCreate := state.AsServicePrincipalUser()

			user, err := userClient.Create(ctx, state.UserId, userToCreate)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.UserId, err)
			}

			_, isValid := TryValidateServicePrincipalUserExistence(user.Value)
			if !isValid {
				return fmt.Errorf("unable to validate existence of user: id = %+v, type = ServicePrincipal after creating user: %+v", state.UserId, userToCreate)
			}

			id := parse.NewUserID(appId.SubscriptionId, appId.ResourceGroupName, appId.IotAppName, state.UserId)

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralServicePrincipalUserResource) Read() sdk.ResourceFunc {
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
			_, isValid := TryValidateServicePrincipalUserExistence(user.Value)
			if err != nil {
				if !isValid {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			servicePrincipalUser, isValid := user.Value.AsServicePrincipalUser()
			if !isValid {
				return fmt.Errorf("unable to convert user to type ServicePrincipalUser")
			}

			var state = IotCentralServicePrincipalUserModel{
				IotCentralApplicationId: appId.ID(),
				UserId:                  id.Name,
				Type:                    "ServicePrincipal",
				TenantId:                *servicePrincipalUser.TenantID,
				ObjectId:                *servicePrincipalUser.ObjectID,
				Role:                    ConvertFromRoleAssignments(servicePrincipalUser.Roles),
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r IotCentralServicePrincipalUserResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralServicePrincipalUserModel
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
			_, isValid := TryValidateServicePrincipalUserExistence(existing.Value)
			if err != nil {
				if !isValid {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			servicePrincipalUser, _ := existing.Value.AsServicePrincipalUser()

			if metadata.ResourceData.HasChange("role") {
				servicePrincipalUser.Roles = ConvertToRoleAssignments(state.Role)
			}

			_, err = userClient.Update(ctx, *servicePrincipalUser.ID, servicePrincipalUser, "*")
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralServicePrincipalUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralServicePrincipalUserModel
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
