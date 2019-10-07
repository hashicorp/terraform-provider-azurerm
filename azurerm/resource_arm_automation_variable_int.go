package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceArmAutomationVariableInt() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationVariableIntCreateUpdate,
		Read:   resourceArmAutomationVariableIntRead,
		Update: resourceArmAutomationVariableIntCreateUpdate,
		Delete: resourceArmAutomationVariableIntDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeInt, nil),
	}
}

func resourceArmAutomationVariableIntCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Int")
}

func resourceArmAutomationVariableIntRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Int")
}

func resourceArmAutomationVariableIntDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Int")
}
