// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SnapshotScheduleName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[^&%#/]{1,90}$`), `Data share snapshot schedule name should have length of 1 - 90, and cannot contain &%#/`,
	)
}
