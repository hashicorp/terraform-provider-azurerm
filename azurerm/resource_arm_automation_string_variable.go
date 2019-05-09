package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func resourceArmAutomationStringVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationStringVariableCreateUpdate,
		Read:   resourceArmAutomationStringVariableRead,
		Update: resourceArmAutomationStringVariableCreateUpdate,
		Delete: resourceArmAutomationStringVariableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeString, validate.NoEmptyStrings),
	}
}

func resourceArmAutomationStringVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "String")
}

func resourceArmAutomationStringVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "String")
}

func resourceArmAutomationStringVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "String")
}
