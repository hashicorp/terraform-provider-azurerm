package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/sdk/2022-04-01/applicationinsights"
)

func WorkbookID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := applicationinsights.ParseWorkbookID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func StringDoesNotContainUpperCaseLetter(input interface{}, k string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if strings.ToLower(v) != v {
		errors = append(errors, fmt.Errorf("expected value of %s to not contain any uppercase letter", k))
		return
	}

	return
}
