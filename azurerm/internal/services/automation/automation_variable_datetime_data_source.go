package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArmAutomationVariableDateTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableDateTimeRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceArmAutomationVariableDateTimeRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Datetime")
}
