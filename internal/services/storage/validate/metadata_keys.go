// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

// MetaDataKeys validate there is no duplicate keys in the map case insensitively.
func MetaDataKeys(value any, _ string) (warnings []string, errors []error) {
	m, ok := value.(map[string]any)
	if !ok {
		return
	}

	mm := map[string]string{}
	for key := range m {
		lnkey := strings.ToLower(key)
		if okey, ok := mm[lnkey]; ok {
			errors = append(errors, fmt.Errorf("duplicated key exists: %q and %q", okey, key))
			continue
		}
		mm[lnkey] = key
	}

	return
}
