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

		Schema: AutomationVariableCommonSchemaFrom(map[string]*schema.Schema{
			"value": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		}),
	}
}

func resourceArmAutomationStringVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableCreateUpdate(d, meta, "String")
}

func resourceArmAutomationStringVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableRead(d, meta, "String")
}

func resourceArmAutomationStringVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableDelete(d, meta, "String")
}
