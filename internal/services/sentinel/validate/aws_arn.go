// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"fmt"
	"strings"
)

func IsARN(i any, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	// Referencing https://github.com/aws/aws-sdk-go/blob/e8afe81156c70d5bf7b6d2ed5aeeb609ea3ba3f8/aws/arn/arn.go#L81
	if !strings.HasPrefix(v, "arn:") || strings.Count(v, ":") < 5 {
		errs = append(errs, errors.New("invalid ARN"))
	}
	return warnings, errs
}
