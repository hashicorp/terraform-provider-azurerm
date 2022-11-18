package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LabServiceLabModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Description       string            `tfschema:"description"`
	LabPlanId         string            `tfschema:"lab_plan_id"`
	Location          string            `tfschema:"location"`
	Title             string            `tfschema:"title"`
	Tags              map[string]string `tfschema:"tags"`
}

type LabServiceLabResource struct{}

var _ sdk.ResourceWithUpdate = LabServiceLabResource{}

func (r LabServiceLabResource) ResourceType() string {
	return "azurerm_lab_service_lab"
}

func (r LabServiceLabResource) ModelObject() interface{} {
	return &LabServiceLabModel{}
}

func (r LabServiceLabResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return lab.ValidateLabID
}

func (r LabServiceLabResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LabName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.LabDescription,
		},

		"lab_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: labplan.ValidateLabPlanID,
		},

		"title": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.LabTitle,
		},

		"tags": commonschema.Tags(),
	}
}

func (r LabServiceLabResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServiceLabResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServiceLabModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LabService.LabClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := lab.NewLabID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := &lab.Lab{
				Location:   location.Normalize(model.Location),
				Properties: lab.LabProperties{},
				Tags:       &model.Tags,
			}

			if model.Description != "" {
				props.Properties.Description = &model.Description
			}

			if model.LabPlanId != "" {
				props.Properties.LabPlanId = &model.LabPlanId
			}

			if model.Title != "" {
				props.Properties.Title = &model.Title
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LabServiceLabResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServiceLabModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = &model.Description
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("lab_plan_id") {
				if model.LabPlanId != "" {
					properties.Properties.LabPlanId = &model.LabPlanId
				} else {
					properties.Properties.LabPlanId = nil
				}
			}

			if metadata.ResourceData.HasChange("title") {
				if model.Title != "" {
					properties.Properties.Title = &model.Title
				} else {
					properties.Properties.Title = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LabServiceLabResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := LabServiceLabModel{
				Name:              id.LabName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.LabPlanId != nil {
				state.LabPlanId = *properties.LabPlanId
			}

			if properties.Title != nil {
				state.Title = *properties.Title
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServiceLabResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.LabClient

			id, err := lab.ParseLabID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
