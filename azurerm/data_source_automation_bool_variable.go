package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmAutomationBoolVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationBoolVariableRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeBool),
	}
}

func dataSourceArmAutomationBoolVariableRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "Bool")
}
