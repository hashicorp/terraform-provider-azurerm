// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package suppress

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CaseDifference will be deprecated and removed in a future release.
// Rather than making the field case-insensitive (which will cause issues down the line)
// this issue can be fixed by normalizing the value being returned from the Azure API
// for example, either by using the `Parse{IDType}Insensitively` function, or by re-casing the value of the constant.
func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	// fields should be case-sensitive, normalize the Azure Resource ID in the Read if required
	return strings.EqualFold(old, new)
}
