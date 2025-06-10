// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DevCenterNetworkConnectionDataSource{}

type DevCenterNetworkConnectionDataSource struct{}

type DevCenterNetworkConnectionDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	DomainJoinType    string            `tfschema:"domain_join_type"`
	SubnetId          string            `tfschema:"subnet_id"`
	DomainName        string            `tfschema:"domain_name"`
	DomainUsername    string            `tfschema:"domain_username"`
	OrganizationUnit  string            `tfschema:"organization_unit"`
	Tags              map[string]string `tfschema:"tags"`
}

func (DevCenterNetworkConnectionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DevCenterNetworkConnectionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"domain_join_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"domain_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"organization_unit": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DevCenterNetworkConnectionDataSource) ModelObject() interface{} {
	return &DevCenterNetworkConnectionDataSourceModel{}
}

func (DevCenterNetworkConnectionDataSource) ResourceType() string {
	return "azurerm_dev_center_network_connection"
}

func (r DevCenterNetworkConnectionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterNetworkConnectionDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := networkconnections.NewNetworkConnectionID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Name = id.NetworkConnectionName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.SubnetId = pointer.From(props.SubnetId)
					state.DomainName = pointer.From(props.DomainName)
					state.DomainUsername = pointer.From(props.DomainUsername)
					state.OrganizationUnit = pointer.From(props.OrganizationUnit)

					if v := props.DomainJoinType; v != "" {
						state.DomainJoinType = string(v)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
