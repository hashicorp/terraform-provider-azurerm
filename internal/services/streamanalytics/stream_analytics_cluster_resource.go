package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ClusterResource struct{}

type ClusterModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroup     string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	StreamingCapacity int32                  `tfschema:"streaming_capacity"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = ClusterResource{}

func (r ClusterResource) ModelObject() interface{} {
	return &ClusterModel{}
}

func (r ClusterResource) ResourceType() string {
	return "azurerm_stream_analytics_cluster"
}

func (r ClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ClusterID
}

func (r ClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"streaming_capacity": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.All(
				validation.IntBetween(36, 216),
				validation.IntDivisibleBy(36),
			),
		},

		"tags": tags.Schema(),
	}
}

func (r ClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ClusterModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.ClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewClusterID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := streamanalytics.Cluster{
				Name:     utils.String(model.Name),
				Location: utils.String(model.Location),
				Sku: &streamanalytics.ClusterSku{
					Name:     streamanalytics.Default,
					Capacity: utils.Int32(model.StreamingCapacity),
				},
				Tags: tags.Expand(model.Tags),
			}

			future, err := client.CreateOrUpdate(ctx, props, id.ResourceGroup, id.Name, "", "")
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.ClustersClient
			id, err := parse.ClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := ClusterModel{
				Name:              id.Name,
				ResourceGroup:     id.ResourceGroup,
				Location:          *resp.Location,
				StreamingCapacity: *resp.Sku.Capacity,
				Tags:              tags.Flatten(resp.Tags),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.ClustersClient
			id, err := parse.ClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
				if !response.WasNotFound(resp.Response()) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r ClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.ClustersClient

			var state ClusterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("streaming_capacity") || metadata.ResourceData.HasChange("tags") {
				props := streamanalytics.Cluster{
					Sku: &streamanalytics.ClusterSku{
						Capacity: utils.Int32(state.StreamingCapacity),
					},
					Tags: tags.Expand(state.Tags),
				}

				future, err := client.Update(ctx, props, id.ResourceGroup, id.Name, "")
				if err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}

				if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("waiting for update to %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
