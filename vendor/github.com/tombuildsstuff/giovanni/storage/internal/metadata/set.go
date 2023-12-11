// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metadata

import "fmt"

// SetIntoHeaders sets the provided MetaData into the headers
func SetIntoHeaders(headers map[string]interface{}, metaData map[string]string) map[string]interface{} {
	for k, v := range metaData {
		key := fmt.Sprintf("x-ms-meta-%s", k)
		headers[key] = v
	}

	return headers
}
