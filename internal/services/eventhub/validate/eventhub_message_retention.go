// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "errors"

func ValidateEventHubMessageRetentionCount(v interface{}, _ string) (warnings []string, errs []error) {
	value := v.(int)

	if 90 < value || value < 1 {
		errs = append(errs, errors.New("EventHub Retention Count has to be between 1 and 7 or between 1 and 90 if using a dedicated Event Hubs Cluster"))
	}

	return warnings, errs
}
