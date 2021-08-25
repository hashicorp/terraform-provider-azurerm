package validate

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices/sdk/servers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func QueryPoolConnectionMode() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(servers.ConnectionModeAll),
		string(servers.ConnectionModeReadOnly),
	}, true)
}
