// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_service

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).ServicesClient
	configServersClient := springcloudClient.From(meta).ConfigServersClient
	monitoringSettingsClient := springcloudClient.From(meta).MonitoringSettingsClient
	serviceRegistryClient := springcloudClient.From(meta).ServiceRegistryClient
	agentPoolClient := springcloudClient.From(meta).BuildServiceAgentPoolClient
	buildServiceClient := springcloudClient.From(meta).BuildServiceClient
	containerRegistryClient := springcloudClient.From(meta).ContainerRegistryClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := commonids.NewSpringCloudServiceID(subscriptionId, resourceGroup, name)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.ServiceName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_service", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	resource := appplatform.ServiceResource{
		Location: pointer.To(location),
		Properties: &appplatform.ClusterResourceProperties{
			NetworkProfile:      expandSpringCloudNetwork(d.Get("network").([]any)),
			ZoneRedundant:       pointer.To(d.Get("zone_redundant").(bool)),
			MarketplaceResource: expandSpringCloudMarketplaceResource(d.Get("marketplace").([]any)),
		},
		Sku: &appplatform.Sku{
			Name: pointer.To(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if enabled := d.Get("log_stream_public_endpoint_enabled").(bool); enabled {
		resource.Properties.VnetAddons = &appplatform.ServiceVNetAddons{
			LogStreamPublicEndpoint: pointer.To(enabled),
		}
	}

	// we set sku_tier only when managed_environment_id is set
	// otherwise we break the existing flows where sku_name is set to E0/Basic etc.,
	if v, ok := d.GetOk("managed_environment_id"); ok {
		resource.Properties.ManagedEnvironmentID = pointer.To(v.(string))
		resource.Sku.Tier = pointer.To(d.Get("sku_tier").(string))
	}

	gitProperty, err := expandSpringCloudConfigServerGitProperty(d.Get("config_server_git_setting").([]any))
	if err != nil {
		return err
	}

	// current create api doesn't take care parameters of config server.
	// so we need to invoke create api first and then update api
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.ServiceName, resource)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	skuName := d.Get("sku_name").(string)
	if skuName == "E0" && gitProperty != nil {
		return errors.New("`config_server_git_setting` is not supported for sku `E0`")
	}

	if skuName != "E0" {
		if err := updateConfigServerSettings(ctx, configServersClient, id, gitProperty); err != nil {
			return err
		}
	}

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

	if d.Get("service_registry_enabled").(bool) {
		future, err := serviceRegistryClient.CreateOrUpdate(ctx, id.ResourceGroupName, id.ServiceName, "default")
		if err != nil {
			return fmt.Errorf("creating service registry %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation service registry of %s: %+v", id, err)
		}
	}

	if skuName == "E0" {
		new := expandSpringCloudContainerRegistries(d.Get("container_registry").([]any))
		if err = applyContainerRegistries(ctx, containerRegistryClient, id, nil, new); err != nil {
			return fmt.Errorf("applying container registries for %s: %+v", id, err)
		}
		buildResource := appplatform.BuildService{
			Properties: pointer.To(expandSpringCloudBuildService(d.Get("default_build_service").([]any), id)),
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
