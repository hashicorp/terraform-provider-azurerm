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

		Schema: AutomationVariableCommonSchemaFrom(map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameSchema(),

			"value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		}),
	}
}

func resourceArmAutomationIntVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceArmAutomationVariableCreateUpdate(d, meta, "Int"); err != nil {
		return err
	}
	return resourceArmAutomationIntVariableRead(d, meta)
}

func resourceArmAutomationIntVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableRead(d, meta, "Int")
}

func resourceArmAutomationIntVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableDelete(d, meta, "Int")
}
