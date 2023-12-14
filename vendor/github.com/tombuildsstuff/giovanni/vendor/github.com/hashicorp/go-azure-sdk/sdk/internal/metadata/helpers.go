// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metadata

import "strings"

func normalizeResourceId(resourceId string) string {
	return strings.TrimRight(resourceId, "/")
}
