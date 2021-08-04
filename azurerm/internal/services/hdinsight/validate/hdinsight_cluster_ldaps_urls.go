package validate

import (
	"fmt"
	"strings"
)

func HDInsightClusterLdapsUrls(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !strings.HasPrefix(v, "ldaps://") {
		errors = append(errors, fmt.Errorf(`%s should start with "ldaps://"`, k))
		return
	}

	return
}
