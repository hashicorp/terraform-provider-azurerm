// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualDesktopHostPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualDesktopHostPoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceVirtualDesktopHostPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hostpool.NewHostPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.HostPoolName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		props := model.Properties
		maxSessionLimit := 0
		if props.MaxSessionLimit != nil {
			maxSessionLimit = int(*props.MaxSessionLimit)
		}

		d.Set("custom_rdp_properties", props.CustomRdpProperty)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("maximum_sessions_allowed", maxSessionLimit)
		d.Set("load_balancer_type", string(props.LoadBalancerType))
		personalDesktopAssignmentType := ""
		if props.PersonalDesktopAssignmentType != nil {
			personalDesktopAssignmentType = string(*props.PersonalDesktopAssignmentType)
		}
		d.Set("personal_desktop_assignment_type", personalDesktopAssignmentType)
		d.Set("preferred_app_group_type", string(props.PreferredAppGroupType))
		d.Set("start_vm_on_connect", props.StartVMOnConnect)
		d.Set("type", string(props.HostPoolType))
		d.Set("validate_environment", props.ValidationEnvironment)
		d.Set("scheduled_agent_updates", flattenAgentUpdate(props.AgentUpdate))
	}

	return nil
}
