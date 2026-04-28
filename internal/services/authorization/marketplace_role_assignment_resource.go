// Copyright IBM Corp. 2014, 2025
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

const (
	MarketplaceScope = "/providers/Microsoft.Marketplace"
)

var _ sdk.Resource = RoleAssignmentMarketplaceResource{}

type RoleAssignmentModel struct {
	Name                               string `tfschema:"name"`
	PrincipalId                        string `tfschema:"principal_id"`
	RoleDefinitionId                   string `tfschema:"role_definition_id"`
	RoleDefinitionName                 string `tfschema:"role_definition_name"`
	SkipServicePrincipalAadCheck       bool   `tfschema:"skip_service_principal_aad_check"`
	DelegatedManagedIdentityResourceId string `tfschema:"delegated_managed_identity_resource_id"`
	Description                        string `tfschema:"description"`
	Condition                          string `tfschema:"condition"`
	ConditionVersion                   string `tfschema:"condition_version"`
	PrincipalType                      string `tfschema:"principal_type"`
}

type RoleAssignmentMarketplaceResource struct{}

func (r RoleAssignmentMarketplaceResource) ModelObject() interface{} {
	return &RoleAssignmentModel{}
}

func (r RoleAssignmentMarketplaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateScopedRoleAssignmentID
}

func (r RoleAssignmentMarketplaceResource) ResourceType() string {
	return "azurerm_marketplace_role_assignment"
}

func (r RoleAssignmentMarketplaceResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r RoleAssignmentMarketplaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RoleAssignmentMarketplaceResource) Create() sdk.ResourceFunc {
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
				tenantId, err = getTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
				if err != nil {
					return err
				}
			}

			id := parse.NewScopedRoleAssignmentID(MarketplaceScope, name, tenantId)
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
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
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

func (r RoleAssignmentMarketplaceResource) Read() sdk.ResourceFunc {
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
						state.PrincipalType = string(pointer.From(props.PrincipalType))
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

func (r RoleAssignmentMarketplaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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
	}
}

func retryMarketplaceRoleAssignmentsClient(ctx context.Context, metadata sdk.ResourceMetaData, id parse.ScopedRoleAssignmentId, properties *roleassignments.RoleAssignmentCreateParameters) func() *pluginsdk.RetryError {
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
			Refresh:                   marketplaceRoleAssignmentCreateStateRefreshFunc(ctx, roleAssignmentsClient, id),
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

func marketplaceRoleAssignmentCreateStateRefreshFunc(ctx context.Context, client *roleassignments.RoleAssignmentsClient, id parse.ScopedRoleAssignmentId) pluginsdk.StateRefreshFunc {
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
