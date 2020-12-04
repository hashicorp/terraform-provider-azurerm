package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

func VirtualMachineScaleSetID(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	id, err := parse.VirtualMachineScaleSetID(v)
	if err != nil {
		es = append(es, fmt.Errorf("Error parsing %q as a VM Scale Set Resource ID: %s", v, err))
		return
	}

	if id.Name == "" {
		es = append(es, fmt.Errorf("Error parsing %q as a VM Scale Set Resource ID: `virtualMachineScaleSets` segment was empty", v))
		return
	}

	return
}
