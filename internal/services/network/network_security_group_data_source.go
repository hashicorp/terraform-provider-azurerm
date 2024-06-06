// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceNetworkSecurityGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNetworkSecurityGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"security_rule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"source_port_range": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"source_port_ranges": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"destination_port_range": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"destination_port_ranges": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"source_address_prefix": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"source_address_prefixes": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"source_application_security_group_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"destination_address_prefix": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"destination_address_prefixes": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"destination_application_security_group_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"access": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"priority": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"direction": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceNetworkSecurityGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkSecurityGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := networksecuritygroups.NewNetworkSecurityGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, networksecuritygroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.SetId(id.ID())

	d.Set("name", id.NetworkSecurityGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			flattenedRules := flattenNetworkSecurityRules(props.SecurityRules)
			if err := d.Set("security_rule", flattenedRules); err != nil {
				return fmt.Errorf("setting `security_rule`: %+v", err)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
