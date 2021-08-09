package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/sdk/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func QueryPoolConnectionMode() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(servers.ConnectionModeAll),
		string(servers.ConnectionModeReadOnly),
	}, true)
}
