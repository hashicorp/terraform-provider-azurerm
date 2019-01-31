package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-08-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerOutboundRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Read:   resourceArmLoadBalancerOutboundRuleRead,
		Update: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Delete: resourceArmLoadBalancerOutboundRuleDelete,

		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmLoadBalancerOutboundRuleName,
			},

			"location": deprecatedLocationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"frontend_ip_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_address_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protocol": {
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: ignoreCaseStateFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Protocol1All),
					string(network.Protocol1TCP),
					string(network.Protocol1UDP),
				}, false),
			},

			"enable_tcp_reset": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"idle_timeout_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(4, 30),
			},
		},
	}
}

func resourceArmLoadBalancerOutboundRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmLoadBalancerOutboundRuleRead(d, meta)
}

func resourceArmLoadBalancerOutboundRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmLoadBalancerOutboundRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
