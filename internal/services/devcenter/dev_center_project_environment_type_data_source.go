// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DevCenterProjectEnvironmentTypeDataSource{}

type DevCenterProjectEnvironmentTypeDataSource struct{}

type DevCenterProjectEnvironmentTypeDataSourceModel struct {
	Name                       string                                                        `tfschema:"name"`
	Location                   string                                                        `tfschema:"location"`
	DevCenterProjectId         string                                                        `tfschema:"dev_center_project_id"`
	DeploymentTargetId         string                                                        `tfschema:"deployment_target_id"`
	CreatorRoleAssignmentRoles []string                                                      `tfschema:"creator_role_assignment_roles"`
	Identity                   []identity.ModelSystemAssignedUserAssigned                    `tfschema:"identity"`
	UserRoleAssignment         []DevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment `tfschema:"user_role_assignment"`
	Tags                       map[string]string                                             `tfschema:"tags"`
}

type DevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment struct {
	UserId string   `tfschema:"user_id"`
	Roles  []string `tfschema:"roles"`
}

func (DevCenterProjectEnvironmentTypeDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterProjectEnvironmentTypeName,
		},

		"dev_center_project_id": commonschema.ResourceIDReferenceRequired(&environmenttypes.ProjectId{}),
	}
}

func (DevCenterProjectEnvironmentTypeDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"deployment_target_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"creator_role_assignment_roles": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"user_role_assignment": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"roles": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DevCenterProjectEnvironmentTypeDataSource) ModelObject() interface{} {
	return &DevCenterProjectEnvironmentTypeDataSourceModel{}
}

func (DevCenterProjectEnvironmentTypeDataSource) ResourceType() string {
	return "azurerm_dev_center_project_environment_type"
}

func (r DevCenterProjectEnvironmentTypeDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterProjectEnvironmentTypeDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterProjectId, err := projects.ParseProjectID(state.DevCenterProjectId)
			if err != nil {
				return err
			}

			id := environmenttypes.NewEnvironmentTypeID(subscriptionId, devCenterProjectId.ResourceGroupName, devCenterProjectId.ProjectName, state.Name)

			resp, err := client.ProjectEnvironmentTypesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Name = id.EnvironmentTypeName
				state.DevCenterProjectId = projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.ProjectName).ID()
				state.Location = location.NormalizeNilable(model.Location)
				state.Tags = pointer.From(model.Tags)

				identity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				state.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					state.DeploymentTargetId = pointer.From(props.DeploymentTargetId)
					state.UserRoleAssignment = flattenDevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment(props.UserRoleAssignments)

					if v := props.CreatorRoleAssignment; v != nil {
						state.CreatorRoleAssignmentRoles = flattenDevCenterProjectEnvironmentTypeDataSourceCreatorRoleAssignmentRoles(v.Roles)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment(input *map[string]environmenttypes.UserRoleAssignment) []DevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment {
	results := make([]DevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment, 0)

	if input == nil {
		return results
	}

	for k, v := range *input {
		result := DevCenterProjectEnvironmentTypeDataSourceUserRoleAssignment{
			UserId: k,
			Roles:  flattenDevCenterProjectEnvironmentTypeDataSourceUserRoleAssignmentRoles(v.Roles),
		}

		results = append(results, result)
	}

	return results
}

func flattenDevCenterProjectEnvironmentTypeDataSourceUserRoleAssignmentRoles(input *map[string]environmenttypes.EnvironmentRole) []string {
	result := make([]string, 0)

	if input == nil {
		return result
	}

	for k := range *input {
		result = append(result, k)
	}

	return result
}

func flattenDevCenterProjectEnvironmentTypeDataSourceCreatorRoleAssignmentRoles(input *map[string]environmenttypes.EnvironmentRole) []string {
	result := make([]string, 0)

	if input == nil {
		return result
	}

	for k := range *input {
		result = append(result, k)
	}

	return result
}
