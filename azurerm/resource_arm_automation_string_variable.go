package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
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
			"resource_group_name": resourceGroupNameSchema(),

			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceArmAutomationStringVariableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceArmAutomationVariableCreateUpdate(d, meta, "String"); err != nil {
		return err
	}
	return resourceArmAutomationStringVariableRead(d, meta)
}

func resourceArmAutomationStringVariableRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableRead(d, meta, "String")
}

func resourceArmAutomationStringVariableDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmAutomationVariableDelete(d, meta, "String")
}
