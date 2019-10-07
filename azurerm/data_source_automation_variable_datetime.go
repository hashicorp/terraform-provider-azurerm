package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArmAutomationVariableDateTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableDateTimeRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceArmAutomationVariableDateTimeRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "Datetime")
}
