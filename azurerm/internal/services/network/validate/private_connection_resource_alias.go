package validate

import (
	"fmt"
	"strings"
)

func PrivateConnectionResourceAlias(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
		return
	}

	if !strings.HasSuffix(v, ".azure.privatelinkservice") {
		errors = append(errors, fmt.Errorf("expected %q to have suffix `.azure.privatelinkservice`", key))
	}

	return
}
