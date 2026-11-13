// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_sso_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogsinglesignonresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func datadogMonitorSsoConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"datadog_monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: datadogsinglesignonresources.ValidateMonitorID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "default",
		},

		"enterprise_application_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"single_sign_on": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(datadogsinglesignonresources.PossibleValuesForSingleSignOnStates(), false),
		},

		"login_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
