// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_role_assignment

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			id, err := parse.ScopedRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			options := roleassignments.DefaultGetByIdOperationOptions()
			if id.TenantId != "" {
				options.TenantId = &id.TenantId
			}

			resp, err := client.GetById(ctx, commonids.NewScopeID(id.ScopedId.ID()), options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state", id)
					metadata.ResourceData.SetId("")
					return nil
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := RoleAssignmentModel{}

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.Name)
				if props := model.Properties; props != nil {
					state.RoleDefinitionId = props.RoleDefinitionId
					state.PrincipalId = props.PrincipalId
					state.DelegatedManagedIdentityResourceId = pointer.From(props.DelegatedManagedIdentityResourceId)
					state.Description = pointer.From(props.Description)
					state.Condition = pointer.From(props.Condition)
					state.ConditionVersion = pointer.From(props.ConditionVersion)

					if props.PrincipalType != nil {
						state.PrincipalType = pointer.FromEnum(props.PrincipalType)
					}

					// allows for import when role name is used (also if the role name changes a plan will show a diff)
					// The tenant level role definitions do not have a scope
					roleId := fmt.Sprintf("%s%s", MarketplaceScope, props.RoleDefinitionId)

					roleDefinitionId, err := roledefinitions.ParseScopedRoleDefinitionID(roleId)
					if err != nil {
						return err
					}

					roleResp, err := roleDefinitionsClient.Get(ctx, *roleDefinitionId)
					if err != nil {
						return fmt.Errorf("retrieving %s: %s", roleDefinitionId, err)
					}

					if roleModel := roleResp.Model; roleModel != nil {
						if roleProps := roleModel.Properties; roleProps != nil {
							state.RoleDefinitionName = pointer.From(roleProps.RoleName)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

const (
	MarketplaceScope = "/providers/Microsoft.Marketplace"
)

func retryMarketplaceRoleAssignmentsClient(ctx context.Context, metadata sdk.ResourceMetaData, id parse.ScopedRoleAssignmentId, properties *roleassignments.RoleAssignmentCreateParameters) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		roleAssignmentsClient := metadata.Client.Authorization.ScopedRoleAssignmentsClient
		resp, err := roleAssignmentsClient.Create(ctx, id.ScopedId, *properties)
		if err != nil {
			if response.WasStatusCode(resp.HttpResponse, 400) && strings.Contains(err.Error(), "PrincipalNotFound") {
				// When waiting for service principal to become available
				return pluginsdk.RetryableError(err)
			}

			return pluginsdk.NonRetryableError(err)
		}

		if resp.Model == nil || resp.Model.Id == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("creation of Role Assignment %s did not return an id value", id))
		}

		pollerOpts := &custompollers.EventualConsistencyPollerOptions{
			Interval:              5 * time.Second,
			RetryErrorStatusCodes: []int{http.StatusNotFound},
		}
		options := roleassignments.DefaultGetByIdOperationOptions()
		if id.TenantId != "" {
			options.TenantId = &id.TenantId
		}

		poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
			resp, err := roleAssignmentsClient.GetById(pollerCtx, commonids.NewScopeID(id.ScopedId.ID()), options)
			return resp.HttpResponse, err
		}, pollerOpts)
		if err := poller.PollUntilDone(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("failed waiting for Role Assignment %s to finish replicating: %+v", id, err))
		}

		return nil
	}
}
