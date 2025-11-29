// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DesktopVirtualizationHostPoolDataSource{}

type DesktopVirtualizationHostPoolDataSource struct{}

func (DesktopVirtualizationHostPoolDataSource) ModelObject() interface{} {
	return &DesktopVirtualizationHostPoolModel{}
}

func (DesktopVirtualizationHostPoolDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hostpool.ValidateHostPoolID
}

func (DesktopVirtualizationHostPoolDataSource) ResourceType() string {
	return "azurerm_virtual_desktop_host_pool"
}

func (r DesktopVirtualizationHostPoolDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DesktopVirtualizationHostPoolDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"load_balancer_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"validate_environment": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"custom_rdp_properties": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"personal_desktop_assignment_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"maximum_sessions_allowed": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"start_vm_on_connect": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"preferred_app_group_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scheduled_agent_updates": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"timezone": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"use_session_host_timezone": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"schedule": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day_of_week": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"hour_of_day": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"vm_template": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DesktopVirtualizationHostPoolDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.HostPoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DesktopVirtualizationHostPoolModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := hostpool.NewHostPoolID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.HostPoolName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props := model.Properties

				state.CustomRdpProperties = pointer.From(props.CustomRdpProperty)
				state.Description = pointer.From(props.Description)
				state.FriendlyName = pointer.From(props.FriendlyName)
				state.PublicNetworkAccess = string(pointer.From(props.PublicNetworkAccess))
				state.MaximumSessionsAllowed = pointer.From(props.MaxSessionLimit)
				state.LoadBalancerType = string(props.LoadBalancerType)
				personalDesktopAssignmentType := ""
				if props.PersonalDesktopAssignmentType != nil {
					personalDesktopAssignmentType = string(*props.PersonalDesktopAssignmentType)
				}
				state.PersonalDesktopAssignmentType = personalDesktopAssignmentType
				state.PreferredAppGroupType = string(props.PreferredAppGroupType)
				state.StartVmOnConnect = pointer.From(props.StartVMOnConnect)
				state.Type = string(props.HostPoolType)
				state.ValidateEnvironment = pointer.From(props.ValidationEnvironment)
				state.ScheduledAgentUpdates = flattenAgentUpdate(props.AgentUpdate)
				state.VmTemplate = pointer.From(props.VMTemplate)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
