// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"time"

	appplatform2 "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
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

func resourceSpringCloudBuildDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_build_deployment` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudBuildDeploymentCreateUpdate,
		Read:   resourceSpringCloudBuildDeploymentRead,
		Update: resourceSpringCloudBuildDeploymentCreateUpdate,
		Delete: resourceSpringCloudBuildDeploymentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudBuildDeploymentV0ToV1{},
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

		Schema: map[string]*pluginsdk.Schema{
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

			"build_result_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"application_performance_monitoring_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: appplatform2.ValidateApmID,
				},
			},

			"addon_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"environment_variables": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"instance_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 500),
			},

			"quota": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cpu": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
							ValidateFunc: validation.StringIsNotEmpty,
						},

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
		},
	}
}

func resourceSpringCloudBuildDeploymentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_build_deployment", id.ID())
		}
	}

	service, err := servicesClient.Get(ctx, appId.ResourceGroup, appId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", appId, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for Spring Cloud Service %q (Resource Group %q)", appId.SpringName, appId.ResourceGroup)
	}

	addonConfig, err := expandSpringCloudAppAddon(d.Get("addon_json").(string))
	if err != nil {
		return err
	}

	deployment := appplatform.DeploymentResource{
		Sku: &appplatform.Sku{
			Name:     service.Sku.Name,
			Tier:     service.Sku.Tier,
			Capacity: utils.Int32(int32(d.Get("instance_count").(int))),
		},
		Properties: &appplatform.DeploymentResourceProperties{
			Source: appplatform.BuildResultUserSourceInfo{
				BuildResultID: utils.String(d.Get("build_result_id").(string)),
			},
			DeploymentSettings: &appplatform.DeploymentSettings{
				AddonConfigs:         addonConfig,
				Apms:                 expandSpringCloudDeploymentApms(d.Get("application_performance_monitoring_ids").([]interface{})),
				EnvironmentVariables: expandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]interface{})),
				ResourceRequests:     expandSpringCloudBuildDeploymentResourceRequests(d.Get("quota").([]interface{})),
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

	return resourceSpringCloudBuildDeploymentRead(d, meta)
}

func resourceSpringCloudBuildDeploymentRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("reading %s: %+v", id, err)
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
			if err := d.Set("addon_json", flattenSpringCloudAppAddon(settings.AddonConfigs)); err != nil {
				return fmt.Errorf("setting `addon_json`: %s", err)
			}
			apmIds, err := flattenSpringCloudDeploymentApms(settings.Apms)
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

func resourceSpringCloudBuildDeploymentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudDeploymentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

func expandSpringCloudBuildDeploymentResourceRequests(input []interface{}) *appplatform.ResourceRequests {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	cpuResult := "1"
	memResult := "1Gi"

	v := input[0].(map[string]interface{})
	if cpuNew := v["cpu"].(string); cpuNew != "" {
		cpuResult = cpuNew
	}

	if memoryNew := v["memory"].(string); memoryNew != "" {
		memResult = memoryNew
	}

	result := appplatform.ResourceRequests{
		CPU:    utils.String(cpuResult),
		Memory: utils.String(memResult),
	}

	return &result
}
