package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceTagModel struct {
	ResourceId string `tfschema:"resource_id"`
	Key        string `tfschema:"key"`
	Value      string `tfschema:"value"`
}

type ResourceTagResource struct{}

var _ sdk.ResourceWithUpdate = ResourceTagResource{}

func (r ResourceTagResource) ResourceType() string {
	return "azurerm_resource_tag"
}

func (r ResourceTagResource) ModelObject() interface{} {
	return &ResourceTagModel{}
}

func (r ResourceTagResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateResourceTagID
}

func (r ResourceTagResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.TagKey,
		},
		"value": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.TagValue,
		},
	}
}

func (r ResourceTagResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceTagResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ResourceTagModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Resource.TagsClient
			id := parse.NewTagID(model.ResourceId, model.Key)
			existing, err := client.GetAtScope(ctx, model.ResourceId)

			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if _, ok := existing.Properties.Tags[model.Key]; ok {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			existing.Properties.Tags[model.Key] = &model.Value

			if _, err := client.CreateOrUpdateAtScope(ctx, model.ResourceId, existing); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceTagResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.TagsClient

			id, err := parse.TagID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ResourceTagModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.GetAtScope(ctx, id.ResourceId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("value") {
				resp.Properties.Tags[model.Key] = &model.Value
			}

			if _, err := client.CreateOrUpdateAtScope(ctx, id.ResourceId, resp); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ResourceTagResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.TagsClient

			id, err := parse.TagID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAtScope(ctx, id.ResourceId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if _, ok := resp.Properties.Tags[id.Key]; !ok {
				return metadata.MarkAsGone(id)
			}

			state := ResourceTagModel{
				ResourceId: id.ResourceId,
				Key:        id.Key,
				Value:      *resp.Properties.Tags[id.Key],
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ResourceTagResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.TagsClient

			id, err := parse.TagID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ResourceTagModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.GetAtScope(ctx, id.ResourceId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			delete(resp.Properties.Tags, model.Key)

			if _, err := client.CreateOrUpdateAtScope(ctx, id.ResourceId, resp); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
