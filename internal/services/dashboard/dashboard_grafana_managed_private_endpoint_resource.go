package dashboard

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/grafanaresource"
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
	GrafanaId                 string            `tfschema:"grafana_id"`
	PrivateLinkResourceId     string            `tfschema:"private_link_resource_id"`
	PrivateLinkResourceRegion string            `tfschema:"private_link_resource_region"`
	Tags                      map[string]string `tfschema:"tags"`
	GroupIds                  []string          `tfschema:"group_ids"`
	RequestMessage            string            `tfschema:"request_message"`
}

type ManagedPrivateEndpointId struct {
	SubscriptionId             string
	ResourceGroupName          string
	GrafanaName                string
	ManagedPrivateEndpointName string
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`\A([a-zA-Z]{1}[a-zA-Z0-9\-]{1,19}[a-zA-Z0-9]{1})\z`),
				`Name length can only consist of alphanumeric characters or dashes, and must be between 2 and 20 characters long. It must begin with a letter and end with a letter or digit.`,
			),
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"grafana_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: grafanaresource.ValidateGrafanaID,
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
			grafanaId, err := grafanaresource.ParseGrafanaID(model.GrafanaId)
			if err != nil {
				return err
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
				Location: location.Normalize(model.Location),
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
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			grafanaId := grafanaresource.NewGrafanaID(id.SubscriptionId, id.ResourceGroupName, id.GrafanaName)
			state := ManagedPrivateEndpointModel{
				Name:      id.ManagedPrivateEndpointName,
				GrafanaId: grafanaId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.GroupIds = pointer.From(props.GroupIds)
					state.PrivateLinkResourceId = pointer.From(props.PrivateLinkResourceId)
					state.PrivateLinkResourceRegion = pointer.From(props.PrivateLinkResourceRegion)
					state.RequestMessage = pointer.From(props.RequestMessage)
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

			var mpe ManagedPrivateEndpointModel
			if err := metadata.Decode(&mpe); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &mpe.Tags
			}

			if err := client.CreateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
