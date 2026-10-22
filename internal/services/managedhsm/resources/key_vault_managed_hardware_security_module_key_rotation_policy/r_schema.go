// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key_rotation_policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_hsm_key_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"expire_after": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.ISO8601DurationBetween("P28D", "P100Y"),
		},

		// notify not supported in HSM Key, only rotate is supported
		"time_after_creation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601DurationBetween("P28D", "P100Y"),
			ExactlyOneOf: []string{
				"time_after_creation",
				"time_before_expiry",
			},
		},

		"time_before_expiry": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601Duration,
			ExactlyOneOf: []string{
				"time_after_creation",
				"time_before_expiry",
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
