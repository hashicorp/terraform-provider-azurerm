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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DatabricksAccessConnectorDataSource{}

type DatabricksAccessConnectorDataSource struct{}

type DatabricksAccessConnectorDataSourceModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
}

func (DatabricksAccessConnectorDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AccessConnectorName,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (DatabricksAccessConnectorDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),
		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),
		"tags":     commonschema.TagsDataSource(),
	}
}

func (DatabricksAccessConnectorDataSource) ModelObject() interface{} {
	return &DatabricksAccessConnectorDataSourceModel{}
}

func (DatabricksAccessConnectorDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DatabricksAccessConnectorDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.DataBricks.AccessConnectorClient

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := accessconnector.NewAccessConnectorID(subscriptionId, model.ResourceGroup, model.Name)

			resp, err := client.Get(ctx, id)
			resourceData := metadata.ResourceData
			// Handle fetch errors
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			// Set the attributes
			metadata.SetID(id)
			resourceData.Set("location", location.NormalizeNilable(&resp.Model.Location))
			identity, err := identity.FlattenLegacySystemAndUserAssignedMap(resp.Model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := resourceData.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}
			return tags.FlattenAndSet(resourceData, resp.Model.Tags)
		},
	}
}

func (DatabricksAccessConnectorDataSource) ResourceType() string {
	return "azurerm_databricks_access_connector"
}
