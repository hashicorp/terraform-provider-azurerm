// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// DataCollectionRuleImmutableId validates a Data Collection Rule immutable ID,
// which the Logs Ingestion API requires to identify the DCR (a value starting with `dcr-`
// followed by a GUID).
// https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal
func DataCollectionRuleImmutableId(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %q to be string", k))
	}

	if !regexp.MustCompile(`^dcr-[0-9a-fA-F]{32}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be a Data Collection Rule immutable ID in the form `dcr-<guid>`, got %q", k, v))
	}

	return
}
