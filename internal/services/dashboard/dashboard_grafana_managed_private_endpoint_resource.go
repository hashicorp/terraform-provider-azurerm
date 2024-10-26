package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedPrivateEndpointResource struct{}

type ManagedPrivateEndpointModel struct {
	Name                      string            `tfschema:"name"`
	Location                  string            `tfschema:"location"`
	ManagedGrafanaId          string            `tfschema:"managed_grafana_id"`
	PrivateLinkResourceId     string            `tfschema:"private_link_resource_id"`
	PrivateLinkResourceRegion string            `tfschema:"private_link_resource_region"`
	Tags                      map[string]string `tfschema:"tags"`
	GroupIds                  []string          `tfschema:"group_ids"`
	RequestMessage            string            `tfschema:"request_message"`
}

func (r ManagedPrivateEndpointResource) ModelObject() interface{} {
	return &ManagedPrivateEndpointModel{}
}

func (r ManagedPrivateEndpointResource) ResourceType() string {
	return "azurerm_dashboard_grafana_managed_private_endpoint"
}

func (r ManagedPrivateEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedprivateendpoints.ValidateManagedPrivateEndpointID
}

func (r ManagedPrivateEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"managed_grafana_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedprivateendpoints.ValidateGrafanaID,
		},

		"private_link_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"group_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			// no validation here because an empty group id is valid for a generic private link resource
		},

		"private_link_resource_region": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"request_message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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

			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			grafanaId, err := managedprivateendpoints.ParseGrafanaID(model.ManagedGrafanaId)

			if err != nil {
				return fmt.Errorf("invalid grafana id %s: %+v", model.ManagedGrafanaId, err)
			}

			id := managedprivateendpoints.NewManagedPrivateEndpointID(subscriptionId, grafanaId.ResourceGroupName, grafanaId.GrafanaName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := managedprivateendpoints.ManagedPrivateEndpointModel{
				Location: model.Location,
				Name:     &model.Name,
				Properties: &managedprivateendpoints.ManagedPrivateEndpointModelProperties{
					GroupIds:                  &model.GroupIds,
					PrivateLinkResourceId:     &model.PrivateLinkResourceId,
					PrivateLinkResourceRegion: &model.PrivateLinkResourceRegion,
					RequestMessage:            &model.RequestMessage,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateThenPoll(ctx, id, props); err != nil {
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
			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient
			id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := ManagedPrivateEndpointModel{
				Name:     id.ManagedPrivateEndpointName,
				Location: model.Location,
			}

			if props := model.Properties; props != nil {
				if props.GroupIds != nil {
					state.GroupIds = *props.GroupIds
				}

				if props.PrivateLinkResourceId != nil {
					state.PrivateLinkResourceId = *props.PrivateLinkResourceId
				}

				if props.PrivateLinkResourceRegion != nil {
					state.PrivateLinkResourceRegion = *props.PrivateLinkResourceRegion
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedPrivateEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient
			id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagedPrivateEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient

			id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedPrivateEndpointModel
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

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

// func (r ManagedPrivateEndpointResource) StateUpgraders() sdk.StateUpgradeData {
// 	return sdk.StateUpgradeData{
// 		SchemaVersion: 1,
// 		Upgraders: map[int]pluginsdk.StateUpgrade{
// 			0: migration.StreamAnalyticsManagedPrivateEndpointV0ToV1{},
// 		},
// 	}
// }
