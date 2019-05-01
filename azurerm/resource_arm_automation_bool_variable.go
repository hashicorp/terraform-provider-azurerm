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

		Schema: AutomationVariableCommonSchemaFrom(map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameSchema(),

			"value": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		}),
	}
}

func resourceArmAutomationBoolVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceArmAutomationVariableCreateUpdate(d, meta, "Bool"); err != nil {
		return err
	}
	return resourceArmAutomationBoolVariableRead(d, meta)
}

func resourceArmAutomationBoolVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableRead(d, meta, "Bool")
}

func resourceArmAutomationBoolVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableDelete(d, meta, "Bool")
}
