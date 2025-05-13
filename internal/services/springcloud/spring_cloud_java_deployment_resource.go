// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudJavaDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_java_deployment` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudJavaDeploymentCreate,
		Read:   resourceSpringCloudJavaDeploymentRead,
		Update: resourceSpringCloudJavaDeploymentUpdate,
		Delete: resourceSpringCloudJavaDeploymentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudJavaDeploymentV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudDeploymentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceSprintCloudJavaDeploymentSchema(),
	}
}

func resourceSpringCloudJavaDeploymentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	servicesClient := meta.(*clients.Client).AppPlatform.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSpringCloudDeploymentID(subscriptionId, appId.ResourceGroup, appId.SpringName, appId.AppName, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_spring_cloud_java_deployment", id.ID())
	}

	service, err := servicesClient.Get(ctx, appId.ResourceGroup, appId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Spring Cloud Service %q (Resource Group %q): %+v", appId.SpringName, appId.ResourceGroup, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for Spring Cloud Service %q (Resource Group %q)", appId.SpringName, appId.ResourceGroup)
	}

	deployment := appplatform.DeploymentResource{
		Sku: &appplatform.Sku{
			Name:     service.Sku.Name,
			Tier:     service.Sku.Tier,
			Capacity: utils.Int32(int32(d.Get("instance_count").(int))),
		},
		Properties: &appplatform.DeploymentResourceProperties{
			Source: appplatform.JarUploadedUserSourceInfo{
				RuntimeVersion: utils.String(d.Get("runtime_version").(string)),
				JvmOptions:     utils.String(d.Get("jvm_options").(string)),
				RelativePath:   utils.String("<default>"),
				Type:           appplatform.TypeBasicUserSourceInfoTypeJar,
			},
			DeploymentSettings: &appplatform.DeploymentSettings{
				EnvironmentVariables: expandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]interface{})),
				ResourceRequests:     expandSpringCloudDeploymentResourceRequests(d.Get("quota").([]interface{})),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSpringCloudJavaDeploymentRead(d, meta)
}

func resourceSpringCloudJavaDeploymentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudDeploymentID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("reading existing %s: %+v", id, err)
	}
	if existing.Sku == nil || existing.Properties == nil || existing.Properties.DeploymentSettings == nil {
		return fmt.Errorf("nil `sku`, `properties` or `properties.deploymentSettings` for %s: %+v", id, err)
	}

	if d.HasChange("instance_count") {
		existing.Sku.Capacity = utils.Int32(int32(d.Get("instance_count").(int)))
	}

	if d.HasChange("cpu") {
		if existing.Properties.DeploymentSettings.ResourceRequests != nil {
			existing.Properties.DeploymentSettings.ResourceRequests.CPU = utils.String(strconv.Itoa(d.Get("cpu").(int)))
		}
	}

	if d.HasChange("environment_variables") {
		existing.Properties.DeploymentSettings.EnvironmentVariables = expandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]interface{}))
	}

	if d.HasChange("jvm_options") {
		if source, ok := existing.Properties.Source.AsJarUploadedUserSourceInfo(); ok {
			source.JvmOptions = utils.String(d.Get("jvm_options").(string))
			existing.Properties.Source = source
		}
	}

	if d.HasChange("memory_in_gb") {
		if existing.Properties.DeploymentSettings.ResourceRequests != nil {
			existing.Properties.DeploymentSettings.ResourceRequests.Memory = utils.String(fmt.Sprintf("%dGi", d.Get("memory_in_gb").(int)))
		}
	}

	if d.HasChange("quota") {
		if existing.Properties.DeploymentSettings.ResourceRequests == nil {
			return fmt.Errorf("nil `properties.deploymentSettings.resourceRequests` for %s: %+v", id, err)
		}

		existing.Properties.DeploymentSettings.ResourceRequests = expandSpringCloudDeploymentResourceRequests(d.Get("quota").([]interface{}))
	}

	if d.HasChange("runtime_version") {
		if source, ok := existing.Properties.Source.AsJarUploadedUserSourceInfo(); ok {
			source.RuntimeVersion = utils.String(d.Get("runtime_version").(string))
			existing.Properties.Source = source
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName, existing)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceSpringCloudJavaDeploymentRead(d, meta)
}

func resourceSpringCloudJavaDeploymentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud deployment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / resource Group %q): %+v", id.DeploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	d.Set("name", id.DeploymentName)
	d.Set("spring_cloud_app_id", parse.NewSpringCloudAppID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName).ID())
	if resp.Sku != nil {
		d.Set("instance_count", resp.Sku.Capacity)
	}
	if resp.Properties != nil {
		if settings := resp.Properties.DeploymentSettings; settings != nil {
			d.Set("environment_variables", flattenSpringCloudDeploymentEnvironmentVariables(settings.EnvironmentVariables))
			if err := d.Set("quota", flattenSpringCloudDeploymentResourceRequests(settings.ResourceRequests)); err != nil {
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

func resourceSpringCloudJavaDeploymentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudDeploymentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / resource Group %q): %+v", id.DeploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}

func expandSpringCloudDeploymentEnvironmentVariables(envMap map[string]interface{}) map[string]*string {
	output := make(map[string]*string, len(envMap))

	for k, v := range envMap {
		output[k] = utils.String(v.(string))
	}

	return output
}

func flattenSpringCloudDeploymentEnvironmentVariables(envMap map[string]*string) map[string]interface{} {
	output := make(map[string]interface{}, len(envMap))
	for i, v := range envMap {
		if v == nil {
			continue
		}
		output[i] = *v
	}
	return output
}

func expandSpringCloudDeploymentResourceRequests(input []interface{}) *appplatform.ResourceRequests {
	cpuResult := "1"   // default value that's aligned with previous behavior used to be defined in schema.
	memResult := "1Gi" // default value that's aligned with previous behavior used to be defined in schema.

	if len(input) > 0 && input[0] != nil {
		v := input[0].(map[string]interface{})
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
		CPU:    utils.String(cpuResult),
		Memory: utils.String(memResult),
	}

	return &result
}

func flattenSpringCloudDeploymentResourceRequests(input *appplatform.ResourceRequests) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	cpu := ""
	if input.CPU != nil {
		cpu = *input.CPU
	}

	memory := ""
	if input.Memory != nil {
		memory = *input.Memory
	}

	return []interface{}{
		map[string]interface{}{
			"cpu":    cpu,
			"memory": memory,
		},
	}
}

func resourceSprintCloudJavaDeploymentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudDeploymentName,
		},

		"spring_cloud_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudAppID,
		},

		"environment_variables": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 500),
		},

		"jvm_options": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"quota": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// The value returned in GET will be recalculated by the service if the deprecated "cpu" is honored, so make this property as Computed.
					"cpu": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
						ValidateFunc: validation.StringIsNotEmpty,
					},

					// The value returned in GET will be recalculated by the service if the deprecated "memory_in_gb" is honored, so make this property as Computed.
					"memory": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"runtime_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appplatform.SupportedRuntimeValueJava8),
				string(appplatform.SupportedRuntimeValueJava11),
				string(appplatform.SupportedRuntimeValueJava17),
			}, false),
			Default: appplatform.SupportedRuntimeValueJava8,
		},
	}
}
