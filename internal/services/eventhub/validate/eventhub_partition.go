// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "errors"

func ValidateEventHubPartitionCount(v interface{}, _ string) (warnings []string, errs []error) {
	value := v.(int)

	if 1024 < value || value < 1 {
		errs = append(errs, errors.New("EventHub Partition Count has to be between 1 and 32 or between 1 and 1024 if using a dedicated Event Hubs Cluster"))
	}

	return warnings, errs
}
