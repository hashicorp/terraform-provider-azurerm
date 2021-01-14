package base64

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func OptionalSchema(isVirtualMachine bool) *schema.Schema {
	// Virtual Machine's don't allow updating the Custom Data
	// Code="PropertyChangeNotAllowed" Message="Changing property 'customData' is not allowed."

	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     isVirtualMachine,
		Sensitive:    true,
		ValidateFunc: validation.StringIsBase64,
	}
}
