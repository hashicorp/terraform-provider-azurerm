package validate

import "fmt"

const expectedID = "/providers/Microsoft.Web/sourceControls/GitHub"

func AppServiceSourceControlTokenID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}
	if v != expectedID {
		errors = append(errors, fmt.Errorf("ID must be exactly %q", expectedID))
	}
	return
}
