// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsLinkedServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := linkedservices.ParseLinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	workspace := linkedservices.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.LinkedServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_id", workspace.ID())

	if model := resp.Model; model != nil {
		d.Set("read_access_id", pointer.From(model.Properties.ResourceId))

		d.Set("write_access_id", pointer.From(model.Properties.WriteAccessResourceId))
	}

	return nil
}

func LogAnalyticsLinkedServiceType(readAccessId string) string {
	if readAccessId != "" {
		return "Automation"
	}

	return "Cluster"
}
