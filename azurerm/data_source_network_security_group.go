package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetworkSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetworkSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"security_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"source_port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"source_port_ranges": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"destination_port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"destination_port_ranges": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"source_address_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"source_address_prefixes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"source_application_security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"destination_address_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"destination_address_prefixes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"destination_application_security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"access": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmNetworkSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Network Security Group %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Network Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.SecurityGroupPropertiesFormat; props != nil {
		flattenedRules := flattenNetworkSecurityRules(props.SecurityRules)
		if err := d.Set("security_rule", flattenedRules); err != nil {
			return fmt.Errorf("Error flattening `security_rule`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
