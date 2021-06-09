package validate

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func BotMSTeamsCallingWebHook() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		value := i.(string)
		if !strings.HasPrefix(value, "https://") || !strings.HasSuffix(value, "/") {
			errors = append(errors, fmt.Errorf("invalid `calling_web_hook`, must start with `https://` and end with `/`"))
		}

		return warnings, errors
	}
}
