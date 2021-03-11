package automation

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAutomationVariableDateTime() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationVariableDateTimeCreateUpdate,
		Read:   resourceAutomationVariableDateTimeRead,
		Update: resourceAutomationVariableDateTimeCreateUpdate,
		Delete: resourceAutomationVariableDateTimeDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAutomationVariableCommonSchema(schema.TypeString, validation.IsRFC3339Time),
	}
}

func resourceAutomationVariableDateTimeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableCreateUpdate(d, meta, "Datetime")
}

func resourceAutomationVariableDateTimeRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableRead(d, meta, "Datetime")
}

func resourceAutomationVariableDateTimeDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAutomationVariableDelete(d, meta, "Datetime")
}
