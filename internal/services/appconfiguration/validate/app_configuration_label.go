package validate

import (
	"fmt"
	"strings"
)

func AppConfigurationLabel(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", key)}
	}

	if idx := strings.Index(v, "/Label/"); idx != -1 {
		return nil, []error{fmt.Errorf(`'/Label/' is not allowed in %q`, key)}
	}

	return nil, nil
}
