// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = NetworkSecurityPerimeterDataSource{}

type NetworkSecurityPerimeterDataSource struct{}

type NetworkSecurityPerimeterDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (NetworkSecurityPerimeterDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`(^[a-zA-Z0-9]+[a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_]+$)|(^[a-zA-Z0-9]$)`),
				"`name` must be between 1 and 80 characters long, start with a letter or number, end with a letter, number, or underscore, and may contain only letters, numbers, underscores (_), periods (.), or hyphens (-).",
			),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (NetworkSecurityPerimeterDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (NetworkSecurityPerimeterDataSource) ModelObject() interface{} {
	return &NetworkSecurityPerimeterDataSourceModel{}
}

func (NetworkSecurityPerimeterDataSource) ResourceType() string {
	return "azurerm_network_security_perimeter"
}

func (NetworkSecurityPerimeterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{

		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimetersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state NetworkSecurityPerimeterDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := networksecurityperimeters.NewNetworkSecurityPerimeterID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}
