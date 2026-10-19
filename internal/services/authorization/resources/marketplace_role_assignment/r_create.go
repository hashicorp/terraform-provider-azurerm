// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_role_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			roleAssignmentsClient := metadata.Client.Authorization.ScopedRoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.ScopedRoleDefinitionsClient
			subscriptionClient := metadata.Client.Subscription.SubscriptionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config RoleAssignmentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			name := config.Name
			roleDefinitionId := config.RoleDefinitionId

			if config.RoleDefinitionName != "" {
				roleName := config.RoleDefinitionName
				roleDefinitions, err := roleDefinitionsClient.List(ctx, commonids.NewScopeID(MarketplaceScope), roledefinitions.ListOperationOptions{Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", roleName))})
				if err != nil {
					return fmt.Errorf("loading Role Definition List: %+v", err)
				}

				if roleDefinitions.Model == nil || len(*roleDefinitions.Model) != 1 || (*roleDefinitions.Model)[0].Id == nil {
					return fmt.Errorf("loading Role Definition List: failed to find role '%s'", roleName)
				}

				roleDefinitionId = *(*roleDefinitions.Model)[0].Id
			}

			if roleDefinitionId == "" {
				return fmt.Errorf("either 'role_definition_id' or 'role_definition_name' needs to be set")
			}

			metadata.ResourceData.Set("role_definition_id", roleDefinitionId)

			var err error
			if name == "" {
				name, err = uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
				}
			}

			tenantId := ""
			if len(config.DelegatedManagedIdentityResourceId) > 0 {
				tenantId, err = helpers.GetTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
				if err != nil {
					return err
				}
			}

			id := parse.NewScopedRoleAssignmentID(MarketplaceScope, name, tenantId)
			options := roleassignments.DefaultGetOperationOptions()
			if tenantId != "" {
				options.TenantId = &tenantId
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := roleAssignmentsClient.Get(ctx, id.ScopedId, options)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %s", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			properties := roleassignments.RoleAssignmentCreateParameters{
				Properties: roleassignments.RoleAssignmentProperties{
					RoleDefinitionId: roleDefinitionId,
					PrincipalId:      config.PrincipalId,
					Description:      pointer.To(config.Description),
				},
			}

			if len(config.DelegatedManagedIdentityResourceId) > 0 {
				properties.Properties.DelegatedManagedIdentityResourceId = &config.DelegatedManagedIdentityResourceId
			}

			if config.Condition != "" {
				properties.Properties.Condition = &config.Condition
				properties.Properties.ConditionVersion = &config.ConditionVersion
			}

			if config.SkipServicePrincipalAadCheck {
				properties.Properties.PrincipalType = pointer.To(roleassignments.PrincipalTypeServicePrincipal)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id)
			}

			if err = pluginsdk.Retry(time.Until(deadline), retryMarketplaceRoleAssignmentsClient(ctx, metadata, id, &properties)); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}
