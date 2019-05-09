package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmAutomationStringVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationStringVariableRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeString),
	}
}

func dataSourceArmAutomationStringVariableRead(d *schema.ResourceData, meta interface{}) error {
	return datasourceAutomationVariableRead(d, meta, "String")
}
