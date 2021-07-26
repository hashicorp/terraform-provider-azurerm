package validate

import (
	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func QueryPoolConnectionMode() schema.SchemaValidateFunc {
	connectionModes := make([]string, len(analysisservices.PossibleConnectionModeValues()))
	for i, v := range analysisservices.PossibleConnectionModeValues() {
		connectionModes[i] = string(v)
	}

	return validation.StringInSlice(connectionModes, true)
}
