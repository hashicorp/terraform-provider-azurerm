package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAutomationVariableBool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomationVariableBoolRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeBool),
	}
}

func dataSourceAutomationVariableBoolRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Bool")
}
