package base64

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func OptionalSchema(isVirtualMachine bool) *schema.Schema {
	// Virtual Machine's don't allow updating the Custom Data
	// Code="PropertyChangeNotAllowed" Message="Changing property 'customData' is not allowed."

	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     isVirtualMachine,
		ValidateFunc: validate.Base64String(),
	}
}
