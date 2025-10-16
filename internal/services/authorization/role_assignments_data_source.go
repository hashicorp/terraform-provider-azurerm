// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RoleAssignmentsDataSource struct{}

type RoleAssignmentsDataSourceModel struct {
	Scope           string                 `tfschema:"scope"`
	LimitAtScope    bool                   `tfschema:"limit_at_scope"`
	PrincipalID     string                 `tfschema:"principal_id"`
	TenantID        string                 `tfschema:"tenant_id"`
	RoleAssignments []RoleAssignmentsModel `tfschema:"role_assignments"`
}

type RoleAssignmentsModel struct {
	Condition                          string `tfschema:"condition"`
	ConditionVersion                   string `tfschema:"condition_version"`
	DelegatedManagedIdentityResourceID string `tfschema:"delegated_managed_identity_resource_id"`
	Description                        string `tfschema:"description"`
	PrincipalID                        string `tfschema:"principal_id"`
	PrincipalType                      string `tfschema:"principal_type"`
	RoleAssignmentID                   string `tfschema:"role_assignment_id"`
	RoleAssignmentName                 string `tfschema:"role_assignment_name"`
	RoleAssignmentScope                string `tfschema:"role_assignment_scope"`
	RoleDefinitionID                   string `tfschema:"role_definition_id"`
}

var _ sdk.DataSource = RoleAssignmentsDataSource{}

func (r RoleAssignmentsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"limit_at_scope": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"principal_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},
		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r RoleAssignmentsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_assignments": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"condition": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"condition_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"delegated_managed_identity_resource_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"principal_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"role_assignment_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"role_assignment_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"role_assignment_scope": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"role_definition_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r RoleAssignmentsDataSource) ModelObject() interface{} {
	return &RoleAssignmentsDataSourceModel{}
}

func (r RoleAssignmentsDataSource) ResourceType() string {
	return "azurerm_role_assignments"
}

func (r RoleAssignmentsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleAssignmentsClient

			var state RoleAssignmentsDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewScopeID(state.Scope)

			options := roleassignments.DefaultListForScopeOperationOptions()
			// Root scope requires a filter, default to `atScope()`
			if state.Scope == "/" {
				options.Filter = pointer.To("atScope()")
			}

			if state.TenantID != "" {
				options.TenantId = pointer.To(state.TenantID)
			}

			if state.PrincipalID != "" {
				options.Filter = pointer.To(fmt.Sprintf("principalId eq '%s'", state.PrincipalID))
			}

			resp, err := client.ListForScope(ctx, id, options)
			if err != nil {
				return fmt.Errorf("listing role assignments for %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.RoleAssignments = flattenRoleAssignmentsToModel(model, state.Scope, state.LimitAtScope)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenRoleAssignmentsToModel(input *[]roleassignments.RoleAssignment, scope string, limitAtScope bool) []RoleAssignmentsModel {
	result := make([]RoleAssignmentsModel, 0)

	if len(*input) == 0 {
		return result
	}

	for _, v := range *input {
		assignment := RoleAssignmentsModel{
			RoleAssignmentID:   pointer.From(v.Id),
			RoleAssignmentName: pointer.From(v.Name),
		}

		if props := v.Properties; props != nil {
			// The API returns all role assignments at, above, or below the provided scope regardless of whether `atScope()` is passed as a filter.
			// If user set `limit_at_scope` to `true` and configuration scope != returned scope, discard
			if limitAtScope && !strings.EqualFold(scope, pointer.From(props.Scope)) {
				continue
			}

			assignment.Condition = pointer.From(props.Condition)
			assignment.ConditionVersion = pointer.From(props.ConditionVersion)
			assignment.DelegatedManagedIdentityResourceID = pointer.From(props.DelegatedManagedIdentityResourceId)
			assignment.Description = pointer.From(props.Description)
			assignment.PrincipalID = props.PrincipalId
			assignment.PrincipalType = string(pointer.From(props.PrincipalType))
			assignment.RoleAssignmentScope = pointer.From(props.Scope)
			assignment.RoleDefinitionID = props.RoleDefinitionId
		}

		result = append(result, assignment)
	}

	return result
}
