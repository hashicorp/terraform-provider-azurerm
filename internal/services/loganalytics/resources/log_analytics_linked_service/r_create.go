// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_service

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsLinkedServiceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var workspaceId string
	resourceGroup := d.Get("resource_group_name").(string)
	readAccess := d.Get("read_access_id").(string)
	writeAccess := d.Get("write_access_id").(string)
	workspaceId = d.Get("workspace_id").(string)

	workspace, err := linkedservices.ParseWorkspaceID(workspaceId)
	if err != nil {
		return fmt.Errorf("linked service (Resource Group %q) unable to parse workspace id: %+v", resourceGroup, err)
	}

	id := linkedservices.NewLinkedServiceID(subscriptionId, resourceGroup, workspace.WorkspaceName, LogAnalyticsLinkedServiceType(readAccess))

	if strings.EqualFold(id.LinkedServiceName, "Cluster") && writeAccess == "" {
		return fmt.Errorf("linked service '%s/%s' (Resource Group %q): A linked Log Analytics Cluster requires the 'write_access_id' attribute to be set", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup)
	}

	if strings.EqualFold(id.LinkedServiceName, "Automation") && readAccess == "" {
		return fmt.Errorf("linked service '%s/%s' (Resource Group %q): A linked Automation Account requires the 'read_access_id' attribute to be set", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_log_analytics_linked_service", id.ID())
			}
		}
	}

	parameters := linkedservices.LinkedService{
		Properties: linkedservices.LinkedServiceProperties{},
	}

	if id.LinkedServiceName == "Automation" {
		parameters.Properties.ResourceId = pointer.To(readAccess)
	}

	if id.LinkedServiceName == "Cluster" {
		parameters.Properties.WriteAccessResourceId = pointer.To(writeAccess)
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceLogAnalyticsLinkedServiceRead(d, meta)
}
