package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/applicationinsights"
)

func WorkbookTemplateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := applicationinsights.ParseWorkbookTemplateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
