// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsLinkedService() *pluginsdk.Resource {
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

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			if d.HasChange("read_access_id") {
				if readAccessID := d.Get("read_access_id").(string); readAccessID != "" {
					if _, err := automationaccount.ValidateAutomationAccountID(readAccessID, "read_acces_id"); err != nil {
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

func resourceLogAnalyticsLinkedServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Linked Services creation.")

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

	parameters := linkedservices.LinkedService{
		Properties: linkedservices.LinkedServiceProperties{},
	}

	if id.LinkedServiceName == "Automation" {
		parameters.Properties.ResourceId = utils.String(readAccess)
	}

	if id.LinkedServiceName == "Cluster" {
		parameters.Properties.WriteAccessResourceId = utils.String(writeAccess)
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
	}

	if _, err = client.Get(ctx, id); err != nil {
		return fmt.Errorf("retrieving Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceLogAnalyticsLinkedServiceRead(d, meta)
}

func resourceLogAnalyticsLinkedServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		readAccessId := ""
		if model.Properties.ResourceId != nil {
			readAccessId = *model.Properties.ResourceId
		}
		d.Set("read_access_id", readAccessId)

		writeAccessId := ""
		if model.Properties.WriteAccessResourceId != nil {
			writeAccessId = *model.Properties.WriteAccessResourceId
		}
		d.Set("write_access_id", writeAccessId)
	}

	return nil
}

func resourceLogAnalyticsLinkedServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := linkedservices.ParseLinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// (@WodansSon) - This is a bug in the service API, it returns instantly from the delete call with a 200
	// so we must wait for the state to change before we return from the delete function
	deleteWait := logAnalyticsLinkedServiceDeleteWaitForState(ctx, client, d.Timeout(pluginsdk.TimeoutDelete), *id)

	if _, err := deleteWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s: %+v", *id, err)
	}

	return nil
}

func LogAnalyticsLinkedServiceType(readAccessId string) string {
	if readAccessId != "" {
		return "Automation"
	}

	return "Cluster"
}

func resourceLogAnalyticsLinkedServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: linkedservices.ValidateWorkspaceID,
		},

		"read_access_id": {
			Type:         pluginsdk.TypeString,
			Computed:     true,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			ExactlyOneOf: []string{"read_access_id", "write_access_id"},
		},

		"write_access_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			ExactlyOneOf: []string{"read_access_id", "write_access_id"},
		},
		// Exported properties
		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
