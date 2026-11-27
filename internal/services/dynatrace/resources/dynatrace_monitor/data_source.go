// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	dynatraceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/dynatrace/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

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

var _ sdk.DataSource = DataSource{}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "Name of the Dynatrace monitor.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"monitoring_enabled": {
			Description: "Flag specifying if the resource monitoring is enabled or disabled.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"marketplace_subscription": {
			Description: "Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemAssignedIdentityComputed(),

		"plan": {
			Description: "Billing plan information.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"billing_cycle": {
						Description: "Different billing cycles.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"plan": {
						Description: "Plan ID as published by Dynatrace.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"usage_type": {
						Description: "Different usage type.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"effective_date": {
						Description: "Date when plan was applied.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"user": {
			Description: "User's information.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"country": {
						Description: "Country of the user.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"email": {
						Description: "Email of the user used by Dynatrace for contacting them if needed.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"first_name": {
						Description: "First name of the user.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"last_name": {
						Description: "Last name of the user.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"phone_number": {
						Description: "Phone number of the user by Dynatrace for contacting them if needed.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataSource) ModelObject() any {
	return &MonitorsDataSourceModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_dynatrace_monitor"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Dynatrace Monitor."
}

func (d DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := dynatraceClient.From(metadata.Client).MonitorsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

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
					MarketplaceSubscriptionStatus: pointer.FromEnum(props.MarketplaceSubscriptionStatus),
					Identity:                      identityProps,
					EnvironmentProperties:         FlattenDynatraceEnvironmentProperties(props.DynatraceEnvironmentProperties),
					PlanData:                      FlattenDynatracePlanData(props.PlanData),
					UserInfo:                      FlattenDynatraceUserInfo(props.UserInfo),
				}
				if model.Tags != nil {
					monitorResource.Tags = pointer.From(model.Tags)
				}
				metadata.SetID(id)
				return metadata.Encode(&monitorResource)
			}
			return nil
		},
	}
}
