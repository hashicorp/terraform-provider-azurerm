// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ObjectReplicationCopyBlobsCreatedAfter(i interface{}, k string) (warnings []string, errors []error) {
	// "OnlyNewObjects" and "Everything" are valid in addition to any RFC3339 date
	return validation.Any(
		validation.StringInSlice([]string{"OnlyNewObjects", "Everything"}, false),
		validation.IsRFC3339Time,
	)(i, k)
}
