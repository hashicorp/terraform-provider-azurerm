package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArmAutomationVariableInt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableIntRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeInt),
	}
}

func dataSourceArmAutomationVariableIntRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "Int")
}
