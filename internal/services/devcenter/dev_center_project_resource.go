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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = DevCenterProjectResource{}
	_ sdk.ResourceWithUpdate = DevCenterProjectResource{}
)

type DevCenterProjectResource struct{}

func (r DevCenterProjectResource) ModelObject() interface{} {
	return &DevCenterProjectResourceSchema{}
}

type DevCenterProjectResourceSchema struct {
	Description            string                                     `tfschema:"description"`
	DevCenterId            string                                     `tfschema:"dev_center_id"`
	DevCenterUri           string                                     `tfschema:"dev_center_uri"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location               string                                     `tfschema:"location"`
	MaximumDevBoxesPerUser int64                                      `tfschema:"maximum_dev_boxes_per_user"`
	Name                   string                                     `tfschema:"name"`
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Tags                   map[string]interface{}                     `tfschema:"tags"`
}

func (r DevCenterProjectResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projects.ValidateProjectID
}

func (r DevCenterProjectResource) ResourceType() string {
	return "azurerm_dev_center_project"
}

func (r DevCenterProjectResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"description": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"maximum_dev_boxes_per_user": {
			Optional: true,
			Type:     pluginsdk.TypeInt,
		},
		"tags": commonschema.Tags(),
	}
}

func (r DevCenterProjectResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_uri": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r DevCenterProjectResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects

			var config DevCenterProjectResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := projects.NewProjectID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload projects.Project

			payload.Location = location.Normalize(config.Location)
			payload.Tags = tags.Expand(config.Tags)

			identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			payload.Identity = identity

			if payload.Properties == nil {
				payload.Properties = &projects.ProjectProperties{}
			}

			payload.Properties.Description = &config.Description
			payload.Properties.DevCenterId = &config.DevCenterId
			payload.Properties.MaxDevBoxesPerUser = &config.MaximumDevBoxesPerUser

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterProjectResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects
			schema := DevCenterProjectResourceSchema{}

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
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

			if model := resp.Model; model != nil {
				schema.Name = id.ProjectName
				schema.ResourceGroupName = id.ResourceGroupName
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				identity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				schema.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					schema.Description = pointer.From(props.Description)
					schema.DevCenterId = pointer.From(props.DevCenterId)
					schema.DevCenterUri = pointer.From(props.DevCenterUri)
					schema.MaximumDevBoxesPerUser = pointer.From(props.MaxDevBoxesPerUser)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r DevCenterProjectResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterProjectResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Projects

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config DevCenterProjectResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload projects.ProjectUpdate

			if metadata.ResourceData.HasChanges("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return err
				}
				payload.Identity = identity
			}

			if payload.Properties == nil {
				payload.Properties = &projects.ProjectUpdateProperties{}
			}

			if metadata.ResourceData.HasChanges("maximum_dev_boxes_per_user") {
				payload.Properties.MaxDevBoxesPerUser = pointer.To(config.MaximumDevBoxesPerUser)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
