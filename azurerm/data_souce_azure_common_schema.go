package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"prefix": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceAzureFilterBuilder(set *schema.Set) string {
	var filters []string
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			if m["prefix"].(bool) {
				filterValues = append(filterValues, fmt.Sprintf("startswith(%s,'%s')", m["name"].(string), e))
			} else {
				filterValues = append(filterValues, fmt.Sprintf("%s eq '%s'", m["name"].(string), e))

			}
		}
		filters = append(filters, fmt.Sprintf("%s", strings.Join(filterValues[:], " or ")))

	}
	return strings.Join(filters[:], " and ")
}
