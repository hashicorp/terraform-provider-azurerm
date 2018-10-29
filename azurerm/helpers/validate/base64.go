package validate

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func Base64String() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		// Empty string is not allowed
		if s, es = validation.NoZeroValues(i, k); len(es) > 0 {
			return
		}

		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if _, err := base64.StdEncoding.DecodeString(v); err != nil {
			es = append(es, fmt.Errorf("expect value (%s) of %s is base64 string", v, k))
		}

		return
	}
}
