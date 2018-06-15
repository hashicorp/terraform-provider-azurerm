package validate

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
)

//todo, now in terraform helper, switch over once vended
// -> https://github.com/hashicorp/terraform/blob/master/helper/validation/validation.go#L263
func Rfc3339Time(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := date.ParseTime(time.RFC3339, v); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
	}

	return
}

func Rfc3339DateInFutureBy(d time.Duration) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (_ []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
			return
		}

		t, err := date.ParseTime(time.RFC3339, v)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
			return
		}

		if time.Until(t) < d {
			errors = append(errors, fmt.Errorf("%q is %q and should be at least %q in the future", k, i, d))
		}

		return
	}
}
