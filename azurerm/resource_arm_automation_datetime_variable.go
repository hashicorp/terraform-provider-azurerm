package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
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

		Schema: AutomationVariableCommonSchemaFrom(map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameSchema(),

			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceArmAutomationDatetimeVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceArmAutomationVariableCreateUpdate(d, meta, "Datetime"); err != nil {
		return err
	}
	return resourceArmAutomationDatetimeVariableRead(d, meta)
}

func resourceArmAutomationDatetimeVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableRead(d, meta, "Datetime")
}

func resourceArmAutomationDatetimeVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableDelete(d, meta, "Datetime")
}
