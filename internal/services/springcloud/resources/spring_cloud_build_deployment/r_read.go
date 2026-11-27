// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_build_deployment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudBuildDeploymentRead(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).DeploymentsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseDeploymentIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud deployment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", id.DeploymentName)
	d.Set("spring_cloud_app_id", appplatform.NewAppID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName).ID())
	if resp.Sku != nil {
		d.Set("instance_count", resp.Sku.Capacity)
	}
	if resp.Properties != nil {
		if settings := resp.Properties.DeploymentSettings; settings != nil {
			d.Set("environment_variables", helpers.FlattenSpringCloudDeploymentEnvironmentVariables(settings.EnvironmentVariables))
			if err := d.Set("quota", helpers.FlattenSpringCloudDeploymentResourceRequests(settings.ResourceRequests)); err != nil {
				return fmt.Errorf("setting `quota`: %+v", err)
			}
			if err := d.Set("addon_json", helpers.FlattenSpringCloudAppAddon(settings.AddonConfigs)); err != nil {
				return fmt.Errorf("setting `addon_json`: %s", err)
			}
			apmIds, err := helpers.FlattenSpringCloudDeploymentApms(settings.Apms)
			if err != nil {
				return fmt.Errorf("setting `application_performance_monitoring_ids`: %+v", err)
			}
			d.Set("application_performance_monitoring_ids", apmIds)
		}
		if source, ok := resp.Properties.Source.AsBuildResultUserSourceInfo(); ok && source != nil {
			d.Set("build_result_id", source.BuildResultID)
		}
	}

	return nil
}
