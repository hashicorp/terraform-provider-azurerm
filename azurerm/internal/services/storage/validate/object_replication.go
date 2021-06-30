package validate

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
)

func ObjectReplicationCopyBlobsCreatedAfter(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "OnlyNewObjects" || v == "Everything" {
		return warnings, errors
	}

	_, err := date.ParseTime(time.RFC3339, v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
		return warnings, errors
	}

	return warnings, errors
}
