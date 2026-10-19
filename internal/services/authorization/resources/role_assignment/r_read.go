// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_assignment

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmRoleAssignmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
	roleDefinitionsClient := meta.(*clients.Client).Authorization.ScopedRoleDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScopedRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	options := roleassignments.DefaultGetOperationOptions()
	if id.TenantId != "" {
		options.TenantId = pointer.To(id.TenantId)
	}

	resp, err := client.Get(ctx, id.ScopedId, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)

		if props := model.Properties; props != nil {
			d.Set("scope", normalizeScopeValue(pointer.From(props.Scope)))
			d.Set("role_definition_id", props.RoleDefinitionId)
			d.Set("principal_id", props.PrincipalId)
			d.Set("principal_type", pointer.From(props.PrincipalType))
			d.Set("delegated_managed_identity_resource_id", props.DelegatedManagedIdentityResourceId)
			d.Set("description", props.Description)
			d.Set("condition", props.Condition)
			d.Set("condition_version", props.ConditionVersion)

			if roleDefResourceId := props.RoleDefinitionId; roleDefResourceId != "" {
				// Workaround for https://github.com/hashicorp/pandora/issues/3257
				// The role definition id returned does not contain scope when the role definition was on tenant level (management group or tenant).
				// And adding tenant id as scope will cause 404 response, so just adding a slash to parse that.
				if strings.HasPrefix(roleDefResourceId, "/providers") {
					roleDefResourceId = fmt.Sprintf("/%s", roleDefResourceId)
				}
				parsedRoleDefId, err := roledefinitions.ParseScopedRoleDefinitionID(roleDefResourceId)
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", roleDefResourceId, err)
				}
				roleResp, err := roleDefinitionsClient.Get(ctx, *parsedRoleDefId)
				if err != nil {
					return fmt.Errorf("retrieving Role Definition %q: %+v", roleDefResourceId, err)
				}
				if roleResp.Model != nil && roleResp.Model.Properties != nil {
					d.Set("role_definition_name", pointer.From(roleResp.Model.Properties.RoleName))
				}
			}
		}
	}

	return nil
}

func retryRoleAssignmentsClient(d *pluginsdk.ResourceData, id parse.ScopedRoleAssignmentId, param roleassignments.RoleAssignmentCreateParameters, meta any, retryLinkedAuthorizationFailedError bool) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		roleAssignmentsClient := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
		ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := roleAssignmentsClient.Create(ctx, id.ScopedId, param)
		if err != nil {
			switch {
			case response.WasStatusCode(resp.HttpResponse, 400) && strings.Contains(err.Error(), "PrincipalNotFound"):
				// When waiting for service principal to become available
				return pluginsdk.RetryableError(err)
			case retryLinkedAuthorizationFailedError && response.WasForbidden(resp.HttpResponse) && strings.Contains(err.Error(), "LinkedAuthorizationFailed"):
				return pluginsdk.RetryableError(err)
			default:
				return pluginsdk.NonRetryableError(err)
			}
		}

		if resp.Model == nil || resp.Model.Id == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("creation of %s did not return an id value", id))
		}

		pollerOpts := &custompollers.EventualConsistencyPollerOptions{
			Interval:              5 * time.Second,
			RetryErrorStatusCodes: []int{http.StatusNotFound},
		}
		options := roleassignments.DefaultGetOperationOptions()
		if id.TenantId != "" {
			options.TenantId = pointer.To(id.TenantId)
		}
		poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
			resp, err := roleAssignmentsClient.Get(pollerCtx, id.ScopedId, options)
			return resp.HttpResponse, err
		}, pollerOpts)
		if err := poller.PollUntilDone(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("failed waiting for Role Assignment %s to finish replicating: %+v", id, err))
		}

		return nil
	}
}

func normalizeScopeValue(scope string) (result string) {
	if rg, err := commonids.ParseResourceGroupIDInsensitively(scope); err == nil {
		return rg.ID()
	}

	if ag, err := applicationgroup.ParseApplicationGroupIDInsensitively(scope); err == nil {
		return ag.ID()
	}

	// only check part of IDs, there are may be other specific resource types, like storage account id
	// we may need append these parse logics below when needed
	return scope
}
