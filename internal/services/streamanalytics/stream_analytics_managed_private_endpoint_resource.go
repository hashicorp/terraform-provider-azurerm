package streamanalytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedPrivateEndpointResource struct{}

type ManagedPrivateEndpointModel struct {
	Name                   string `tfschema:"name"`
	ResourceGroup          string `tfschema:"resource_group_name"`
	StreamAnalyticsCluster string `tfschema:"stream_analytics_cluster_name"`
	TargetResourceId       string `tfschema:"target_resource_id"`
	SubResourceName        string `tfschema:"subresource_name"`
}

func (r ManagedPrivateEndpointResource) ModelObject() interface{} {
	return &ManagedPrivateEndpointModel{}
}

func (r ManagedPrivateEndpointResource) ResourceType() string {
	return "azurerm_stream_analytics_managed_private_endpoint"
}

func (r ManagedPrivateEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.PrivateEndpointID
}

func (r ManagedPrivateEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"stream_analytics_cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subresource_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagedPrivateEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedPrivateEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedPrivateEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.EndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewPrivateEndpointID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsCluster, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := streamanalytics.PrivateEndpoint{
				Properties: &streamanalytics.PrivateEndpointProperties{
					ManualPrivateLinkServiceConnections: &[]streamanalytics.PrivateLinkServiceConnection{
						{
							PrivateLinkServiceConnectionProperties: &streamanalytics.PrivateLinkServiceConnectionProperties{
								PrivateLinkServiceID: utils.String(model.TargetResourceId),
								GroupIds:             &[]string{model.SubResourceName},
							},
						},
					},
				},
			}

			if _, err := client.CreateOrUpdate(ctx, props, id.ResourceGroup, id.ClusterName, id.Name, "", ""); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagedPrivateEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.EndpointsClient
			id, err := parse.PrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if resp.Properties.ManualPrivateLinkServiceConnections == nil {
				return fmt.Errorf("TODO")
			}

			state := ManagedPrivateEndpointModel{
				Name:                   id.Name,
				ResourceGroup:          id.ResourceGroup,
				StreamAnalyticsCluster: id.ClusterName,
			}

			for _, mplsc := range *resp.Properties.ManualPrivateLinkServiceConnections {
				state.TargetResourceId = *mplsc.PrivateLinkServiceID
				state.SubResourceName = strings.Join(*mplsc.GroupIds, "")
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedPrivateEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.EndpointsClient
			id, err := parse.PrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}
