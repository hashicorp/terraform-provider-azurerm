// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	sharedKeyWorkspaces "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	loganalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsWorkspaceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := loganalyticsClient.From(meta).SharedKeyWorkspacesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}
	sharedKeyId := sharedKeyWorkspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	permanentlyDeleteOnDestroy := features.From(meta).LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy
	if err = client.DeleteThenPoll(ctx, sharedKeyId, sharedKeyWorkspaces.DeleteOperationOptions{Force: pointer.To(permanentlyDeleteOnDestroy)}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
