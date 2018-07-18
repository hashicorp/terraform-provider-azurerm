package azurerm

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAzureFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAzureFirewallCreate,
		Read:   resourceArmAzureFirewallRead,
		Update: resourceArmAzureFirewallCreate,
		Delete: resourceArmAzureFirewallDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"internal_public_ip_address_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAzureFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Firewall creation.")

	return resourceArmAzureFirewallRead(d, meta)
}

func resourceArmAzureFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	return nil
}

func resourceArmAzureFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	return nil
}
