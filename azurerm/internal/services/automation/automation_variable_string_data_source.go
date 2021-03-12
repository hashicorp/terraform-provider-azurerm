package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAutomationVariableString() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomationVariableStringRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceAutomationVariableStringRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "String")
}
