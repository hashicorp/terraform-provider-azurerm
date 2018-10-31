package validate

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
)

func UUID(i interface{}, k string) (ws []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := uuid.ParseUUID(v); err != nil {
		errors = append(errors, fmt.Errorf("%q isn't a valid UUID (%q): %+v", k, v, err))
	}

	return ws, errors
}
