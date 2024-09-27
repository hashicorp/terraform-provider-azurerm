package devcenter

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/projects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = DevCenterProjectResource{}
var _ sdk.ResourceWithUpdate = DevCenterProjectResource{}

type DevCenterProjectResource struct{}

func (r DevCenterProjectResource) ModelObject() interface{} {
	return &DevCenterProjectResourceSchema{}
}

type DevCenterProjectResourceSchema struct {
	Description            string                 `tfschema:"description"`
	DevCenterId            string                 `tfschema:"dev_center_id"`
	DevCenterUri           string                 `tfschema:"dev_center_uri"`
	Location               string                 `tfschema:"location"`
	MaximumDevBoxesPerUser int64                  `tfschema:"maximum_dev_boxes_per_user"`
	Name                   string                 `tfschema:"name"`
	ResourceGroupName      string                 `tfschema:"resource_group_name"`
	Tags                   map[string]interface{} `tfschema:"tags"`
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
			client := metadata.Client.DevCenter.V20230401.Projects

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
			if err := r.mapDevCenterProjectResourceSchemaToProject(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

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
			client := metadata.Client.DevCenter.V20230401.Projects
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
				if err := r.mapProjectToDevCenterProjectResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
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
			client := metadata.Client.DevCenter.V20230401.Projects

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
			client := metadata.Client.DevCenter.V20230401.Projects

			id, err := projects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config DevCenterProjectResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload projects.ProjectUpdate
			if err := r.mapDevCenterProjectResourceSchemaToProjectUpdate(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterProjectResource) mapDevCenterProjectResourceSchemaToProject(input DevCenterProjectResourceSchema, output *projects.Project) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &projects.ProjectProperties{}
	}
	if err := r.mapDevCenterProjectResourceSchemaToProjectProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "ProjectProperties", "Properties", err)
	}

	return nil
}

func (r DevCenterProjectResource) mapProjectToDevCenterProjectResourceSchema(input projects.Project, output *DevCenterProjectResourceSchema) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &projects.ProjectProperties{}
	}
	if err := r.mapProjectPropertiesToDevCenterProjectResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "ProjectProperties", "Properties", err)
	}

	return nil
}

func (r DevCenterProjectResource) mapDevCenterProjectResourceSchemaToProjectProperties(input DevCenterProjectResourceSchema, output *projects.ProjectProperties) error {
	output.Description = &input.Description
	output.DevCenterId = input.DevCenterId

	output.MaxDevBoxesPerUser = &input.MaximumDevBoxesPerUser
	return nil
}

func (r DevCenterProjectResource) mapProjectPropertiesToDevCenterProjectResourceSchema(input projects.ProjectProperties, output *DevCenterProjectResourceSchema) error {
	output.Description = pointer.From(input.Description)
	output.DevCenterId = input.DevCenterId
	output.DevCenterUri = pointer.From(input.DevCenterUri)
	output.MaximumDevBoxesPerUser = pointer.From(input.MaxDevBoxesPerUser)
	return nil
}

func (r DevCenterProjectResource) mapDevCenterProjectResourceSchemaToProjectUpdate(input DevCenterProjectResourceSchema, output *projects.ProjectUpdate) error {
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &projects.ProjectUpdateProperties{}
	}
	if err := r.mapDevCenterProjectResourceSchemaToProjectUpdateProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "ProjectUpdateProperties", "Properties", err)
	}

	return nil
}

func (r DevCenterProjectResource) mapProjectUpdateToDevCenterProjectResourceSchema(input projects.ProjectUpdate, output *DevCenterProjectResourceSchema) error {
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &projects.ProjectUpdateProperties{}
	}
	if err := r.mapProjectUpdatePropertiesToDevCenterProjectResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "ProjectUpdateProperties", "Properties", err)
	}

	return nil
}

func (r DevCenterProjectResource) mapDevCenterProjectResourceSchemaToProjectUpdateProperties(input DevCenterProjectResourceSchema, output *projects.ProjectUpdateProperties) error {
	output.MaxDevBoxesPerUser = &input.MaximumDevBoxesPerUser
	return nil
}

func (r DevCenterProjectResource) mapProjectUpdatePropertiesToDevCenterProjectResourceSchema(input projects.ProjectUpdateProperties, output *DevCenterProjectResourceSchema) error {
	output.MaximumDevBoxesPerUser = pointer.From(input.MaxDevBoxesPerUser)
	return nil
}
