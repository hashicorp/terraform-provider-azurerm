package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceArmAutomationVariableBool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationVariableBoolCreateUpdate,
		Read:   resourceArmAutomationVariableBoolRead,
		Update: resourceArmAutomationVariableBoolCreateUpdate,
		Delete: resourceArmAutomationVariableBoolDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeBool, nil),
	}
}

func resourceArmAutomationVariableBoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Bool")
}

func resourceArmAutomationVariableBoolRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Bool")
}

func resourceArmAutomationVariableBoolDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Bool")
}
