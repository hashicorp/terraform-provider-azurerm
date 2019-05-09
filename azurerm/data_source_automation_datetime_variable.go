package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmAutomationDatetimeVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationDatetimeVariableRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceArmAutomationDatetimeVariableRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "Datetime")
}
