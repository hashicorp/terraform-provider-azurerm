package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAutomationBoolVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationBoolVariableCreateUpdate,
		Read:   resourceArmAutomationBoolVariableRead,
		Update: resourceArmAutomationBoolVariableCreateUpdate,
		Delete: resourceArmAutomationBoolVariableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeBool, nil),
	}
}

func resourceArmAutomationBoolVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Bool")
}

func resourceArmAutomationBoolVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Bool")
}

func resourceArmAutomationBoolVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Bool")
}
