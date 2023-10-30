package devcenter

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/galleries"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

type DevCenterGalleryResource struct{}

//var _ sdk.ResourceWithUpdate = DevCenterGalleryResource{}

func (r DevCenterGalleryResource) ModelObject() interface{} {
	return &DevCenterGalleryResourceModel{}
}

type DevCenterGalleryResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	DevCenterName     string `tfschema:"dev_center_name"`
	GalleryResourceId string `tfschema:"gallery_resource_id"`
}

func (r DevCenterGalleryResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"gallery_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r DevCenterGalleryResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterGalleryResource) ResourceType() string {
	return "azurerm_dev_center_gallery_resource"
}

func (r DevCenterGalleryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20230401.Galleries

			var model DevCenterGalleryResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := galleries.NewGalleryID(subscriptionId, model.ResourceGroupName, model.DevCenterName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := galleries.Gallery{
				Id:   utils.String(id.ID()),
				Name: utils.String(model.Name),
				Properties: &galleries.GalleryProperties{
					GalleryResourceId: model.GalleryResourceId,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterGalleryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20230401.Galleries
			id, err := galleries.ParseGalleryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				state := DevCenterGalleryResourceModel{
					Name:              pointer.From(model.Name),
					ResourceGroupName: id.ResourceGroupName,
					DevCenterName:     id.DevCenterName,
					GalleryResourceId: props.GalleryResourceId,
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r DevCenterGalleryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20230401.Galleries
			id, err := galleries.ParseGalleryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterGalleryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return galleries.ValidateGalleryID
}

//func (r DevCenterGalleryResource) Update() sdk.ResourceFunc {
//	//TODO implement me
//	panic("implement me")
//}
