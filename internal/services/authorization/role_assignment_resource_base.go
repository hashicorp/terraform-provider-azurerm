// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type roleAssignmentBaseResource struct{}

func (br roleAssignmentBaseResource) arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"role_definition_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"role_definition_name"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"role_definition_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"role_definition_id"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"skip_service_principal_aad_check": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"delegated_managed_identity_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"condition_version"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"condition"},
			ValidateFunc: validation.StringInSlice([]string{
				"1.0",
				"2.0",
			}, false),
		},
	}
}

func (br roleAssignmentBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (br roleAssignmentBaseResource) createFunc(resourceName, scope string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			roleAssignmentsClient := metadata.Client.Authorization.ScopedRoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.ScopedRoleDefinitionsClient
			subscriptionClient := metadata.Client.Subscription.SubscriptionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			name := metadata.ResourceData.Get("name").(string)

			var roleDefinitionId string
			if v, ok := metadata.ResourceData.GetOk("role_definition_id"); ok {
				roleDefinitionId = v.(string)
			}

			if v, ok := metadata.ResourceData.GetOk("role_definition_name"); ok {
				roleName := v.(string)
				roleDefinitions, err := roleDefinitionsClient.List(ctx, commonids.NewScopeID(scope), roledefinitions.ListOperationOptions{Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", roleName))})
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

			principalId := metadata.ResourceData.Get("principal_id").(string)

			var err error
			if name == "" {
				name, err = uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
				}
			}

			tenantId := ""
			delegatedManagedIdentityResourceID := metadata.ResourceData.Get("delegated_managed_identity_resource_id").(string)
			if len(delegatedManagedIdentityResourceID) > 0 {
				tenantId, err = getTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
				if err != nil {
					return err
				}
			}

			id := parse.NewScopedRoleAssignmentID(scope, name, tenantId)
			options := roleassignments.DefaultGetOperationOptions()
			if tenantId != "" {
				options.TenantId = &tenantId
			}

			existing, err := roleAssignmentsClient.Get(ctx, id.ScopedId, options)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			properties := roleassignments.RoleAssignmentCreateParameters{
				Properties: roleassignments.RoleAssignmentProperties{
					RoleDefinitionId: roleDefinitionId,
					PrincipalId:      principalId,
					Description:      pointer.To(metadata.ResourceData.Get("description").(string)),
				},
			}

			if len(delegatedManagedIdentityResourceID) > 0 {
				properties.Properties.DelegatedManagedIdentityResourceId = &delegatedManagedIdentityResourceID
			}

			condition := metadata.ResourceData.Get("condition").(string)
			conditionVersion := metadata.ResourceData.Get("condition_version").(string)

			if condition != "" {
				properties.Properties.Condition = &condition
				properties.Properties.ConditionVersion = &conditionVersion
			}

			skipPrincipalCheck := metadata.ResourceData.Get("skip_service_principal_aad_check").(bool)
			if skipPrincipalCheck {
				properties.Properties.PrincipalType = pointer.To(roleassignments.PrincipalTypeServicePrincipal)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id)
			}

			if err = pluginsdk.Retry(time.Until(deadline), br.retryRoleAssignmentsClient(ctx, metadata, id, &properties)); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},

		Timeout: 30 * time.Minute,
	}
}

func (br roleAssignmentBaseResource) readFunc(scope string, isTenantLevel bool) sdk.ResourceFunc {
	return sdk.ResourceFunc{
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

			if model := resp.Model; model != nil {
				metadata.ResourceData.Set("name", model.Name)
				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("role_definition_id", props.RoleDefinitionId)
					metadata.ResourceData.Set("principal_id", props.PrincipalId)
					metadata.ResourceData.Set("delegated_managed_identity_resource_id", props.DelegatedManagedIdentityResourceId)
					metadata.ResourceData.Set("description", props.Description)
					metadata.ResourceData.Set("condition", props.Condition)
					metadata.ResourceData.Set("condition_version", props.ConditionVersion)

					if props.PrincipalType != nil {
						metadata.ResourceData.Set("principal_type", pointer.From(props.PrincipalType))
					}

					// allows for import when role name is used (also if the role name changes a plan will show a diff)
					roleId := props.RoleDefinitionId
					// The tenant level role definitions do not have a scope
					if isTenantLevel {
						roleId = fmt.Sprintf("%s%s", scope, props.RoleDefinitionId)
					}

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
							metadata.ResourceData.Set("role_definition_name", roleProps.RoleName)
						}
					}
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (br roleAssignmentBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleAssignmentsClient

			id, err := parse.ScopedRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			options := roleassignments.DefaultDeleteOperationOptions()
			if id.TenantId != "" {
				options.TenantId = &id.TenantId
			}

			resp, err := client.Delete(ctx, id.ScopedId, options)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return err
				}
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}

func (br roleAssignmentBaseResource) retryRoleAssignmentsClient(ctx context.Context, metadata sdk.ResourceMetaData, id parse.ScopedRoleAssignmentId, properties *roleassignments.RoleAssignmentCreateParameters) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		roleAssignmentsClient := metadata.Client.Authorization.ScopedRoleAssignmentsClient
		resp, err := roleAssignmentsClient.Create(ctx, id.ScopedId, *properties)
		if err != nil {
			if utils.ResponseErrorIsRetryable(err) {
				return pluginsdk.RetryableError(err)
			} else if response.WasStatusCode(resp.HttpResponse, 400) && strings.Contains(err.Error(), "PrincipalNotFound") {
				// When waiting for service principal to become available
				return pluginsdk.RetryableError(err)
			}

			return pluginsdk.NonRetryableError(err)
		}

		if resp.Model == nil || resp.Model.Id == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("creation of Role Assignment %s did not return an id value", id))
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			return pluginsdk.NonRetryableError(fmt.Errorf("could not retrieve context deadline for %s", metadata.ResourceData.Id()))
		}

		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{
				"pending",
			},
			Target: []string{
				"ready",
			},
			Refresh:                   br.roleAssignmentCreateStateRefreshFunc(ctx, roleAssignmentsClient, id),
			MinTimeout:                5 * time.Second,
			ContinuousTargetOccurence: 5,
			Timeout:                   time.Until(deadline),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("failed waiting for Role Assignment %s to finish replicating: %+v", id, err))
		}

		return nil
	}
}

func (br roleAssignmentBaseResource) roleAssignmentCreateStateRefreshFunc(ctx context.Context, client *roleassignments.RoleAssignmentsClient, id parse.ScopedRoleAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		options := roleassignments.DefaultGetByIdOperationOptions()
		if id.TenantId != "" {
			options.TenantId = &id.TenantId
		}

		resp, err := client.GetById(ctx, commonids.NewScopeID(id.ScopedId.ID()), options)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "pending", nil
			}
			return resp, "failed", err
		}
		return resp, "ready", nil
	}
}
