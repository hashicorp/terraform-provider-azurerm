package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// NoEmptyStrings validates that the string is not just whitespace characters (equal to [\r\n\t\f\v ])
func NoEmptyStrings() schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
		}

		newv := strings.TrimSpace(v)

		if len(newv) == 0 {
			return nil, []error{fmt.Errorf("%q must not be empty", k)}
		}

		return nil, nil
	}
}
