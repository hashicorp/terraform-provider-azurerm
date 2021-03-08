package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func ApplicationGatewayHTTPListenerID(i interface{}, k string) (warnings []string, errors []error) {
	if _, err := parse.ApplicationGatewayHTTPListenerID(k); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as an Application Gateway HTTP Listener ID: %v", k, err))
	}
	return
}

func ApplicationGatewayPathBasedRuleID(i interface{}, k string) (warnings []string, errors []error) {
	if _, err := parse.ApplicationGatewayPathRuleID(k); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as an Application Gateway Path Based Rule ID: %v", k, err))
	}
	return
}
