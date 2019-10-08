package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
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

			"nat_ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_allocation_method": {
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
						"primary": {
							Type:     schema.TypeBool,
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

			"private_endpoint_connection": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_endpoint_location": azure.SchemaLocationForDataSource(),
						"state_action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_link_service_connection_state": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
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
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
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

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.PrivateLinkServiceProperties; props != nil {
		if err := d.Set("alias", props.Alias); err != nil {
			return fmt.Errorf("Error setting `alias`: %+v", err)
		}
		if props.AutoApproval != nil {
			if err := d.Set("auto_approval_subscription_ids", flattenArmPrivateLinkServicePropertiesAutoApproval(props.AutoApproval)); err != nil {
				return fmt.Errorf("Error setting `auto_approval_subscription_ids`: %+v", err)
			}
		}
		if props.Visibility != nil {
			if err := d.Set("visibility_subscription_ids", flattenArmPrivateLinkServicePropertiesVisibility(props.Visibility)); err != nil {
				return fmt.Errorf("Error setting `visibility_subscription_ids`: %+v", err)
			}
		}
		// currently not implemented yet, timeline unknown, exact purpose unknown, maybe coming to a future API near you
		// if props.Fqdns != nil {
		// 	if err := d.Set("fqdns", utils.FlattenStringSlice(props.Fqdns)); err != nil {
		// 		return fmt.Errorf("Error setting `fqdns`: %+v", err)
		// 	}
		// }
		if err := d.Set("nat_ip_configuration", flattenArmPrivateLinkServiceIPConfiguration(props.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `nat_ip_configuration`: %+v", err)
		}
		if err := d.Set("load_balancer_frontend_ip_configuration_ids", flattenArmPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
		}
		if err := d.Set("network_interface_ids", flattenArmPrivateLinkServiceInterface(props.NetworkInterfaces)); err != nil {
			return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
		}
		if connectionProps := props.PrivateEndpointConnections; connectionProps != nil {
			if err := d.Set("private_endpoint_connection", flattenArmPrivateLinkServicePrivateEndpointConnection(connectionProps)); err != nil {
				return fmt.Errorf("Error setting `private_endpoint_connection`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
