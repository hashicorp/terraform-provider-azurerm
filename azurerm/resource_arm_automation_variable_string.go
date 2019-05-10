package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func resourceArmAutomationVariableString() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationVariableStringCreateUpdate,
		Read:   resourceArmAutomationVariableStringRead,
		Update: resourceArmAutomationVariableStringCreateUpdate,
		Delete: resourceArmAutomationVariableStringDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeString, validate.NoEmptyStrings),
	}
}

func resourceArmAutomationVariableStringCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "String")
}

func resourceArmAutomationVariableStringRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "String")
}

func resourceArmAutomationVariableStringDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "String")
}
