package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAutomationVariableString() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationVariableStringCreateUpdate,
		Read:   resourceAutomationVariableStringRead,
		Update: resourceAutomationVariableStringCreateUpdate,
		Delete: resourceAutomationVariableStringDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeString, validation.StringIsNotEmpty),
	}
}

func resourceAutomationVariableStringCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "String")
}

func resourceAutomationVariableStringRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "String")
}

func resourceAutomationVariableStringDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "String")
}
