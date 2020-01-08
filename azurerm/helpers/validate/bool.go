package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func BoolIsTrue() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (_ []string, errors []error) {
		v, ok := i.(bool)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be bool", k))
			return
		}

		if !v {
			errors = append(errors, fmt.Errorf("%q can only be set to true, if not required remove key", k))
			return
		}
		return
	}
}
