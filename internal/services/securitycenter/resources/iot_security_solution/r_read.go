// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_solution

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotSecuritySolutionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.IotSecuritySolutionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotSecuritySolutionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Security Center Iot Security Solution %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Security Center Iot Security Solution %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if prop := resp.IoTSecuritySolutionProperties; prop != nil {
		d.Set("enabled", prop.Status == security.SolutionStatusEnabled)
		d.Set("display_name", prop.DisplayName)
		d.Set("iothub_ids", pluginsdk.FlattenSlice(prop.IotHubs))
		d.Set("log_analytics_workspace_id", prop.Workspace)
		d.Set("log_unmasked_ips_enabled", prop.UnmaskedIPLoggingStatus == security.UnmaskedIPLoggingStatusEnabled)
		if err := d.Set("events_to_export", flattenIotSecuritySolutionExport(prop.Export)); err != nil {
			return fmt.Errorf("setting `events_to_export`: %s", err)
		}
		if err := d.Set("recommendations", flattenIotSecuritySolutionRecommendation(prop.RecommendationsConfiguration)); err != nil {
			return fmt.Errorf("setting `recommendations`: %s", err)
		}
		if prop.UserDefinedResources != nil {
			d.Set("query_for_resources", prop.UserDefinedResources.Query)
			d.Set("query_subscription_ids", pluginsdk.FlattenSlice(prop.UserDefinedResources.QuerySubscriptions))
		}
		if err := d.Set("additional_workspace", flattenIotSecuritySolutionAdditionalWorkspace(prop.AdditionalWorkspaces)); err != nil {
			return fmt.Errorf("setting `additional_workspace`: %+v", err)
		}
		if err := d.Set("disabled_data_sources", flattenIotSecuritySolutionDisabledDataSources(prop.DisabledDataSources)); err != nil {
			return fmt.Errorf("setting `disabled_data_sources`: %s", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func getRecommendationSchemaMap() map[security.RecommendationType]string {
	return map[security.RecommendationType]string{
		security.IoTACRAuthentication:             "acr_authentication",
		security.IoTAgentSendsUnutilizedMessages:  "agent_send_unutilized_msg",
		security.IoTBaseline:                      "baseline",
		security.IoTEdgeHubMemOptimize:            "edge_hub_mem_optimize",
		security.IoTEdgeLoggingOptions:            "edge_logging_option",
		security.IoTInconsistentModuleSettings:    "inconsistent_module_settings",
		security.IoTInstallAgent:                  "install_agent",
		security.IoTIPFilterDenyAll:               "ip_filter_deny_all",
		security.IoTIPFilterPermissiveRule:        "ip_filter_permissive_rule",
		security.IoTOpenPorts:                     "open_ports",
		security.IoTPermissiveFirewallPolicy:      "permissive_firewall_policy",
		security.IoTPermissiveInputFirewallRules:  "permissive_input_firewall_rules",
		security.IoTPermissiveOutputFirewallRules: "permissive_output_firewall_rules",
		security.IoTPrivilegedDockerOptions:       "privileged_docker_options",
		security.IoTSharedCredentials:             "shared_credentials",
		security.IoTVulnerableTLSCipherSuite:      "vulnerable_tls_cipher_suite",
	}
}
