package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAutomationVariableDateTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomationVariableDateTimeRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceAutomationVariableDateTimeRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Datetime")
}
