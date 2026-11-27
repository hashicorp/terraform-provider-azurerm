// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_java_deployment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudJavaDeploymentCreate(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).DeploymentsClient
	servicesClient := springcloudClient.From(meta).ServicesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	appId, err := appplatform_rm.ParseAppIDInsensitively(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := appplatform_rm.NewDeploymentID(subscriptionId, appId.ResourceGroupName, appId.SpringName, appId.AppName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_java_deployment", id.ID())
		}
	}

	service, err := servicesClient.Get(ctx, appId.ResourceGroupName, appId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Spring Cloud Service %q (Resource Group %q): %+v", appId.SpringName, appId.ResourceGroupName, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for Spring Cloud Service %q (Resource Group %q)", appId.SpringName, appId.ResourceGroupName)
	}

	deployment := appplatform.DeploymentResource{
		Sku: &appplatform.Sku{
			Name:     service.Sku.Name,
			Tier:     service.Sku.Tier,
			Capacity: pointer.To(int32(d.Get("instance_count").(int))),
		},
		Properties: &appplatform.DeploymentResourceProperties{
			Source: appplatform.JarUploadedUserSourceInfo{
				RuntimeVersion: pointer.To(d.Get("runtime_version").(string)),
				JvmOptions:     pointer.To(d.Get("jvm_options").(string)),
				RelativePath:   pointer.To("<default>"),
				Type:           appplatform.TypeBasicUserSourceInfoTypeJar,
			},
			DeploymentSettings: &appplatform.DeploymentSettings{
				EnvironmentVariables: helpers.ExpandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]any)),
				ResourceRequests:     expandSpringCloudDeploymentResourceRequests(d.Get("quota").([]any)),
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

	return resourceSpringCloudJavaDeploymentRead(d, meta)
}

func expandSpringCloudDeploymentResourceRequests(input []any) *appplatform.ResourceRequests {
	cpuResult := "1"   // default value that's aligned with previous behavior used to be defined in schema.
	memResult := "1Gi" // default value that's aligned with previous behavior used to be defined in schema.

	if len(input) > 0 && input[0] != nil {
		v := input[0].(map[string]any)
		if v != nil {
			if cpuNew := v["cpu"].(string); cpuNew != "" {
				cpuResult = cpuNew
			}

			if memoryNew := v["memory"].(string); memoryNew != "" {
				memResult = memoryNew
			}
		}
	}

	result := appplatform.ResourceRequests{
		CPU:    pointer.To(cpuResult),
		Memory: pointer.To(memResult),
	}

	return &result
}
