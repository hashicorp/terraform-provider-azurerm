// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_java_deployment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudJavaDeploymentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
		return fmt.Errorf("reading Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / resource Group %q): %+v", id.DeploymentName, id.SpringName, id.AppName, id.ResourceGroupName, err)
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
		}
		if source, ok := resp.Properties.Source.AsJarUploadedUserSourceInfo(); ok && source != nil {
			d.Set("jvm_options", source.JvmOptions)
			d.Set("runtime_version", source.RuntimeVersion)
		}
	}

	return nil
}
