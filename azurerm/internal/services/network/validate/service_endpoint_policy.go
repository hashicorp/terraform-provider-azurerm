package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func ServiceEndpointPolicyID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ServiceEndpointPolicyID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func ServiceEndpointPolicyDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ServiceEndpointPolicyDefinitionID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func ServiceEndpointPolicyName(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_]([\w.\-]{0,78}[\w])?$`), "The name can be up to 80 characters long. It must begin with a alphnum character, and it must end with a alphnum character or with '_'. The name may contain alphnum characters or '.', '-', '_'.")(i, k)
}
