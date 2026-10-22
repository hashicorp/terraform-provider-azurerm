// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster_user

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/users"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameter, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

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
			Role: users.UserRole(v.Name),
		})
	}

	return &result
}
