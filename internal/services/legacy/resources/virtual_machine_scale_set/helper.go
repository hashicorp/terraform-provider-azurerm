// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func userDataDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return helpers.UserDataStateFunc(old) == new
}
