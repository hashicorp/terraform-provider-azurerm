package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func zonesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MinItems: 1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func singleZonesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func zoneValuesToStrings(v []interface{}) *[]string {
	zones := make([]string, 0)
	for _, zone := range v {
		zones = append(zones, zone.(string))
	}
	return &zones
}
