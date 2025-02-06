// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arcresourcebridge

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ArcResourceBridgeApplianceDataSource struct{}

var _ sdk.DataSource = ArcResourceBridgeApplianceDataSource{}

type ApplianceDataSourceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	Distro            appliances.Distro              `tfschema:"distro"`
	Identity          []identity.ModelSystemAssigned `tfschema:"identity"`
	Provider          appliances.Provider            `tfschema:"infrastructure_provider"`
	PublicKeyBase64   string                         `tfschema:"public_key_base64"`
	Tags              map[string]interface{}         `tfschema:"tags"`
}

func (r ArcResourceBridgeApplianceDataSource) ResourceType() string {
	return "azurerm_arc_resource_bridge_appliance"
}

func (r ArcResourceBridgeApplianceDataSource) ModelObject() interface{} {
	return &ApplianceDataSourceModel{}
}

func (r ArcResourceBridgeApplianceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 260),
				validation.StringMatch(regexp.MustCompile(`[^+#%&'?/,%\\]+$`), "any of '+', '#', '%', '&', ''', '?', '/', ',', '%', '&', '\\', are not allowed"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ArcResourceBridgeApplianceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"distro": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedIdentityComputed(),

		"infrastructure_provider": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_key_base64": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ArcResourceBridgeApplianceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcResourceBridge.AppliancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ApplianceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := appliances.NewApplianceID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Identity = identity.FlattenSystemAssignedToModel(model.Identity)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.Distro = pointer.From(props.Distro)
					state.PublicKeyBase64 = pointer.From(props.PublicKey)

					if infraConfig := props.InfrastructureConfig; infraConfig != nil {
						state.Provider = pointer.From(infraConfig.Provider)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
