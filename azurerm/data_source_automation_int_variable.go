package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmAutomationIntVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationIntVariableRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeInt),
	}
}

func dataSourceArmAutomationIntVariableRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "Int")
}
