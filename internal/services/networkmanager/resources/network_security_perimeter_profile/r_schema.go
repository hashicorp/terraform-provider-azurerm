// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_profile

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`(^[a-zA-Z0-9]+[a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_]+$)|(^[a-zA-Z0-9]$)`),
				"`name` must be between 1 and 80 characters long, start with a letter or number, end with a letter, number, or underscore, and may contain only letters, numbers, underscores (_), periods (.), or hyphens (-).",
			),
			ForceNew: true,
		},

		"network_security_perimeter_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networksecurityperimeterprofiles.ValidateNetworkSecurityPerimeterID,
			ForceNew:     true,
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
