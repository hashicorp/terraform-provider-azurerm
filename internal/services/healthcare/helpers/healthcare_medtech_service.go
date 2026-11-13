// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func SuppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return pluginsdk.NormalizeJson(old) == pluginsdk.NormalizeJson(new)
}
