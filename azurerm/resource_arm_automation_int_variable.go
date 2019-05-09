package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAutomationIntVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationIntVariableCreateUpdate,
		Read:   resourceArmAutomationIntVariableRead,
		Update: resourceArmAutomationIntVariableCreateUpdate,
		Delete: resourceArmAutomationIntVariableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeInt, nil),
	}
}

func resourceArmAutomationIntVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Int")
}

func resourceArmAutomationIntVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Int")
}

func resourceArmAutomationIntVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Int")
}
