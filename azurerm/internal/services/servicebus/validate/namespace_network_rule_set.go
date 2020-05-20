package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
)

func ServiceBusNamespaceNetworkRuleSetID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ServiceBusNamespaceNetworkRuleSetID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a Service Bus Namespace Network Rule Set ID: %+v", k, err))
		return
	}

	return warnings, errors
}
