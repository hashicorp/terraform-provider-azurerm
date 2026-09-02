// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// PipelineGroupComponentName validates names used for a pipeline group and its named
// components (exporters, processors, receivers, pipelines and TLS configurations), which
// must be 4-33 characters, alphanumeric or hyphens, and not start or end with a hyphen.
// https://learn.microsoft.com/en-us/azure/templates/microsoft.monitor/pipelinegroups
func PipelineGroupComponentName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %q to be string", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,31}[a-zA-Z0-9]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 33 characters, contain only letters, numbers and hyphens, and must not start or end with a hyphen, got %q", k, v))
	}

	return
}
