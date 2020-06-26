package logic

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func resourceArmIntegrationServiceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIntegrationServiceEnvironmentCreate,
		Read:   resourceArmIntegrationServiceEnvironmentRead,
		Update: resourceArmIntegrationServiceEnvironmentUpdate,
		Delete: resourceArmIntegrationServiceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(logic.IntegrationServiceEnvironmentSkuNameNotSpecified),
				ValidateFunc: validation.StringInSlice([]string{
					string(logic.IntegrationServiceEnvironmentSkuNameNotSpecified),
					string(logic.IntegrationServiceEnvironmentSkuNameDeveloper),
					string(logic.IntegrationServiceEnvironmentSkuNamePremium),
				}, false),
			},

			// Maximum scale units that you can add	10
			// https://docs.microsoft.com/en-US/azure/logic-apps/logic-apps-limits-and-config#integration-service-environment-ise
			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 10),
			},

			"access_endpoint_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(logic.IntegrationServiceEnvironmentAccessEndpointTypeNotSpecified),
				ValidateFunc: validation.StringInSlice([]string{
					string(logic.IntegrationServiceEnvironmentAccessEndpointTypeNotSpecified),
					string(logic.IntegrationServiceEnvironmentAccessEndpointTypeInternal),
					string(logic.IntegrationServiceEnvironmentAccessEndpointTypeExternal),
				}, false),
			},

			"virtual_network_subnet_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      set.HashStringIgnoreCase,
			},

			"connector_endpoint_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"connector_outbound_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"workflow_endpoint_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"workflow_outbound_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmIntegrationServiceEnvironmentCreate(d *schema.ResourceData, m interface{}) error {
	return resourceArmIntegrationServiceEnvironmentRead(d, m)
}

func resourceArmIntegrationServiceEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceArmIntegrationServiceEnvironmentUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceArmIntegrationServiceEnvironmentRead(d, m)
}

func resourceArmIntegrationServiceEnvironmentDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
