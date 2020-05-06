package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArmAutomationVariableString() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableStringRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceArmAutomationVariableStringRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "String")
}
