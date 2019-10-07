package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func resourceArmAutomationVariableDateTime() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationVariableDateTimeCreateUpdate,
		Read:   resourceArmAutomationVariableDateTimeRead,
		Update: resourceArmAutomationVariableDateTimeCreateUpdate,
		Delete: resourceArmAutomationVariableDateTimeDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeString, validate.RFC3339Time),
	}
}

func resourceArmAutomationVariableDateTimeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Datetime")
}

func resourceArmAutomationVariableDateTimeRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Datetime")
}

func resourceArmAutomationVariableDateTimeDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Datetime")
}
