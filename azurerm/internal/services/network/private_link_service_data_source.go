package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"auto_approval_subscription_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

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
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Private Link Service %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Private Link Service %q (Resource Group %q)", name, resourceGroup)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if props := resp.PrivateLinkServiceProperties; props != nil {
		d.Set("alias", props.Alias)
		d.Set("enable_proxy_protocol", props.EnableProxyProtocol)

		if autoApproval := props.AutoApproval; autoApproval != nil {
			if err := d.Set("auto_approval_subscription_ids", utils.FlattenStringSlice(autoApproval.Subscriptions)); err != nil {
				return fmt.Errorf("Error setting `auto_approval_subscription_ids`: %+v", err)
			}
		}
		if visibility := props.Visibility; visibility != nil {
			if err := d.Set("visibility_subscription_ids", utils.FlattenStringSlice(visibility.Subscriptions)); err != nil {
				return fmt.Errorf("Error setting `visibility_subscription_ids`: %+v", err)
			}
		}

		if props.IPConfigurations != nil {
			if err := d.Set("nat_ip_configuration", flattenPrivateLinkServiceIPConfiguration(props.IPConfigurations)); err != nil {
				return fmt.Errorf("Error setting `nat_ip_configuration`: %+v", err)
			}
		}
		if props.LoadBalancerFrontendIPConfigurations != nil {
			if err := d.Set("load_balancer_frontend_ip_configuration_ids", dataSourceFlattenPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
				return fmt.Errorf("Error setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
			}
		}
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

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
