// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var IQN = validation.All(
	validation.StringMatch(
		regexp.MustCompile(`iqn\.(1|2)\d{3}-(0[1-9]|1[0-2])\.[:0-9a-z-.]+\.[:0-9a-z-.]+$`),
		"IQN should follow the format `iqn.yyyy-mm.<abc>.<pqr>[:xyz]`; supported characters include [0-9a-z-.:]",
	),
	validation.StringLenBetween(4, 223),
)
