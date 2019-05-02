package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAutomationNullVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationNullVariableCreateUpdate,
		Read:   resourceArmAutomationNullVariableRead,
		Update: resourceArmAutomationNullVariableCreateUpdate,
		Delete: resourceArmAutomationNullVariableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: AutomationVariableCommonSchemaFrom(map[string]*schema.Schema{}),
	}
}

func resourceArmAutomationNullVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableCreateUpdate(d, meta, "Null")
}

func resourceArmAutomationNullVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableRead(d, meta, "Null")
}

func resourceArmAutomationNullVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableDelete(d, meta, "Null")
}
