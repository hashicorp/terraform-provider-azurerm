package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateLinkServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"auto_approval_subscription_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"visibility_subscription_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// currently not implemented yet, timeline unknown, exact purpose unknown, maybe coming to a future API near you
			// "fqdns": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },

			"primary_nat_ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"secondary_nat_ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancer_frontend_ip_configuration_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_interface_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmPrivateLinkServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

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
		if props.AutoApproval.Subscriptions != nil {
			if err := d.Set("auto_approval_subscription_ids", utils.FlattenStringSlice(props.AutoApproval.Subscriptions)); err != nil {
				return fmt.Errorf("Error setting `auto_approval_subscription_ids`: %+v", err)
			}
		}
		if props.Visibility.Subscriptions != nil {
			if err := d.Set("visibility_subscription_ids", utils.FlattenStringSlice(props.Visibility.Subscriptions)); err != nil {
				return fmt.Errorf("Error setting `visibility_subscription_ids`: %+v", err)
			}
		}
		// currently not implemented yet, timeline unknown, exact purpose unknown, maybe coming to a future API near you
		// if props.Fqdns != nil {
		// 	if err := d.Set("fqdns", utils.FlattenStringSlice(props.Fqdns)); err != nil {
		// 		return fmt.Errorf("Error setting `fqdns`: %+v", err)
		// 	}
		// }
		if props.IPConfigurations != nil {
			primaryIpConfig, secondaryIpConfig := flattenArmPrivateLinkServiceIPConfiguration(props.IPConfigurations)
			if err := d.Set("primary_nat_ip_configuration", primaryIpConfig); err != nil {
				return fmt.Errorf("Error setting `primary_nat_ip_configuration`: %+v", err)
			}
			if err := d.Set("secondary_nat_ip_configuration", secondaryIpConfig); err != nil {
				return fmt.Errorf("Error setting `secondary_nat_ip_configuration`: %+v", err)
			}
		}
		if props.LoadBalancerFrontendIPConfigurations != nil {
			if err := d.Set("load_balancer_frontend_ip_configuration_ids", flattenArmPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
				return fmt.Errorf("Error setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
			}
		}
		if props.NetworkInterfaces != nil {
			if err := d.Set("network_interface_ids", flattenArmPrivateLinkServiceInterface(props.NetworkInterfaces)); err != nil {
				return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
			}
		}
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}
