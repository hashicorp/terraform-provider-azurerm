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
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
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
			roleAssignmentsClient := metadata.Client.Authorization.NewRoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.NewRoleDefinitionsClient
			subscriptionClient := metadata.Client.Subscription.Client
			subscriptionId := metadata.Client.Account.SubscriptionId
			name := metadata.ResourceData.Get("name").(string)

			var roleDefinitionId string
			if v, ok := metadata.ResourceData.GetOk("role_definition_id"); ok {
				roleDefinitionId = v.(string)
			} else if v, ok := metadata.ResourceData.GetOk("role_definition_name"); ok {
				roleName := v.(string)
				roleDefinitions, err := roleDefinitionsClient.List(ctx, commonids.NewScopeID(scope), roledefinitions.ListOperationOptions{Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", roleName))})
				if err != nil {
					return fmt.Errorf("loading Role Definition List: %+v", err)
				}

				if roleDefinitions.Model != nil && len(*roleDefinitions.Model) == 1 && (*roleDefinitions.Model)[0].Id != nil {
					roleDefinitionId = *(*roleDefinitions.Model)[0].Id
				} else {
					return fmt.Errorf("loading Role Definition List: failed to find role '%s'", roleName)
				}
			} else {
				return fmt.Errorf("either 'role_definition_id' or 'role_definition_name' needs to be set")
			}

			metadata.ResourceData.Set("role_definition_id", roleDefinitionId)

			principalId := metadata.ResourceData.Get("principal_id").(string)

			if name == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
				}

				name = uuid
			}

			tenantId := ""
			delegatedManagedIdentityResourceID := metadata.ResourceData.Get("delegated_managed_identity_resource_id").(string)
			if len(delegatedManagedIdentityResourceID) > 0 {
				var err error
				tenantId, err = getTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
				if err != nil {
					return err
				}
			}

			id := roleassignments.NewScopedRoleAssignmentID(scope, name)
			options := roleassignments.DefaultGetOperationOptions()
			if tenantId != "" {
				options.TenantId = &tenantId
			}

			existing, err := roleAssignmentsClient.Get(ctx, id, options)
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

			if condition != "" && conditionVersion != "" {
				properties.Properties.Condition = &condition
				properties.Properties.ConditionVersion = &conditionVersion
			} else if condition != "" || conditionVersion != "" {
				return fmt.Errorf("`condition` and `conditionVersion` should be both set or unset")
			}

			skipPrincipalCheck := metadata.ResourceData.Get("skip_service_principal_aad_check").(bool)
			if skipPrincipalCheck {
				properties.Properties.PrincipalType = pointer.To(roleassignments.PrincipalTypeServicePrincipal)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id)
			}

			if err = pluginsdk.Retry(time.Until(deadline), br.retryRoleAssignmentsClient(ctx, metadata, id, &properties, tenantId)); err != nil {
				return err
			}

			read, err := roleAssignmentsClient.Get(ctx, id, options)
			if err != nil {
				return err
			}

			if read.Model != nil && read.Model.Id != nil {
				metadata.ResourceData.SetId(parse.ConstructRoleAssignmentId(*read.Model.Id, tenantId))
			} else {
				return fmt.Errorf("retrieving %s: %s", id, err)
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}

func (br roleAssignmentBaseResource) readFunc(scope string, isTenantLevel bool) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.NewRoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.NewRoleDefinitionsClient

			azureResourceId, tenantId := parse.DestructRoleAssignmentId(metadata.ResourceData.Id())
			id := commonids.NewScopeID(azureResourceId)
			options := roleassignments.DefaultGetByIdOperationOptions()
			if tenantId != "" {
				options.TenantId = &tenantId
			}

			resp, err := client.GetById(ctx, id, options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] Role Assignment ID %s was not found - removing from state", id)
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
					metadata.ResourceData.Set("principal_type", props.PrincipalType)
					metadata.ResourceData.Set("delegated_managed_identity_resource_id", props.DelegatedManagedIdentityResourceId)
					metadata.ResourceData.Set("description", props.Description)
					metadata.ResourceData.Set("condition", props.Condition)
					metadata.ResourceData.Set("condition_version", props.ConditionVersion)

					// allows for import when role name is used (also if the role name changes a plan will show a diff)
					roleId := props.RoleDefinitionId
					// The tenant level role definitions do not have a scope in URL
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
			client := metadata.Client.Authorization.NewRoleAssignmentsClient

			azureResourceId, tenantId := parse.DestructRoleAssignmentId(metadata.ResourceData.Id())
			id, err := roleassignments.ParseScopedRoleAssignmentID(azureResourceId)
			if err != nil {
				return err
			}

			options := roleassignments.DefaultDeleteOperationOptions()
			if tenantId != "" {
				options.TenantId = &tenantId
			}
			resp, err := client.Delete(ctx, *id, options)
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

func (br roleAssignmentBaseResource) retryRoleAssignmentsClient(ctx context.Context, metadata sdk.ResourceMetaData, id roleassignments.ScopedRoleAssignmentId, properties *roleassignments.RoleAssignmentCreateParameters, tenantId string) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		roleAssignmentsClient := metadata.Client.Authorization.NewRoleAssignmentsClient
		resp, err := roleAssignmentsClient.Create(ctx, id, *properties)
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
			Refresh:                   br.roleAssignmentCreateStateRefreshFunc(ctx, roleAssignmentsClient, id.ID(), tenantId),
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

func (br roleAssignmentBaseResource) roleAssignmentCreateStateRefreshFunc(ctx context.Context, client *roleassignments.RoleAssignmentsClient, id string, tenantId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		options := roleassignments.DefaultGetByIdOperationOptions()
		if tenantId != "" {
			options.TenantId = &tenantId
		}

		resp, err := client.GetById(ctx, commonids.NewScopeID(id), options)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "pending", nil
			}
			return resp, "failed", err
		}
		return resp, "ready", nil
	}
}

func (br roleAssignmentBaseResource) validateRoleAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	azureResourceId, _ := parse.DestructRoleAssignmentId(v)
	return roleassignments.ValidateScopedRoleAssignmentID(azureResourceId, key)
}
