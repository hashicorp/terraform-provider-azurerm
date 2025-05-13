// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynatrace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MonitorsDataSource struct{}

type MonitorsDataSourceModel struct {
	Name                          string                         `tfschema:"name"`
	ResourceGroup                 string                         `tfschema:"resource_group_name"`
	Location                      string                         `tfschema:"location"`
	MonitoringStatus              bool                           `tfschema:"monitoring_enabled"`
	MarketplaceSubscriptionStatus string                         `tfschema:"marketplace_subscription"`
	Identity                      []identity.ModelSystemAssigned `tfschema:"identity"`
	EnvironmentProperties         []EnvironmentProperties        `tfschema:"environment_properties"`
	PlanData                      []PlanData                     `tfschema:"plan"`
	UserInfo                      []UserInfo                     `tfschema:"user"`
	Tags                          map[string]string              `tfschema:"tags"`
}

type EnvironmentProperties struct {
	EnvironmentInfo []EnvironmentInfo `tfschema:"environment_info"`
}

type EnvironmentInfo struct {
	EnvironmentId string `tfschema:"environment_id"`
}

var _ sdk.DataSource = MonitorsDataSource{}

func (d MonitorsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d MonitorsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"marketplace_subscription": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedIdentityComputed(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"billing_cycle": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"plan": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"usage_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"effective_date": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"user": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"country": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"email": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"first_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"phone_number": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"environment_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"environment_info": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"environment_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (d MonitorsDataSource) ModelObject() interface{} {
	return &MonitorsDataSourceModel{}
}

func (d MonitorsDataSource) ResourceType() string {
	return "azurerm_dynatrace_monitor"
}

func (d MonitorsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var monitor MonitorsDataSourceModel
			if err := metadata.Decode(&monitor); err != nil {
				return err
			}
			id := monitors.NewMonitorID(subscriptionId, monitor.ResourceGroup, monitor.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				props := model.Properties
				identityProps, err := flattenDynatraceIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity: %+v", err)
				}
				monitoringStatus := true
				if pointer.From(props.MonitoringStatus) == monitors.MonitoringStatusDisabled {
					monitoringStatus = false
				}

				monitorResource := MonitorsDataSourceModel{
					Name:                          id.MonitorName,
					ResourceGroup:                 id.ResourceGroupName,
					Location:                      model.Location,
					MonitoringStatus:              monitoringStatus,
					MarketplaceSubscriptionStatus: string(pointer.From(props.MarketplaceSubscriptionStatus)),
					Identity:                      identityProps,
					EnvironmentProperties:         FlattenDynatraceEnvironmentProperties(props.DynatraceEnvironmentProperties),
					PlanData:                      FlattenDynatracePlanData(props.PlanData),
					UserInfo:                      FlattenDynatraceUserInfo(props.UserInfo),
				}
				if model.Tags != nil {
					monitorResource.Tags = pointer.From(model.Tags)
				}
				return metadata.Encode(&monitorResource)
			}
			return nil
		},
	}
}
