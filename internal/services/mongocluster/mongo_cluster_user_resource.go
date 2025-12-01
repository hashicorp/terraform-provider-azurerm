// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongocluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/mongoclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/users"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MongoClusterUserResource struct{}

var _ sdk.Resource = MongoClusterUserResource{}

type MongoClusterUserResourceModel struct {
	IdentityProviderType string              `tfschema:"identity_provider_type"`
	MongoClusterId       string              `tfschema:"mongo_cluster_id"`
	ObjectId             string              `tfschema:"object_id"`
	PrincipalType        string              `tfschema:"principal_type"`
	Roles                []DatabaseRoleModel `tfschema:"role"`
}

type DatabaseRoleModel struct {
	Database string `tfschema:"database"`
	Role     string `tfschema:"role"`
}

func (r MongoClusterUserResource) ModelObject() interface{} {
	return &MongoClusterUserResourceModel{}
}

func (r MongoClusterUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return users.ValidateUserID
}

func (r MongoClusterUserResource) ResourceType() string {
	return "azurerm_mongo_cluster_user"
}

func (r MongoClusterUserResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// `object_id` is actually resource name
		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"mongo_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&users.MongoClusterId{}),

		"identity_provider_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				users.PossibleValuesForIdentityProviderType(),
				false,
			),
		},

		"principal_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				users.PossibleValuesForEntraPrincipalType(),
				false,
			),
		},

		"role": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"database": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"role": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice(
							users.PossibleValuesForUserRole(),
							false,
						),
					},
				},
			},
		},
	}
}

func (r MongoClusterUserResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MongoClusterUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.UsersClient

			var state MongoClusterUserResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			mongoClusterId, err := mongoclusters.ParseMongoClusterID(state.MongoClusterId)
			if err != nil {
				return err
			}

			id := users.NewUserID(mongoClusterId.SubscriptionId, mongoClusterId.ResourceGroupName, mongoClusterId.MongoClusterName, state.ObjectId)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityProvider := users.EntraIdentityProvider{
				Type: users.IdentityProviderType(state.IdentityProviderType),
				Properties: users.EntraIdentityProviderProperties{
					PrincipalType: users.EntraPrincipalType(state.PrincipalType),
				},
			}

			parameter := users.User{
				Properties: &users.UserProperties{
					IdentityProvider: identityProvider,
					Roles:            expandDatabaseRoles(state.Roles),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameter); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MongoClusterUserResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.UsersClient

			id, err := users.ParseUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MongoClusterUserResourceModel{
				ObjectId:       id.UserName,
				MongoClusterId: mongoclusters.NewMongoClusterID(id.SubscriptionId, id.ResourceGroupName, id.MongoClusterName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Roles = flattenDatabaseRoles(props.Roles)

					if identityProvider, ok := props.IdentityProvider.(users.EntraIdentityProvider); ok {
						state.PrincipalType = string(identityProvider.Properties.PrincipalType)
						state.IdentityProviderType = string(identityProvider.Type)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MongoClusterUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.UsersClient

			id, err := users.ParseUserID(metadata.ResourceData.Id())
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

func expandDatabaseRoles(input []DatabaseRoleModel) *[]users.DatabaseRole {
	if len(input) == 0 {
		return nil
	}

	result := make([]users.DatabaseRole, 0)
	for _, v := range input {
		result = append(result, users.DatabaseRole{
			Db:   v.Database,
			Role: users.UserRole(v.Role),
		})
	}

	return &result
}

func flattenDatabaseRoles(input *[]users.DatabaseRole) []DatabaseRoleModel {
	results := make([]DatabaseRoleModel, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		results = append(results, DatabaseRoleModel{
			Database: v.Db,
			Role:     string(v.Role),
		})
	}

	return results
}
