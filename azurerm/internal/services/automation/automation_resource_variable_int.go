package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceArmAutomationVariableInt() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationVariableIntCreateUpdate,
		Read:   resourceArmAutomationVariableIntRead,
		Update: resourceArmAutomationVariableIntCreateUpdate,
		Delete: resourceArmAutomationVariableIntDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeInt, nil),
	}
}

func resourceArmAutomationVariableIntCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Int")
}

func resourceArmAutomationVariableIntRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Int")
}

func resourceArmAutomationVariableIntDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Int")
}
