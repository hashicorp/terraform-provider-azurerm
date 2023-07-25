// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

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
			var model AccessConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.DataBricks.AccessConnectorClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := accessconnector.NewAccessConnectorID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Databricks %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			accessConnector := accessconnector.AccessConnector{
				Name:     &model.Name,
				Location: model.Location,
				Tags:     &model.Tags,
				Identity: expandedIdentity,
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
				return fmt.Errorf("reading Databricks %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("identity") {
				// TODO: Switch this to 'identity.ExpandSystemOrSingleUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))'
				// once SDK Helpers PR #164 has been merged and integrated into the provider...
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
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
			id, err := accessconnector.ParseAccessConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.DataBricks.AccessConnectorClient

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

			client := metadata.Client.DataBricks.AccessConnectorClient

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Databricks %s: %+v", *id, err)
			}

			return nil
		},
	}
}
