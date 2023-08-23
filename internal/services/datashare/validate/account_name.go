// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AccountName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[^<>%&:\\?/#*$^();,.\|+={}\[\]!~@]{3,90}$`), `Data share account name should have length of 3 - 90, and cannot contain <>%&:\?/#*$^();,.|+={}[]!~@.`,
	)
}
