// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_service

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	configServersClient := meta.(*clients.Client).AppPlatform.ConfigServersClient
	monitoringSettingsClient := meta.(*clients.Client).AppPlatform.MonitoringSettingsClient
	serviceRegistryClient := meta.(*clients.Client).AppPlatform.ServiceRegistryClient
	agentPoolClient := meta.(*clients.Client).AppPlatform.BuildServiceAgentPoolClient
	buildServiceClient := meta.(*clients.Client).AppPlatform.BuildServiceClient
	containerRegistryClient := meta.(*clients.Client).AppPlatform.ContainerRegistryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := commonids.ParseSpringCloudServiceIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.ServiceName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Spring Cloud Service %q (Resource Group %q): %+v", id.ServiceName, id.ResourceGroupName, err)
	}

	monitoringSettings, err := monitoringSettingsClient.Get(ctx, id.ResourceGroupName, id.ServiceName)
	if err != nil {
		return fmt.Errorf("retrieving monitoring settings for %s: %+v", id, err)
	}

	serviceRegistryEnabled := true
	serviceRegistry, err := serviceRegistryClient.Get(ctx, id.ResourceGroupName, id.ServiceName, "default")
	if err != nil {
		if !response.WasNotFound(serviceRegistry.Response.Response) {
			return fmt.Errorf("retrieving service registry of %s: %+v", id, err)
		}
		serviceRegistryEnabled = false
	}
	if response.WasNotFound(serviceRegistry.Response.Response) {
		serviceRegistryEnabled = false
	}

	containerRegistryList, err := containerRegistryClient.ListComplete(ctx, id.ResourceGroupName, id.ServiceName)
	if err == nil {
		containerRegistries := make([]appplatform.ContainerRegistryResource, 0)
		for containerRegistryList.NotDone() {
			containerRegistries = append(containerRegistries, containerRegistryList.Value())
			if err := containerRegistryList.NextWithContext(ctx); err != nil {
				return fmt.Errorf("going to next container registry value of %s: %+v", id, err)
			}
		}
		containerRegistriesState := expandSpringCloudContainerRegistries(d.Get("container_registry").([]any))
		d.Set("container_registry", flattenSpringCloudContainerRegistries(containerRegistriesState, containerRegistries))
	} else {
		log.Printf("[WARN] unable to list container registries for %s: %+v", id, err)
	}

	buildService, err := buildServiceClient.GetBuildService(ctx, id.ResourceGroupName, id.ServiceName, "default")
	if err == nil {
		d.Set("default_build_service", flattenSpringCloudBuildService(buildService.Properties))
	} else {
		log.Printf("[WARN] unable to get build service for %s: %+v", id, err)
	}

	agentPool, err := agentPoolClient.Get(ctx, id.ResourceGroupName, id.ServiceName, "default", "default")
	if err == nil && agentPool.Properties != nil && agentPool.Properties.PoolSize != nil {
		d.Set("build_agent_pool_size", agentPool.Properties.PoolSize.Name)
	} else {
		if err != nil {
			log.Printf("[WARN] error retrieving build agent pool of %q: %+v", id, err)
		}
		d.Set("build_agent_pool_size", "")
	}

	d.Set("name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
		d.Set("sku_tier", resp.Sku.Tier)
	}

	d.Set("service_registry_enabled", serviceRegistryEnabled)
	if serviceRegistryEnabled {
		d.Set("service_registry_id", appplatform_rm.NewServiceRegistryID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, "default").ID())
	} else {
		d.Set("service_registry_id", "")
	}

	if resp.Sku != nil && resp.Sku.Name != nil && *resp.Sku.Name != "E0" {
		configServer, err := configServersClient.Get(ctx, id.ResourceGroupName, id.ServiceName)
		if err != nil {
			return fmt.Errorf("retrieving config server configuration for %s: %+v", id, err)
		}
		if err := d.Set("config_server_git_setting", flattenSpringCloudConfigServerGitProperty(configServer.Properties, d)); err != nil {
			return fmt.Errorf("setting `config_server_git_setting`: %+v", err)
		}
	}

	if err := d.Set("trace", flattenSpringCloudTrace(monitoringSettings.Properties)); err != nil {
		return fmt.Errorf("failure setting `trace`: %+v", err)
	}

	if props := resp.Properties; props != nil {
		if err := d.Set("network", flattenSpringCloudNetwork(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `network`: %+v", err)
		}

		if err := d.Set("outbound_public_ip_addresses", flattenOutboundPublicIPAddresses(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `outbound_public_ip_addresses`: %+v", err)
		}

		if err := d.Set("required_network_traffic_rules", flattenRequiredTraffic(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `required_network_traffic_rules`: %+v", err)
		}

		if err := d.Set("marketplace", flattenSpringCloudMarketplaceResource(props.MarketplaceResource)); err != nil {
			return fmt.Errorf("setting `marketplace`: %+v", err)
		}

		if vnetAddons := props.VnetAddons; vnetAddons != nil {
			if err := d.Set("log_stream_public_endpoint_enabled", pointer.To(*vnetAddons.LogStreamPublicEndpoint)); err != nil {
				return fmt.Errorf("setting `log_stream_public_endpoint_enabled`: %+v", err)
			}
		}

		if managedEnvironmentID := props.ManagedEnvironmentID; managedEnvironmentID != nil {
			if err := d.Set("managed_environment_id", pointer.To(*props.ManagedEnvironmentID)); err != nil {
				return fmt.Errorf("setting `managed_environment_id`: %+v", err)
			}
		}

		d.Set("zone_redundant", props.ZoneRedundant)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func updateConfigServerSettings(ctx context.Context, client *appplatform.ConfigServersClient, id commonids.SpringCloudServiceId, gitProperty *appplatform.ConfigServerGitProperty) error {
	configServer := appplatform.ConfigServerResource{
		Properties: &appplatform.ConfigServerProperties{
			ConfigServer: &appplatform.ConfigServerSettings{
				GitProperty: gitProperty,
			},
		},
	}
	updateFuture, err := client.UpdatePut(ctx, id.ResourceGroupName, id.ServiceName, configServer)
	if err != nil {
		return fmt.Errorf("updating config server for %s: %+v", id, err)
	}
	if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of config server for %s: %+v", id, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.ServiceName)
	if err != nil {
		return fmt.Errorf("retrieving config server for %s: %+v", id, err)
	}
	if resp.Properties != nil && resp.Properties.Error != nil {
		if err := resp.Properties.Error; err != nil {
			return fmt.Errorf("setting config server for %s: %+v", id, err)
		}
	}
	return nil
}

func applyContainerRegistries(ctx context.Context, client *appplatform.ContainerRegistriesClient, springId commonids.SpringCloudServiceId, old []appplatform.ContainerRegistryResource, new []appplatform.ContainerRegistryResource) error {
	containerRegistriesToRemove := make(map[string]bool)
	for _, oldRegistry := range old {
		containerRegistriesToRemove[*oldRegistry.Name] = true
	}
	for _, newRegistry := range new {
		containerRegistriesToRemove[*newRegistry.Name] = false
	}
	for name := range containerRegistriesToRemove {
		if containerRegistriesToRemove[name] {
			id := appplatform_rm.NewContainerRegistryID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, name)
			future, err := client.Delete(ctx, springId.ResourceGroupName, springId.ServiceName, name)
			if err != nil {
				return fmt.Errorf("removing %s: %+v", id, err)
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", id, err)
			}
		}
	}
	for _, newRegistry := range new {
		id := appplatform_rm.NewContainerRegistryID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, *newRegistry.Name)
		future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, *newRegistry.Name, newRegistry)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation of %s: %+v", id, err)
		}
	}
	return nil
}
