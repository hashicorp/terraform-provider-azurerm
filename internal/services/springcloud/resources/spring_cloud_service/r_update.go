// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_service

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).ServicesClient
	configServersClient := springcloudClient.From(meta).ConfigServersClient
	monitoringSettingsClient := springcloudClient.From(meta).MonitoringSettingsClient
	serviceRegistryClient := springcloudClient.From(meta).ServiceRegistryClient
	agentPoolClient := springcloudClient.From(meta).BuildServiceAgentPoolClient
	buildServiceClient := springcloudClient.From(meta).BuildServiceClient
	containerRegistryClient := springcloudClient.From(meta).ContainerRegistryClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := commonids.ParseSpringCloudServiceIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		model := appplatform.ServiceResource{
			Sku: &appplatform.Sku{
				Name: pointer.To(d.Get("sku_name").(string)),
			},
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}

		future, err := client.Update(ctx, id.ResourceGroupName, id.ServiceName, model)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	if d.HasChange("config_server_git_setting") {
		gitPropertyRaw := d.Get("config_server_git_setting").([]any)
		gitProperty, err := expandSpringCloudConfigServerGitProperty(gitPropertyRaw)
		if err != nil {
			return err
		}
		skuName := d.Get("sku_name").(string)
		if skuName == "E0" && gitProperty != nil {
			return errors.New("`config_server_git_setting` is not supported for sku `E0`")
		}
		if skuName != "E0" {
			if err := updateConfigServerSettings(ctx, configServersClient, *id, gitProperty); err != nil {
				return err
			}
		}
	}

	if d.HasChange("trace") {
		monitorSettings := appplatform.MonitoringSettingResource{
			Properties: expandSpringCloudTrace(d.Get("trace").([]any)),
		}
		updateFuture, err := monitoringSettingsClient.UpdatePut(ctx, id.ResourceGroupName, id.ServiceName, monitorSettings)
		if err != nil {
			return fmt.Errorf("updating monitor settings for %s: %+v", id, err)
		}
		if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of monitor settings for %s: %+v", id, err)
		}
	}

	if d.HasChange("service_registry_enabled") {
		if d.Get("service_registry_enabled").(bool) {
			future, err := serviceRegistryClient.CreateOrUpdate(ctx, id.ResourceGroupName, id.ServiceName, "default")
			if err != nil {
				return fmt.Errorf("creating service registry of %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation service registry of %s: %+v", id, err)
			}
		} else {
			future, err := serviceRegistryClient.Delete(ctx, id.ResourceGroupName, id.ServiceName, "default")
			if err != nil {
				return fmt.Errorf("deleting service registry of %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion service registry of %s: %+v", id, err)
			}
		}
	}

	skuName := d.Get("sku_name").(string)
	if skuName == "E0" {
		new := expandSpringCloudContainerRegistries(d.Get("container_registry").([]any))
		if err = applyContainerRegistries(ctx, containerRegistryClient, *id, nil, new); err != nil {
			return fmt.Errorf("applying container registries for %s: %+v", id, err)
		}
		buildResource := appplatform.BuildService{
			Properties: pointer.To(expandSpringCloudBuildService(d.Get("default_build_service").([]any), *id)),
		}
		buildServiceCreateFuture, err := buildServiceClient.CreateOrUpdate(ctx, id.ResourceGroupName, id.ServiceName, "default", buildResource)
		if err != nil {
			return fmt.Errorf("creating build service %s: %+v", id, err)
		}
		if err := buildServiceCreateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation build service of %s: %+v", id, err)
		}
	}

	if size := d.Get("build_agent_pool_size").(string); len(size) > 0 {
		agentPoolResource := appplatform.BuildServiceAgentPoolResource{
			Properties: &appplatform.BuildServiceAgentPoolProperties{
				PoolSize: &appplatform.BuildServiceAgentPoolSizeProperties{
					Name: pointer.To(size),
				},
			},
		}
		future, err := agentPoolClient.UpdatePut(ctx, id.ResourceGroupName, id.ServiceName, "default", "default", agentPoolResource)
		if err != nil {
			return fmt.Errorf("creating default build agent of %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation default build agent of %s: %+v", id, err)
		}
	}

	return resourceSpringCloudServiceRead(d, meta)
}
