// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package playwright_workspace

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Playwright Workspace.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`),
				"length of `name` must be between 3 and 24 characters (inclusive) and contain only numbers, letters, and hyphens (-)"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dataplane_uri": {
			Description: "The data plane service API URI of the Playwright Workspace. The format of URI is `https://{location}.api.playwright.microsoft.com/playwrightworkspaces/{uuid}`. To use the Playwright Workspace remote browsers, apply the URI according to the [documentation](https://learn.microsoft.com/azure/app-testing/playwright-workspaces/quickstart-run-end-to-end-tests?tabs=playwrightcli&pivots=playwright-test-runner#configure-the-browser-endpoint).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"uuid": {
			Description: "The ID in GUID format of the Playwright Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
