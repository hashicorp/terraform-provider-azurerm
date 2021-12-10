package validate

import (
	"encoding/pem"
	"fmt"
)

func IsCert(i interface{}, k string) (warning []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	block, _ := pem.Decode([]byte(v))
	if block == nil {
		errors = append(errors, fmt.Errorf("%s is an invalid X.509 certificate", k))
	}

	return
}
