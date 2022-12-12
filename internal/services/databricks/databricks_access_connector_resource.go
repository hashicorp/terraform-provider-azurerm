package databricks

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-04-01-preview/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AccessConnectorResource struct {
}

var _ sdk.ResourceWithUpdate = AccessConnectorResource{}

type AccessConnectorResourceModel struct {
	Name          string                         `tfschema:"name"`
	ResourceGroup string                         `tfschema:"resource_group_name"`
	Location      string                         `tfschema:"location"`
	Tags          map[string]string              `tfschema:"tags"`
	Identity      []identity.ModelSystemAssigned `tfschema:"identity"`
}

func (r AccessConnectorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccessConnectorName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemAssignedIdentityRequired(),

		"tags": commonschema.Tags(),
	}
}

func (r AccessConnectorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AccessConnectorResource) ModelObject() interface{} {
	return &AccessConnectorResourceModel{}
}

func (r AccessConnectorResource) ResourceType() string {
	return "azurerm_databricks_access_connector"
}

func (r AccessConnectorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accessconnector.ValidateAccessConnectorID
}

func (r AccessConnectorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AccessConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.DataBricks.AccessConnectorClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := accessconnector.NewAccessConnectorID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			accessConnector := accessconnector.AccessConnector{
				Name:     &model.Name,
				Location: model.Location,
				Tags:     &model.Tags,
				Identity: expandedIdentity,
			}

			_, err = client.CreateOrUpdate(ctx, id, accessConnector)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r AccessConnectorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataBricks.AccessConnectorClient
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AccessConnectorResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &state.Tags
			}

			_, err = client.CreateOrUpdate(ctx, *id, *existing.Model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},

		Timeout: 5 * time.Minute,
	}
}

func (r AccessConnectorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.DataBricks.AccessConnectorClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AccessConnectorResourceModel{
				Name:          id.ConnectorName,
				Location:      location.NormalizeNilable(utils.String(resp.Model.Location)),
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				if model.Identity != nil {
					state.Identity = identity.FlattenSystemAssignedToModel(model.Identity)
				}
			}
			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r AccessConnectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())

			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.DataBricks.AccessConnectorClient

			_, err = client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
