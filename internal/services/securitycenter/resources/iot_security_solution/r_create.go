// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_solution

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotSecuritySolutionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.IotSecuritySolutionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))

	resourceId := parse.NewIotSecuritySolutionID(subscriptionId, resourceGroup, name).ID()
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, resourceGroup, name)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Security Center Iot Security Solution %q (Resource Group %q): %+v", name, resourceGroup, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_iot_security_solution", resourceId)
			}
		}
	}

	status := security.SolutionStatusDisabled
	if d.Get("enabled").(bool) {
		status = security.SolutionStatusEnabled
	}

	unmaskedIPLoggingStatus := security.UnmaskedIPLoggingStatusDisabled
	if d.Get("log_unmasked_ips_enabled").(bool) {
		unmaskedIPLoggingStatus = security.UnmaskedIPLoggingStatusEnabled
	}
	solution := security.IoTSecuritySolutionModel{
		Location: pointer.To(location),
		IoTSecuritySolutionProperties: &security.IoTSecuritySolutionProperties{
			DisplayName:                  pointer.To(d.Get("display_name").(string)),
			Status:                       status,
			Export:                       expandIotSecuritySolutionExport(d.Get("events_to_export").(*pluginsdk.Set).List()),
			IotHubs:                      pluginsdk.ExpandStringSlice(d.Get("iothub_ids").(*pluginsdk.Set).List()),
			RecommendationsConfiguration: expandIotSecuritySolutionRecommendation(d.Get("recommendations").([]any)),
			UnmaskedIPLoggingStatus:      unmaskedIPLoggingStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("additional_workspace"); ok {
		solution.AdditionalWorkspaces = expandIotSecuritySolutionAdditionalWorkspace(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("disabled_data_sources"); ok {
		solution.DisabledDataSources = expandIotSecuritySolutionDisabledDataSources(v.(*pluginsdk.Set).List())
	}

	logAnalyticsWorkspaceId := d.Get("log_analytics_workspace_id").(string)
	if logAnalyticsWorkspaceId != "" {
		solution.Workspace = pointer.To(logAnalyticsWorkspaceId)
	}

	query := d.Get("query_for_resources").(string)
	querySubscriptions := d.Get("query_subscription_ids").(*pluginsdk.Set).List()
	if query != "" || len(querySubscriptions) > 0 {
		if query != "" && len(querySubscriptions) > 0 {
			solution.UserDefinedResources = &security.UserDefinedResourcesProperties{
				Query:              pointer.To(query),
				QuerySubscriptions: pluginsdk.ExpandStringSlice(querySubscriptions),
			}
		} else {
			return fmt.Errorf("`query_for_resources` and `query_subscription_ids` must be set togetther")
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, solution); err != nil {
		return fmt.Errorf("creating/updating Security Center Iot Security Solution %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceIotSecuritySolutionRead(d, meta)
}
