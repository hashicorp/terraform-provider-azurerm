package resourcegraph

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2022-10-01/graphquery"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceGraphQueryModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	Query             string            `tfschema:"query"`
	Description       string            `tfschema:"description"`
}

type ResourceGraphQueryResource struct{}

var (
	_ sdk.Resource           = ResourceGraphQueryResource{}
	_ sdk.ResourceWithUpdate = ResourceGraphQueryResource{}
)

func (r ResourceGraphQueryResource) ResourceType() string {
	return "azurerm_resource_graph_query"
}

func (r ResourceGraphQueryResource) ModelObject() interface{} {
	return &ResourceGraphQueryModel{}
}

func (r ResourceGraphQueryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"query": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ResourceGraphQueryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceGraphQueryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ResourceGraphQueryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ResourceGraph.ResourceGraphQueryClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := graphquery.NewQueryID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &graphquery.GraphQueryResource{
				Name:     &model.Name,
				Location: &(model.Location),
				Properties: &graphquery.GraphQueryProperties{
					Query: model.Query,
				},
				Tags: &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceGraphQueryResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ResourceGraph.ResourceGraphQueryClient

			id, err := graphquery.ParseQueryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var resourceModel ResourceGraphQueryModel
			if err := metadata.Decode(&resourceModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if metadata.ResourceData.HasChange("query") {
				model.Properties.Query = resourceModel.Query
			}

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = &resourceModel.Description
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &resourceModel.Tags
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ResourceGraphQueryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ResourceGraph.ResourceGraphQueryClient

			id, err := graphquery.ParseQueryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// var resourceModel ResourceGraphQueryModel
			// if err := metadata.Decode(&resourceModel); err != nil {
			// 	return fmt.Errorf("decoding: %+v", err)
			// }

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			props := model.Properties
			if props == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ResourceGraphQueryModel{
				Name:              id.QueryName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(*model.Location),
				Query:             props.Query,
			}

			if props.Description != nil {
				state.Description = *props.Description
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ResourceGraphQueryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ResourceGraph.ResourceGraphQueryClient

			id, err := graphquery.ParseQueryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ResourceGraphQueryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return graphquery.ValidateQueryID
}
