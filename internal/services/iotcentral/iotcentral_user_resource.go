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
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	TenantId                string `tfschema:"tenant_id"`
	ObjectId                string `tfschema:"object_id"`
	Role                    []Role `tfschema:"role"`
}

type Role struct {
	RoleId         string `tfschema:"role_id"`
	OrganizationId string `tfschema:"organization_id"`
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
		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.UserType,
		},
		"email": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"object_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"role": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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
	return map[string]*pluginsdk.Schema{}
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

func (u IotCentralUserModel) AsADGroupUser() (*dataplane.ADGroupUser, bool) {
	if u.Type == "Group" {
		return &dataplane.ADGroupUser{
			ObjectID: &u.ObjectId,
			TenantID: &u.TenantId,
			Type:     dataplane.TypeBasicUserTypeAdGroup,
			Roles:    convertToRoleAssignments(u.Role),
		}, true
	}
	return nil, false
}

func (u IotCentralUserModel) AsEmailUser() (*dataplane.EmailUser, bool) {
	if u.Type == "Email" {
		return &dataplane.EmailUser{
			Email: &u.Email,
			Type:  dataplane.TypeBasicUserTypeEmail,
			Roles: convertToRoleAssignments(u.Role),
		}, true
	}
	return nil, false
}

func (u IotCentralUserModel) AsServicePrincipalUser() (*dataplane.ServicePrincipalUser, bool) {
	if u.Type == "ServicePrincipal" {
		return &dataplane.ServicePrincipalUser{
			ObjectID: &u.ObjectId,
			TenantID: &u.TenantId,
			Type:     dataplane.TypeBasicUserTypeServicePrincipal,
			Roles:    convertToRoleAssignments(u.Role),
		}, true
	}
	return nil, false
}

func (u IotCentralUserModel) AsUser() (*dataplane.User, bool) {
	return nil, false
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

			user, err := userClient.Create(ctx, state.UserId, state)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.UserId, err)
			}

			userValue, valid := user.Value.AsUser()
			if !valid {
				return fmt.Errorf("unable to convert user to type User")
			}

			orgId := parse.NewUserID(appId.SubscriptionId, appId.ResourceGroupName, appId.IotAppName, *userValue.ID)

			metadata.SetID(orgId)
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
			if err != nil {
				userValue, isValid := user.Value.AsUser()
				if !isValid || userValue.ID == nil || *userValue.ID == "" {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			userValue, isValid := user.Value.AsUser()
			if !isValid {
				return fmt.Errorf("unable to convert user to type User")
			}

			var state IotCentralUserModel
			switch userValue.Type {
			case dataplane.TypeBasicUserTypeAdGroup:
				adGroupUser, isValid := user.Value.AsADGroupUser()
				if !isValid {
					return fmt.Errorf("unable to convert user to type ADGroupUser")
				}
				state = IotCentralUserModel{
					IotCentralApplicationId: id.IotAppName,
					UserId:                  id.Name,
					Type:                    "Group",
					TenantId:                *adGroupUser.TenantID,
					ObjectId:                *adGroupUser.ObjectID,
					Role:                    convertFromRoleAssignments(adGroupUser.Roles),
				}
			case dataplane.TypeBasicUserTypeServicePrincipal:
				servicePrincipalUser, isValid := user.Value.AsServicePrincipalUser()
				if !isValid {
					return fmt.Errorf("unable to convert user to type ServicePrincipalUser")
				}
				state = IotCentralUserModel{
					IotCentralApplicationId: id.IotAppName,
					UserId:                  id.Name,
					Type:                    "ServicePrincipal",
					TenantId:                *servicePrincipalUser.TenantID,
					ObjectId:                *servicePrincipalUser.ObjectID,
					Role:                    convertFromRoleAssignments(servicePrincipalUser.Roles),
				}
			case dataplane.TypeBasicUserTypeEmail:
				emailUser, isValid := user.Value.AsEmailUser()
				if !isValid {
					return fmt.Errorf("unable to convert user to type EmailUser")
				}
				state = IotCentralUserModel{
					IotCentralApplicationId: id.IotAppName,
					UserId:                  id.Name,
					Type:                    "Email",
					Email:                   *emailUser.Email,
					Role:                    convertFromRoleAssignments(emailUser.Roles),
				}
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
			if err != nil {
				userValue, isValid := existing.Value.AsUser()
				if !isValid || userValue.ID == nil || *userValue.ID == "" {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			userValue, isValid := existing.Value.AsUser()
			if !isValid {
				return fmt.Errorf("unable to convert user to type User")
			}

			if metadata.ResourceData.HasChange("role") {
				userValue.Roles = convertToRoleAssignments(state.Role)
			}

			_, err = userClient.Update(ctx, *userValue.ID, userValue, "*")
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

func convertToRoleAssignments(input []Role) *[]dataplane.RoleAssignment {
	if input == nil {
		return nil
	}

	results := make([]dataplane.RoleAssignment, 0)
	for _, item := range input {
		results = append(results, dataplane.RoleAssignment{
			Organization: utils.String(item.OrganizationId),
			Role:         utils.String(item.RoleId),
		})
	}
	return &results
}

func convertFromRoleAssignments(input *[]dataplane.RoleAssignment) []Role {
	if input == nil {
		return nil
	}

	results := make([]Role, 0)
	for _, item := range *input {
		obj := Role{}
		if item.Organization != nil {
			obj.OrganizationId = *item.Organization
		}
		if item.Role != nil {
			obj.RoleId = *item.Role
		}
		results = append(results, obj)
	}
	return results
}
