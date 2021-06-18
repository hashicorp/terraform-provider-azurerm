package validate

import (
	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func QueryPoolConnectionMode() pluginsdk.SchemaValidateFunc {
	connectionModes := make([]string, len(analysisservices.PossibleConnectionModeValues()))
	for i, v := range analysisservices.PossibleConnectionModeValues() {
		connectionModes[i] = string(v)
	}

	return validation.StringInSlice(connectionModes, true)
}
