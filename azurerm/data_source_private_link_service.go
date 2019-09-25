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

			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"auto_approval": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"fqdns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ip_configurations": {
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
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancer_frontend_ip_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_endpoint_connections": {
				Type:     schema.TypeList,
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
						"private_endpoint": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"location": azure.SchemaLocationForDataSource(),
									"tags":     tagsForDataSourceSchema(),
								},
							},
						},
						"private_link_service_connection_state": {
							Type:     schema.TypeList,
							Computed: true,
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
					},
				},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"visibility": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
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
	if privateLinkServiceProperties := resp.PrivateLinkServiceProperties; privateLinkServiceProperties != nil {
		d.Set("alias", privateLinkServiceProperties.Alias)
		if err := d.Set("auto_approval", flattenArmPrivateLinkServicePrivateLinkServicePropertiesAutoApproval(privateLinkServiceProperties.AutoApproval)); err != nil {
			return fmt.Errorf("Error setting `auto_approval`: %+v", err)
		}
		d.Set("fqdns", utils.FlattenStringSlice(privateLinkServiceProperties.Fqdns))
		if err := d.Set("ip_configurations", flattenArmPrivateLinkServicePrivateLinkServiceIPConfiguration(privateLinkServiceProperties.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configurations`: %+v", err)
		}
		if err := d.Set("load_balancer_frontend_ip_configurations", flattenArmPrivateLinkServiceFrontendIPConfiguration(privateLinkServiceProperties.LoadBalancerFrontendIPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `load_balancer_frontend_ip_configurations`: %+v", err)
		}
		if err := d.Set("network_interfaces", flattenArmPrivateLinkServiceInterface(privateLinkServiceProperties.NetworkInterfaces)); err != nil {
			return fmt.Errorf("Error setting `network_interfaces`: %+v", err)
		}
		if err := d.Set("private_endpoint_connections", flattenArmPrivateLinkServicePrivateEndpointConnection(privateLinkServiceProperties.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("Error setting `private_endpoint_connections`: %+v", err)
		}
		if err := d.Set("visibility", flattenArmPrivateLinkServicePrivateLinkServicePropertiesVisibility(privateLinkServiceProperties.Visibility)); err != nil {
			return fmt.Errorf("Error setting `visibility`: %+v", err)
		}
	}
	d.Set("type", resp.Type)

	return tags.FlattenAndSet(d, resp.Tags)
}
