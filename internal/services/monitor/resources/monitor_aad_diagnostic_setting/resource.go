// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_aad_diagnostic_setting

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorAADDiagnosticSettingCreate,
		Read:   resourceMonitorAADDiagnosticSettingRead,
		Update: resourceMonitorAADDiagnosticSettingUpdate,
		Delete: resourceMonitorAADDiagnosticSettingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := diagnosticsettings.ParseDiagnosticSettingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: monitorAadDiagnosticSettingSchema(),
	}
}
