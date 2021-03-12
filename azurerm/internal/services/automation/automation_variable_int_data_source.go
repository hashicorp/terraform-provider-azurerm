package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAutomationVariableInt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomationVariableIntRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeInt),
	}
}

func dataSourceAutomationVariableIntRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Int")
}
