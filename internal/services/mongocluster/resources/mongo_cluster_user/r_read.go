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

func (r Resource) Read() sdk.ResourceFunc {
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

func flattenDatabaseRoles(input *[]users.DatabaseRole) []DatabaseRoleModel {
	results := make([]DatabaseRoleModel, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		results = append(results, DatabaseRoleModel{
			Database: v.Db,
			Name:     string(v.Role),
		})
	}

	return results
}
