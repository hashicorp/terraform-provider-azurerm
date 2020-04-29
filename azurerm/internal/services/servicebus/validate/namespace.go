package validate

import (
	"fmt"
	"regexp"
	"strings"

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

	if matched := regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{4,48}[a-zA-Z0-9]$").MatchString(v); !matched {
		errors = append(errors, fmt.Errorf("%q can contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number and be between 6 and 50 characters long", k))
		return
	}

	// The name cannot end with "-", "-sb" or "-mgmt".
	// See more details from link https://docs.microsoft.com/en-us/rest/api/servicebus/create-namespace.
	illegalSuffixes := []string{"-", "-sb", "-mgmt"}
	for _, illegalSuffix := range illegalSuffixes {
		if strings.HasSuffix(v, illegalSuffix) {
			errors = append(errors, fmt.Errorf("%q cannot end with a hyphen, -sb, or -mgmt", k))
			return
		}
	}

	return warnings, errors
}
