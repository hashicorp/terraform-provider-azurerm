package streamanalytics

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ClusterResource struct{}

type ClusterModel struct {
	Name string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location string `tfschema:"location"`
	SkuName string `tfschema:"sku_name"`
	StreamingCapacity int32 `tfschema:"streaming_capacity"`
	Tags map[string]interface{} `tfschema:"tags"`
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

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"streaming_capacity": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
		},

		"tags": tags.Schema(),
	}
}

func (r ClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
	}
}

func (r ClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
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
				Name: utils.String(model.Name),
				Location: utils.String(model.Location),
				Sku: &streamanalytics.ClusterSku{
					Name: streamanalytics.ClusterSkuName(model.SkuName),
					Capacity: utils.Int32(model.StreamingCapacity),
				},
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
				Name: id.Name,
				ResourceGroup: id.ResourceGroup,
				Location: *resp.Location,
				SkuName: string(resp.Sku.Name),
				StreamingCapacity: *resp.Sku.Capacity,
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
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
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}
