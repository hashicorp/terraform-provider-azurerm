package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
)

func ServiceBusNamespaceID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ServiceBusNamespaceID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a Service Bus Namespace ID: %+v", k, err))
		return
	}

	return warnings, errors
}

func ServiceBusNamespaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{0,100}[a-zA-Z0-9]$").MatchString(v); !matched {
		errors = append(errors, fmt.Errorf("%q can contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number", k))
		return
	}

	return warnings, errors
}
