// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// BucketPath validates that the given value is an absolute POSIX-style path used
// as the mount path inside a NetApp Files bucket.
func BucketPath(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^/`), `must be an absolute POSIX-style path starting with "/"`),
		validation.StringDoesNotMatch(regexp.MustCompile(`\\`), "must not contain backslashes"),
	)(v, k)
}
