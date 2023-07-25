// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func WebPubSubName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{1,61}[a-zA-Z0-9]$"),
		"The web pubsub name can contain only letters, numbers and hyphens. The first character must be a letter. The last character must be a letter or number. The value must be between 3 and 63 characters long.",
	)
}

func WebPubSubHubName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z][A-Za-z0-9_`,.\\[\\]]{0,127}$"),
		"The web pubsub hub name can contain only letters, numbers and special characters including `,` , `_` , `.` , `[` ,` . The first character must be a letter. The value must be between 1 and 128 characters long.",
	)
}
