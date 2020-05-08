package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArmAutomationVariableInt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableIntRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeInt),
	}
}

func dataSourceArmAutomationVariableIntRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Int")
}
