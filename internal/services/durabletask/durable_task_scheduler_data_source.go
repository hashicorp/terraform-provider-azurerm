// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SchedulerDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	SkuName           string            `tfschema:"sku_name"`
	IpAllowList       []string          `tfschema:"ip_allow_list"`
	Capacity          int64             `tfschema:"capacity"`
	Tags              map[string]string `tfschema:"tags"`
	Endpoint          string            `tfschema:"endpoint"`
	RedundancyState   string            `tfschema:"redundancy_state"`
}

type SchedulerDataSource struct{}

var _ sdk.DataSource = SchedulerDataSource{}

func (d SchedulerDataSource) ResourceType() string {
	return "azurerm_durable_task_scheduler"
}

func (d SchedulerDataSource) ModelObject() interface{} {
	return &SchedulerDataSourceModel{}
}

func (d SchedulerDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: ValidateSchedulerName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d SchedulerDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_allow_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"redundancy_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d SchedulerDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.SchedulersClient

			var state SchedulerDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := schedulers.NewSchedulerID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.SkuName = string(props.Sku.Name)
					state.Capacity = pointer.From(props.Sku.Capacity)
					state.IpAllowList = props.IPAllowlist
					state.Endpoint = pointer.From(props.Endpoint)

					if props.Sku.RedundancyState != nil {
						state.RedundancyState = string(*props.Sku.RedundancyState)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
