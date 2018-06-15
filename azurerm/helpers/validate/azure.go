package validate

import (
	"fmt"
	"regexp"
)

func AzureResourceId(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	//todo, match subscription/{guid} and provider/name.name perhaps?
	if matched := regexp.MustCompile(`(/subscriptions/)|(/providers/)`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("azure resource ID %q should start with with '/subscriptions/{subscriptionId}' or '/providers/{resourceProviderNamespace}/'", k))
	}

	return
}
