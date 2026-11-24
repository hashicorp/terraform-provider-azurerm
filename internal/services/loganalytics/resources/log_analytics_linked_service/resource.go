// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsLinkedServiceCreateUpdate,
		Read:   resourceLogAnalyticsLinkedServiceRead,
		Update: resourceLogAnalyticsLinkedServiceCreateUpdate,
		Delete: resourceLogAnalyticsLinkedServiceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := linkedservices.ParseLinkedServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceLogAnalyticsLinkedServiceSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			if d.HasChange("read_access_id") {
				if readAccessID := d.Get("read_access_id").(string); readAccessID != "" {
					if _, err := automationaccount.ValidateAutomationAccountID(readAccessID, "read_access_id"); err != nil {
						return fmt.Errorf("'read_access_id' must be an Automation Account resource ID, got %q", readAccessID)
					}
				}
			}

			if d.HasChange("write_access_id") {
				if writeAccessID := d.Get("write_access_id").(string); writeAccessID != "" {
					if _, err := clusters.ValidateClusterID(writeAccessID, "write_access_id"); err != nil {
						return fmt.Errorf("'write_access_id' must be a Log Analytics Cluster resource ID, got %q", writeAccessID)
					}
				}
			}

			return nil
		}),
	}
}
