package validate

import (
	"fmt"
	"strings"
)

func ParameterNames(v interface{}, _ string) (warnings []string, errors []error) {
	m := v.(map[string]interface{})
	for k := range m {
		if k != strings.ToLower(k) {
			errors = append(errors, fmt.Errorf("Due to a bug in the implementation of Runbooks in Azure, the parameter names need to be specified in lowercase only. See: \"https://github.com/Azure/azure-sdk-for-go/issues/4780\" for more information."))
		}
	}

	return warnings, errors
}
