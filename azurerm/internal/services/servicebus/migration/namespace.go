package migration

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
)

func NamespaceV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    namespaceSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: namespaceUpgradeV0ToV1,
		Version: 0,
	}
}

func namespaceSchemaForV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"default_primary_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func namespaceUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	skuName := strings.ToLower(rawState["sku"].(string))
	premiumSku := strings.ToLower(string(servicebus.Premium))

	if skuName != premiumSku {
		delete(rawState, "capacity")
	}

	return rawState, nil
}
