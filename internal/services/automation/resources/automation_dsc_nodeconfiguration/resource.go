// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_nodeconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationDscNodeConfigurationCreate,
		Read:   resourceAutomationDscNodeConfigurationRead,
		Update: resourceAutomationDscNodeConfigurationUpdate,
		Delete: resourceAutomationDscNodeConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dscnodeconfiguration.ParseNodeConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: automationDscNodeconfigurationSchema(),
	}
}
