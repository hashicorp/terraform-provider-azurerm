// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_build_deployment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudBuildDeploymentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	servicesClient := meta.(*clients.Client).AppPlatform.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	appId, err := appplatform_rm.ParseAppIDInsensitively(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := appplatform_rm.NewDeploymentID(subscriptionId, appId.ResourceGroupName, appId.SpringName, appId.AppName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_build_deployment", id.ID())
			}
		}
	}

	service, err := servicesClient.Get(ctx, appId.ResourceGroupName, appId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", appId, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for Spring Cloud Service %q (Resource Group %q)", appId.SpringName, appId.ResourceGroupName)
	}

	addonConfig, err := helpers.ExpandSpringCloudAppAddon(d.Get("addon_json").(string))
	if err != nil {
		return err
	}

	deployment := appplatform.DeploymentResource{
		Sku: &appplatform.Sku{
			Name:     service.Sku.Name,
			Tier:     service.Sku.Tier,
			Capacity: pointer.To(int32(d.Get("instance_count").(int))),
		},
		Properties: &appplatform.DeploymentResourceProperties{
			Source: appplatform.BuildResultUserSourceInfo{
				BuildResultID: pointer.To(d.Get("build_result_id").(string)),
			},
			DeploymentSettings: &appplatform.DeploymentSettings{
				AddonConfigs:         addonConfig,
				Apms:                 helpers.ExpandSpringCloudDeploymentApms(d.Get("application_performance_monitoring_ids").([]any)),
				EnvironmentVariables: helpers.ExpandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]any)),
				ResourceRequests:     expandSpringCloudBuildDeploymentResourceRequests(d.Get("quota").([]any)),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	return resourceSpringCloudBuildDeploymentRead(d, meta)
}

func expandSpringCloudBuildDeploymentResourceRequests(input []any) *appplatform.ResourceRequests {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	cpuResult := "1"
	memResult := "1Gi"

	v := input[0].(map[string]any)
	if cpuNew := v["cpu"].(string); cpuNew != "" {
		cpuResult = cpuNew
	}

	if memoryNew := v["memory"].(string); memoryNew != "" {
		memResult = memoryNew
	}

	result := appplatform.ResourceRequests{
		CPU:    pointer.To(cpuResult),
		Memory: pointer.To(memResult),
	}

	return &result
}
