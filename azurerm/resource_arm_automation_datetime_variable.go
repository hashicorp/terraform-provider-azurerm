package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func resourceArmAutomationDatetimeVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationDatetimeVariableCreateUpdate,
		Read:   resourceArmAutomationDatetimeVariableRead,
		Update: resourceArmAutomationDatetimeVariableCreateUpdate,
		Delete: resourceArmAutomationDatetimeVariableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeString, validate.RFC3339Time),
	}
}

func resourceArmAutomationDatetimeVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Datetime")
}

func resourceArmAutomationDatetimeVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Datetime")
}

func resourceArmAutomationDatetimeVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Datetime")
}
