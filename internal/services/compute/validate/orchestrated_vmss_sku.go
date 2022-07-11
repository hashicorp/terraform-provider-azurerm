package validate

import (
	"fmt"
	"strings"
)

func OrchestratedVirtualMachineScaleSetSku(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	skuParts := strings.Split(v, "_")

	if len(skuParts) < 2 || strings.Contains(v, "__") || strings.Contains(v, " ") {
		errors = append(errors, fmt.Errorf("%q(%q) is not formatted properly.", key, v))
	}

	return
}
