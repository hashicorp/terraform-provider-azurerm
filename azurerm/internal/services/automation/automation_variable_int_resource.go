package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutomationVariableInt() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationVariableIntCreateUpdate,
		Read:   resourceAutomationVariableIntRead,
		Update: resourceAutomationVariableIntCreateUpdate,
		Delete: resourceAutomationVariableIntDelete,

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

func resourceAutomationVariableIntCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Int")
}

func resourceAutomationVariableIntRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Int")
}

func resourceAutomationVariableIntDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Int")
}
