package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
)

func AppServiceCustomHostnameBindingID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.HostnameBindingID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as App Service Custom Hostname ID: %+v", k, err))
	}
	return
}
