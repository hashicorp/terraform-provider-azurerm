package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmAutomationVariableBool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableBoolRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeBool),
	}
}

func dataSourceArmAutomationVariableBoolRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "Bool")
}
