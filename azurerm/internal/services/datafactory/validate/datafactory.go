package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func DataFactoryPipelineAndTriggerName() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		value := i.(string)
		if !regexp.MustCompile(`^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`).MatchString(value) {
			errors = append(errors, fmt.Errorf("invalid name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules %q: %q", k, value))
		}

		return warnings, errors
	}
}

func DataFactoryName() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		value := i.(string)
		if !regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`).MatchString(value) {
			errors = append(errors, fmt.Errorf("invalid data_factory_name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules %q: %q", k, value))
		}

		return warnings, errors
	}
}
