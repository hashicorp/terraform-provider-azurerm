// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = DevCenterProjectEnvironmentTypeResource{}
	_ sdk.ResourceWithUpdate = DevCenterProjectEnvironmentTypeResource{}
)

type DevCenterProjectEnvironmentTypeResource struct{}

func (r DevCenterProjectEnvironmentTypeResource) ModelObject() interface{} {
	return &DevCenterProjectEnvironmentTypeResourceModel{}
}

type DevCenterProjectEnvironmentTypeResourceModel struct {
	Name                       string                                              `tfschema:"name"`
	Location                   string                                              `tfschema:"location"`
	DevCenterProjectId         string                                              `tfschema:"dev_center_project_id"`
	DeploymentTargetId         string                                              `tfschema:"deployment_target_id"`
	CreatorRoleAssignmentRoles []string                                            `tfschema:"creator_role_assignment_roles"`
	Identity                   []identity.ModelSystemAssignedUserAssigned          `tfschema:"identity"`
	UserRoleAssignment         []DevCenterProjectEnvironmentTypeUserRoleAssignment `tfschema:"user_role_assignment"`
	Tags                       map[string]string                                   `tfschema:"tags"`
}

type DevCenterProjectEnvironmentTypeUserRoleAssignment struct {
	UserId string   `tfschema:"user_id"`
	Roles  []string `tfschema:"roles"`
}

func (r DevCenterProjectEnvironmentTypeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return environmenttypes.ValidateEnvironmentTypeID
}

func (r DevCenterProjectEnvironmentTypeResource) ResourceType() string {
	return "azurerm_dev_center_project_environment_type"
}

func (r DevCenterProjectEnvironmentTypeResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterProjectEnvironmentTypeName,
		},

		"location": commonschema.Location(),

		"dev_center_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&environmenttypes.ProjectId{}),

		"deployment_target_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"creator_role_assignment_roles": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"user_role_assignment": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"roles": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r DevCenterProjectEnvironmentTypeResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterProjectEnvironmentTypeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterProjectEnvironmentTypeResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterProjectId, err := projects.ParseProjectID(model.DevCenterProjectId)
			if err != nil {
				return err
			}

			id := environmenttypes.NewEnvironmentTypeID(subscriptionId, devCenterProjectId.ResourceGroupName, devCenterProjectId.ProjectName, model.Name)

			existing, err := client.ProjectEnvironmentTypesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return err
			}

			parameters := environmenttypes.ProjectEnvironmentType{
				Location: pointer.To(location.Normalize(model.Location)),
				Properties: &environmenttypes.ProjectEnvironmentTypeProperties{
					DeploymentTargetId: pointer.To(model.DeploymentTargetId),
					CreatorRoleAssignment: &environmenttypes.ProjectEnvironmentTypeUpdatePropertiesCreatorRoleAssignment{
						Roles: expandDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(model.CreatorRoleAssignmentRoles),
					},
					Status: pointer.To(environmenttypes.EnvironmentTypeEnableStatusEnabled),
				},
				Identity: identity,
				Tags:     pointer.To(model.Tags),
			}

			userRoleAssignment, err := expandDevCenterProjectEnvironmentTypeUserRoleAssignment(model.UserRoleAssignment)
			if err != nil {
				return err
			}
			parameters.Properties.UserRoleAssignments = userRoleAssignment

			if _, err := client.ProjectEnvironmentTypesCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterProjectEnvironmentTypeResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectEnvironmentTypesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterProjectEnvironmentTypeResourceModel{
				Name:               id.EnvironmentTypeName,
				DevCenterProjectId: projects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.ProjectName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(pointer.From(model.Location))
				state.Tags = pointer.From(model.Tags)

				identity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				state.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					state.DeploymentTargetId = pointer.From(props.DeploymentTargetId)
					state.UserRoleAssignment = flattenDevCenterProjectEnvironmentTypeUserRoleAssignment(props.UserRoleAssignments)

					if v := props.CreatorRoleAssignment; v != nil {
						state.CreatorRoleAssignmentRoles = flattenDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(v.Roles)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DevCenterProjectEnvironmentTypeResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterProjectEnvironmentTypeResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ProjectEnvironmentTypesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			payload := resp.Model

			if metadata.ResourceData.HasChange("creator_role_assignment_roles") {
				payload.Properties.CreatorRoleAssignment.Roles = expandDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(model.CreatorRoleAssignmentRoles)
			}

			if metadata.ResourceData.HasChange("deployment_target_id") {
				payload.Properties.DeploymentTargetId = pointer.To(model.DeploymentTargetId)
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return err
				}
				payload.Identity = identity
			}

			if metadata.ResourceData.HasChange("user_role_assignment") {
				userRoleAssignment, err := expandDevCenterProjectEnvironmentTypeUserRoleAssignment(model.UserRoleAssignment)
				if err != nil {
					return err
				}
				payload.Properties.UserRoleAssignments = userRoleAssignment
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.ProjectEnvironmentTypesCreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterProjectEnvironmentTypeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ProjectEnvironmentTypesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(input []string) *map[string]environmenttypes.EnvironmentRole {
	if len(input) == 0 {
		return nil
	}

	result := map[string]environmenttypes.EnvironmentRole{}

	for _, v := range input {
		result[v] = environmenttypes.EnvironmentRole{}
	}

	return &result
}

func expandDevCenterProjectEnvironmentTypeUserRoleAssignment(input []DevCenterProjectEnvironmentTypeUserRoleAssignment) (*map[string]environmenttypes.UserRoleAssignment, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := map[string]environmenttypes.UserRoleAssignment{}

	for _, v := range input {
		if _, exists := result[v.UserId]; exists {
			return nil, fmt.Errorf("`user_id` is duplicate")
		}

		result[v.UserId] = environmenttypes.UserRoleAssignment{
			Roles: expandDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(v.Roles),
		}
	}

	return &result, nil
}

func expandDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(input []string) *map[string]environmenttypes.EnvironmentRole {
	if len(input) == 0 {
		return nil
	}

	result := map[string]environmenttypes.EnvironmentRole{}

	for _, v := range input {
		result[v] = environmenttypes.EnvironmentRole{}
	}

	return &result
}

func flattenDevCenterProjectEnvironmentTypeCreatorRoleAssignmentRoles(input *map[string]environmenttypes.EnvironmentRole) []string {
	result := make([]string, 0)

	if input == nil {
		return result
	}

	for k := range *input {
		result = append(result, k)
	}

	return result
}

func flattenDevCenterProjectEnvironmentTypeUserRoleAssignment(input *map[string]environmenttypes.UserRoleAssignment) []DevCenterProjectEnvironmentTypeUserRoleAssignment {
	results := make([]DevCenterProjectEnvironmentTypeUserRoleAssignment, 0)

	if input == nil {
		return results
	}

	for k, v := range *input {
		result := DevCenterProjectEnvironmentTypeUserRoleAssignment{
			UserId: k,
			Roles:  flattenDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(v.Roles),
		}

		results = append(results, result)
	}

	return results
}

func flattenDevCenterProjectEnvironmentTypeUserRoleAssignmentRoles(input *map[string]environmenttypes.EnvironmentRole) []string {
	result := make([]string, 0)

	if input == nil {
		return result
	}

	for k := range *input {
		result = append(result, k)
	}

	return result
}
