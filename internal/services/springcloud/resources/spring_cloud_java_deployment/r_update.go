// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_java_deployment

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudJavaDeploymentUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).DeploymentsClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseDeploymentIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("reading existing %s: %+v", id, err)
	}
	if existing.Sku == nil || existing.Properties == nil || existing.Properties.DeploymentSettings == nil {
		return fmt.Errorf("nil `sku`, `properties` or `properties.deploymentSettings` for %s: %+v", id, err)
	}

	if d.HasChange("instance_count") {
		existing.Sku.Capacity = pointer.To(int32(d.Get("instance_count").(int)))
	}

	if d.HasChange("cpu") {
		if existing.Properties.DeploymentSettings.ResourceRequests != nil {
			existing.Properties.DeploymentSettings.ResourceRequests.CPU = pointer.To(strconv.Itoa(d.Get("cpu").(int)))
		}
	}

	if d.HasChange("environment_variables") {
		existing.Properties.DeploymentSettings.EnvironmentVariables = helpers.ExpandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]any))
	}

	if d.HasChange("jvm_options") {
		if source, ok := existing.Properties.Source.AsJarUploadedUserSourceInfo(); ok {
			source.JvmOptions = pointer.To(d.Get("jvm_options").(string))
			existing.Properties.Source = source
		}
	}

	if d.HasChange("memory_in_gb") {
		if existing.Properties.DeploymentSettings.ResourceRequests != nil {
			existing.Properties.DeploymentSettings.ResourceRequests.Memory = pointer.To(fmt.Sprintf("%dGi", d.Get("memory_in_gb").(int)))
		}
	}

	if d.HasChange("quota") {
		if existing.Properties.DeploymentSettings.ResourceRequests == nil {
			return fmt.Errorf("nil `properties.deploymentSettings.resourceRequests` for %s: %+v", id, err)
		}

		existing.Properties.DeploymentSettings.ResourceRequests = expandSpringCloudDeploymentResourceRequests(d.Get("quota").([]any))
	}

	if d.HasChange("runtime_version") {
		if source, ok := existing.Properties.Source.AsJarUploadedUserSourceInfo(); ok {
			source.RuntimeVersion = pointer.To(d.Get("runtime_version").(string))
			existing.Properties.Source = source
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.DeploymentName, existing)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceSpringCloudJavaDeploymentRead(d, meta)
}
