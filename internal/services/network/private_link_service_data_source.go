// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func dataSourcePrivateLinkService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePrivateLinkServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PrivateLinkName,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"auto_approval_subscription_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_proxy_protocol": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"visibility_subscription_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"nat_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"private_ip_address_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"primary": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"load_balancer_frontend_ip_configuration_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"alias": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourcePrivateLinkServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPrivateLinkServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.PrivateLinkServiceProperties; props != nil {
		d.Set("alias", props.Alias)
		d.Set("enable_proxy_protocol", props.EnableProxyProtocol)

		if autoApproval := props.AutoApproval; autoApproval != nil {
			if err := d.Set("auto_approval_subscription_ids", utils.FlattenStringSlice(autoApproval.Subscriptions)); err != nil {
				return fmt.Errorf("setting `auto_approval_subscription_ids`: %+v", err)
			}
		}
		if visibility := props.Visibility; visibility != nil {
			if err := d.Set("visibility_subscription_ids", utils.FlattenStringSlice(visibility.Subscriptions)); err != nil {
				return fmt.Errorf("setting `visibility_subscription_ids`: %+v", err)
			}
		}

		if props.IPConfigurations != nil {
			if err := d.Set("nat_ip_configuration", flattenPrivateLinkServiceIPConfiguration(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `nat_ip_configuration`: %+v", err)
			}
		}
		if props.LoadBalancerFrontendIPConfigurations != nil {
			if err := d.Set("load_balancer_frontend_ip_configuration_ids", dataSourceFlattenPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
				return fmt.Errorf("setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
			}
		}
	}

	d.SetId(id.ID())

	return tags.FlattenAndSet(d, resp.Tags)
}

func dataSourceFlattenPrivateLinkServiceFrontendIPConfiguration(input *[]network.FrontendIPConfiguration) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if id := item.ID; id != nil {
			results = append(results, *id)
		}
	}

	return results
}
