// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func IsARN(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	// Referencing https://github.com/aws/aws-sdk-go/blob/e8afe81156c70d5bf7b6d2ed5aeeb609ea3ba3f8/aws/arn/arn.go#L81
	if !(strings.HasPrefix(v, "arn:") && strings.Count(v, ":") >= 5) {
		errors = append(errors, fmt.Errorf("invalid ARN"))
	}
	return warnings, errors
}
