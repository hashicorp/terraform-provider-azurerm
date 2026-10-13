// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package analysis_services_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAnalysisServicesServerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AnalysisServices.Servers
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	serverResp, err := client.GetDetails(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if serverResp.Model == nil {
		return fmt.Errorf("retrieving %s: response model was nil", *id)
	}
	if serverResp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", *id)
	}
	if serverResp.Model.Properties.State == nil {
		return fmt.Errorf("retrieving %s: state was nil", *id)
	}
	state := *serverResp.Model.Properties.State
	if state != servers.StateSucceeded && state != servers.StatePaused {
		return fmt.Errorf("updating %s: State must be either Succeeded or Paused but got %q", *id, string(state))
	}
	isPaused := state == servers.StatePaused

	if isPaused {
		if err := client.ResumeThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("starting %s: %+v", *id, err)
		}
	}

	analysisServicesServer := servers.AnalysisServicesServerUpdateParameters{
		Sku: &servers.ResourceSku{
			Name: d.Get("sku").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &servers.AnalysisServicesServerMutableProperties{
			AsAdministrators:        expandAnalysisServicesServerAdminUsers(d),
			IPV4FirewallSettings:    expandAnalysisServicesServerFirewallSettings(d),
			QuerypoolConnectionMode: pointer.ToEnum[servers.ConnectionMode](d.Get("querypool_connection_mode").(string)),
		},
	}

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		analysisServicesServer.Properties.BackupBlobContainerUri = pointer.To(containerUri.(string))
	}

	if err := client.UpdateThenPoll(ctx, *id, analysisServicesServer); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if isPaused {
		if err := client.SuspendThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("re-pausing %s: %+v", *id, err)
		}
	}

	return resourceAnalysisServicesServerRead(d, meta)
}
