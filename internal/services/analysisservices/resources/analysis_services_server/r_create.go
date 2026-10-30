// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package analysis_services_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAnalysisServicesServerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AnalysisServices.Servers
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		server, err := client.GetDetails(ctx, id)
		if err != nil {
			if !response.WasNotFound(server.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(server.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_analysis_services_server", id.ID())
		}
	}

	analysisServicesServer := servers.AnalysisServicesServer{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: servers.ResourceSku{
			Name: d.Get("sku").(string),
		},
		Properties: &servers.AnalysisServicesServerProperties{
			AsAdministrators:     expandAnalysisServicesServerAdminUsers(d),
			IPV4FirewallSettings: expandAnalysisServicesServerFirewallSettings(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if querypoolConnectionMode, ok := d.GetOk("querypool_connection_mode"); ok {
		analysisServicesServer.Properties.QuerypoolConnectionMode = pointer.ToEnum[servers.ConnectionMode](querypoolConnectionMode.(string))
	}

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		analysisServicesServer.Properties.BackupBlobContainerUri = pointer.To(containerUri.(string))
	}

	if err := client.CreateCallbackThenPoll(ctx, id, analysisServicesServer, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceAnalysisServicesServerRead(d, meta)
}

func expandAnalysisServicesServerAdminUsers(d *pluginsdk.ResourceData) *servers.ServerAdministrators {
	adminUsers := d.Get("admin_users").(*pluginsdk.Set)
	members := make([]string, 0)

	for _, admin := range adminUsers.List() {
		if adm, ok := admin.(string); ok {
			members = append(members, adm)
		}
	}

	return &servers.ServerAdministrators{
		Members: &members,
	}
}

func expandAnalysisServicesServerFirewallSettings(d *pluginsdk.ResourceData) *servers.IPv4FirewallSettings {
	fwRules := make([]servers.IPv4FirewallRule, 0)
	result := servers.IPv4FirewallSettings{}

	if !pluginsdk.IsExplicitlyNullInConfig(d, "power_bi_service_enabled") {
		result.EnablePowerBIService = pointer.To(d.Get("power_bi_service_enabled").(bool))
		// when `power_bi_service_enabled` is specified, we must send at least an empty array for `FirewallRules`
		// otherwise the API errors out with a 400.
		result.FirewallRules = &fwRules
	}

	firewallRules := d.Get("ipv4_firewall_rule").(*pluginsdk.Set).List()
	if len(firewallRules) == 0 {
		return &result
	}

	for _, v := range firewallRules {
		fwRule := v.(map[string]any)
		fwRules = append(fwRules, servers.IPv4FirewallRule{
			FirewallRuleName: pointer.To(fwRule["name"].(string)),
			RangeStart:       pointer.To(fwRule["range_start"].(string)),
			RangeEnd:         pointer.To(fwRule["range_end"].(string)),
		})
	}
	result.FirewallRules = &fwRules

	return &result
}
