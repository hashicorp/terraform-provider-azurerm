package databricks

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AccessConnectorResource struct {
}

var _ sdk.ResourceWithUpdate = AccessConnectorResource{}

type AccessConnectorResourceModel struct {
	Name          string            `tfschema:"name"`
	ResourceGroup string            `tfschema:"resource_group_name"`
	Location      string            `tfschema:"location"`
	Tags          map[string]string `tfschema:"tags"`
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

		// cannot use common schema 'SystemAssignedUserAssignedIdentityOptional' because
		// the Databricks implementation is slightly different as it only allows
		// 'SystemAssigned' or 'UserAssigned' (e.g. 'Only SystemAssigned or
		// UserAssigned Identity is supported for an Access Connector
		// resource, not both together.') and only allows for a single 'identity_ids'
		// to be passed...
		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(identity.TypeUserAssigned),
							string(identity.TypeSystemAssigned),
						}, false),
					},
					"identity_ids": {
						Type:     schema.TypeSet,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

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
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataBricks.AccessConnectorClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			metadata.Logger.Info("preparing arguments for AzureRM Databricks Access Connector creation")

			var model AccessConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := accessconnector.NewAccessConnectorID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Databricks %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			if identityValue.Type == identity.TypeUserAssigned && len(identityValue.IdentityIds) == 0 {
				return fmt.Errorf("`identity_ids` must be specified when `type` is set to %q", string(identity.TypeUserAssigned))
			}

			accessConnector := accessconnector.AccessConnector{
				Name:     &model.Name,
				Location: model.Location,
				Tags:     &model.Tags,
				Identity: identityValue,
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, accessConnector); err != nil {
				return fmt.Errorf("creating Databricks %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r AccessConnectorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataBricks.AccessConnectorClient
			var state AccessConnectorResourceModel

			metadata.Logger.Info("preparing arguments for AzureRM Databricks Access Connector update")

			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Databricks %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				if identityValue.Type == identity.TypeUserAssigned && len(identityValue.IdentityIds) == 0 {
					return fmt.Errorf("`identity_ids` must be specified when `type` is set to %q", string(identity.TypeUserAssigned))
				}

				existing.Model.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &state.Tags
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating Databricks %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AccessConnectorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataBricks.AccessConnectorClient

			metadata.Logger.Info("preparing arguments for AzureRM Databricks Access Connector read")

			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving Databricks %s: %+v", *id, err)
			}

			state := AccessConnectorResourceModel{
				Name:          id.AccessConnectorName,
				Location:      location.NormalizeNilable(utils.String(resp.Model.Location)),
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if model.Identity != nil {
					identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %+v", err)
					}

					if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
						return fmt.Errorf("setting `identity`: %+v", err)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AccessConnectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			metadata.Logger.Info("preparing arguments for AzureRM Databricks Access Connector deletion")

			client := metadata.Client.DataBricks.AccessConnectorClient

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Databricks %s: %+v", *id, err)
			}

			return nil
		},
	}
}
