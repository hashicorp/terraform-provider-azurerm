package validate

import (
	"fmt"
	"strings"
)

func RoleResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	idParts := strings.Split(v, "roleDefinitions/")
	if len(idParts) != 2 {
		errors = append(errors, fmt.Errorf("failed to parse Role Definition ID from resource ID %q", input))
		return
	}

	return
}
