package validate

import (
	"fmt"
	"strings"
)

func ConfigServerURI(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// the config server URI should be started with http://, https://, git@, or ssh://
	if !strings.HasPrefix(v, "http://") &&
		!strings.HasPrefix(v, "https://") &&
		!strings.HasPrefix(v, "git@") &&
		!strings.HasPrefix(v, "ssh://") {
		errors = append(errors, fmt.Errorf("%s should be started with http://, https://, git@, or ssh://", k))
	}
	return nil, errors
}
