// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_workspace/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsWorkspaceCreate,
		Read:   resourceLogAnalyticsWorkspaceRead,
		Update: resourceLogAnalyticsWorkspaceUpdate,
		Delete: resourceLogAnalyticsWorkspaceDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceLogAnalyticsWorkspaceCustomDiff),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&workspaces.WorkspaceId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&workspaces.WorkspaceId{}),

		SchemaVersion: 3,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceV0ToV1{},
			1: migration.WorkspaceV1ToV2{},
			2: migration.WorkspaceV2ToV3{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: logAnalyticsWorkspaceSchema(),
	}
}
