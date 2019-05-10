package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAutomationVariableBool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationVariableBoolCreateUpdate,
		Read:   resourceArmAutomationVariableBoolRead,
		Update: resourceArmAutomationVariableBoolCreateUpdate,
		Delete: resourceArmAutomationVariableBoolDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeBool, nil),
	}
}

func resourceArmAutomationVariableBoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Bool")
}

func resourceArmAutomationVariableBoolRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Bool")
}

func resourceArmAutomationVariableBoolDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Bool")
}
