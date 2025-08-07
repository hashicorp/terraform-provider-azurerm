// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_service` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudServiceCreate,
		Read:   resourceSpringCloudServiceRead,
		Update: resourceSpringCloudServiceUpdate,
		Delete: resourceSpringCloudServiceDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			// Spring Cloud Service only supports following locations, we are still supporting more locations (Wednesday, November 20, 2019 4:20 PM):
			// `East US`, `Southeast Asia`, `West Europe`, `West US 2`
			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "S0",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"B0",
					"S0",
					"E0",
				}, false),
			},

			"sku_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Basic",
					"Enterprise",
					"Standard",
					"StandardGen2",
				}, false),
			},

			"managed_environment_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"build_agent_pool_size": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"S1",
					"S2",
					"S3",
					"S4",
					"S5",
				}, false),
			},

			"container_registry": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"password": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"server": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"username": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"default_build_service": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"container_registry_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"log_stream_public_endpoint_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"marketplace": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"plan": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"publisher": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"product": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"network": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"app_subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},

						"service_runtime_subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},

						"cidr_ranges": {
							Type:     pluginsdk.TypeList,
							Required: true,
							ForceNew: true,
							MinItems: 3,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"app_network_resource_group": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"outbound_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "loadBalancer",
							ValidateFunc: validation.StringInSlice([]string{
								"loadBalancer",
								"userDefinedRouting",
							}, false),
						},

						"read_timeout_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"service_runtime_network_resource_group": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},

			"config_server_git_setting": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"uri": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ConfigServerURI,
						},

						"label": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"search_paths": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"http_basic_auth": SchemaConfigServerHttpBasicAuth("config_server_git_setting.0.ssh_auth"),

						"ssh_auth": SchemaConfigServerSSHAuth("config_server_git_setting.0.http_basic_auth"),

						"repository": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"uri": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.ConfigServerURI,
									},

									"label": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},

									"pattern": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"search_paths": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"http_basic_auth": SchemaConfigServerHttpBasicAuth(),

									"ssh_auth": SchemaConfigServerSSHAuth(),
								},
							},
						},
					},
				},
			},

			"trace": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"connection_string": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"sample_rate": {
							Type:         pluginsdk.TypeFloat,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.FloatBetween(0, 100),
						},
					},
				},
			},

			"service_registry_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"outbound_public_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"required_network_traffic_rules": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"fqdns": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"direction": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tags.Schema(),

			"service_registry_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSpringCloudServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	configServersClient := meta.(*clients.Client).AppPlatform.ConfigServersClient
	monitoringSettingsClient := meta.(*clients.Client).AppPlatform.MonitoringSettingsClient
	serviceRegistryClient := meta.(*clients.Client).AppPlatform.ServiceRegistryClient
	agentPoolClient := meta.(*clients.Client).AppPlatform.BuildServiceAgentPoolClient
	buildServiceClient := meta.(*clients.Client).AppPlatform.BuildServiceClient
	containerRegistryClient := meta.(*clients.Client).AppPlatform.ContainerRegistryClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewSpringCloudServiceID(subscriptionId, resourceGroup, name)
	existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_spring_cloud_service", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	resource := appplatform.ServiceResource{
		Location: utils.String(location),
		Properties: &appplatform.ClusterResourceProperties{
			NetworkProfile:      expandSpringCloudNetwork(d.Get("network").([]interface{})),
			ZoneRedundant:       utils.Bool(d.Get("zone_redundant").(bool)),
			MarketplaceResource: expandSpringCloudMarketplaceResource(d.Get("marketplace").([]interface{})),
		},
		Sku: &appplatform.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if enabled := d.Get("log_stream_public_endpoint_enabled").(bool); enabled {
		resource.Properties.VnetAddons = &appplatform.ServiceVNetAddons{
			LogStreamPublicEndpoint: utils.Bool(enabled),
		}
	}

	// we set sku_tier only when managed_environment_id is set
	// otherwise we break the existing flows where sku_name is set to E0/Basic etc.,
	if v, ok := d.GetOk("managed_environment_id"); ok {
		resource.Properties.ManagedEnvironmentID = utils.String(v.(string))
		resource.Sku.Tier = utils.String(d.Get("sku_tier").(string))
	}

	gitProperty, err := expandSpringCloudConfigServerGitProperty(d.Get("config_server_git_setting").([]interface{}))
	if err != nil {
		return err
	}

	// current create api doesn't take care parameters of config server.
	// so we need to invoke create api first and then update api
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, resource)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}
	d.SetId(id.ID())

	skuName := d.Get("sku_name").(string)
	if skuName == "E0" && gitProperty != nil {
		return fmt.Errorf("`config_server_git_setting` is not supported for sku `E0`")
	}

	if skuName != "E0" {
		log.Printf("[DEBUG] Updating Config Server Settings for %s..", id)
		if err := updateConfigServerSettings(ctx, configServersClient, id, gitProperty); err != nil {
			return err
		}
		log.Printf("[DEBUG] Updated Config Server Settings for %s.", id)
	}

	log.Printf("[DEBUG] Updating Monitor Settings for %s..", id)
	monitorSettings := appplatform.MonitoringSettingResource{
		Properties: expandSpringCloudTrace(d.Get("trace").([]interface{})),
	}
	updateFuture, err := monitoringSettingsClient.UpdatePut(ctx, id.ResourceGroup, id.SpringName, monitorSettings)
	if err != nil {
		return fmt.Errorf("updating monitor settings for %s: %+v", id, err)
	}
	if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of monitor settings for %s: %+v", id, err)
	}
	log.Printf("[DEBUG] Updated Monitor Settings for %s.", id)

	if d.Get("service_registry_enabled").(bool) {
		future, err := serviceRegistryClient.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, "default")
		if err != nil {
			return fmt.Errorf("creating service registry %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation service registry of %s: %+v", id, err)
		}
	}

	if skuName == "E0" {
		new := expandSpringCloudContainerRegistries(d.Get("container_registry").([]interface{}))
		if err = applyContainerRegistries(ctx, containerRegistryClient, id, nil, new); err != nil {
			return fmt.Errorf("applying container registries for %s: %+v", id, err)
		}
		buildResource := appplatform.BuildService{
			Properties: pointer.To(expandSpringCloudBuildService(d.Get("default_build_service").([]interface{}), id)),
		}
		buildServiceCreateFuture, err := buildServiceClient.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, "default", buildResource)
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
					Name: utils.String(size),
				},
			},
		}
		future, err := agentPoolClient.UpdatePut(ctx, id.ResourceGroup, id.SpringName, "default", "default", agentPoolResource)
		if err != nil {
			return fmt.Errorf("creating default build agent of %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation default build agent of %s: %+v", id, err)
		}
	}

	return resourceSpringCloudServiceRead(d, meta)
}

func resourceSpringCloudServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	configServersClient := meta.(*clients.Client).AppPlatform.ConfigServersClient
	monitoringSettingsClient := meta.(*clients.Client).AppPlatform.MonitoringSettingsClient
	serviceRegistryClient := meta.(*clients.Client).AppPlatform.ServiceRegistryClient
	agentPoolClient := meta.(*clients.Client).AppPlatform.BuildServiceAgentPoolClient
	buildServiceClient := meta.(*clients.Client).AppPlatform.BuildServiceClient
	containerRegistryClient := meta.(*clients.Client).AppPlatform.ContainerRegistryClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		model := appplatform.ServiceResource{
			Sku: &appplatform.Sku{
				Name: utils.String(d.Get("sku_name").(string)),
			},
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}

		future, err := client.Update(ctx, id.ResourceGroup, id.SpringName, model)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	if d.HasChange("config_server_git_setting") {
		gitPropertyRaw := d.Get("config_server_git_setting").([]interface{})
		gitProperty, err := expandSpringCloudConfigServerGitProperty(gitPropertyRaw)
		if err != nil {
			return err
		}
		skuName := d.Get("sku_name").(string)
		if skuName == "E0" && gitProperty != nil {
			return fmt.Errorf("`config_server_git_setting` is not supported for sku `E0`")
		}
		if skuName != "E0" {
			log.Printf("[DEBUG] Updating Config Server Settings for %s..", *id)
			if err := updateConfigServerSettings(ctx, configServersClient, *id, gitProperty); err != nil {
				return err
			}
			log.Printf("[DEBUG] Updated Config Server Settings for %s.", *id)
		}
	}

	if d.HasChange("trace") {
		log.Printf("[DEBUG] Updating Monitor Settings for %s..", id)
		monitorSettings := appplatform.MonitoringSettingResource{
			Properties: expandSpringCloudTrace(d.Get("trace").([]interface{})),
		}
		updateFuture, err := monitoringSettingsClient.UpdatePut(ctx, id.ResourceGroup, id.SpringName, monitorSettings)
		if err != nil {
			return fmt.Errorf("updating monitor settings for %s: %+v", id, err)
		}
		if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of monitor settings for %s: %+v", id, err)
		}
		log.Printf("[DEBUG] Updated Monitor Settings for %s.", id)
	}

	if d.HasChange("service_registry_enabled") {
		if d.Get("service_registry_enabled").(bool) {
			future, err := serviceRegistryClient.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, "default")
			if err != nil {
				return fmt.Errorf("creating service registry of %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation service registry of %s: %+v", id, err)
			}
		} else {
			future, err := serviceRegistryClient.Delete(ctx, id.ResourceGroup, id.SpringName, "default")
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
		new := expandSpringCloudContainerRegistries(d.Get("container_registry").([]interface{}))
		if err = applyContainerRegistries(ctx, containerRegistryClient, *id, nil, new); err != nil {
			return fmt.Errorf("applying container registries for %s: %+v", id, err)
		}
		buildResource := appplatform.BuildService{
			Properties: pointer.To(expandSpringCloudBuildService(d.Get("default_build_service").([]interface{}), *id)),
		}
		buildServiceCreateFuture, err := buildServiceClient.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, "default", buildResource)
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
					Name: utils.String(size),
				},
			},
		}
		future, err := agentPoolClient.UpdatePut(ctx, id.ResourceGroup, id.SpringName, "default", "default", agentPoolResource)
		if err != nil {
			return fmt.Errorf("creating default build agent of %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation default build agent of %s: %+v", id, err)
		}
	}

	return resourceSpringCloudServiceRead(d, meta)
}

func resourceSpringCloudServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	configServersClient := meta.(*clients.Client).AppPlatform.ConfigServersClient
	monitoringSettingsClient := meta.(*clients.Client).AppPlatform.MonitoringSettingsClient
	serviceRegistryClient := meta.(*clients.Client).AppPlatform.ServiceRegistryClient
	agentPoolClient := meta.(*clients.Client).AppPlatform.BuildServiceAgentPoolClient
	buildServiceClient := meta.(*clients.Client).AppPlatform.BuildServiceClient
	containerRegistryClient := meta.(*clients.Client).AppPlatform.ContainerRegistryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Spring Cloud Service %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	monitoringSettings, err := monitoringSettingsClient.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return fmt.Errorf("retrieving monitoring settings for %s: %+v", id, err)
	}

	serviceRegistryEnabled := true
	serviceRegistry, err := serviceRegistryClient.Get(ctx, id.ResourceGroup, id.SpringName, "default")
	if err != nil {
		if !utils.ResponseWasNotFound(serviceRegistry.Response) {
			return fmt.Errorf("retrieving service registry of %s: %+v", id, err)
		}
		serviceRegistryEnabled = false
	}
	if utils.ResponseWasNotFound(serviceRegistry.Response) {
		serviceRegistryEnabled = false
	}

	containerRegistryList, err := containerRegistryClient.ListComplete(ctx, id.ResourceGroup, id.SpringName)
	if err == nil {
		containerRegistries := make([]appplatform.ContainerRegistryResource, 0)
		for containerRegistryList.NotDone() {
			containerRegistries = append(containerRegistries, containerRegistryList.Value())
			if err := containerRegistryList.NextWithContext(ctx); err != nil {
				return fmt.Errorf("going to next container registry value of %s: %+v", id, err)
			}
		}
		containerRegistriesState := expandSpringCloudContainerRegistries(d.Get("container_registry").([]interface{}))
		d.Set("container_registry", flattenSpringCloudContainerRegistries(containerRegistriesState, containerRegistries))
	} else {
		log.Printf("[WARN] unable to list container registries for %s: %+v", id, err)
	}

	buildService, err := buildServiceClient.GetBuildService(ctx, id.ResourceGroup, id.SpringName, "default")
	if err == nil {
		d.Set("default_build_service", flattenSpringCloudBuildService(buildService.Properties))
	} else {
		log.Printf("[WARN] unable to get build service for %s: %+v", id, err)
	}

	agentPool, err := agentPoolClient.Get(ctx, id.ResourceGroup, id.SpringName, "default", "default")
	if err == nil && agentPool.Properties != nil && agentPool.Properties.PoolSize != nil {
		d.Set("build_agent_pool_size", agentPool.Properties.PoolSize.Name)
	} else {
		if err != nil {
			log.Printf("[WARN] error retrieving build agent pool of %q: %+v", id, err)
		}
		d.Set("build_agent_pool_size", "")
	}

	d.Set("name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
		d.Set("sku_tier", resp.Sku.Tier)
	}

	d.Set("service_registry_enabled", serviceRegistryEnabled)
	if serviceRegistryEnabled {
		d.Set("service_registry_id", parse.NewSpringCloudServiceRegistryID(id.SubscriptionId, id.ResourceGroup, id.SpringName, "default").ID())
	} else {
		d.Set("service_registry_id", "")
	}

	if resp.Sku != nil && resp.Sku.Name != nil && *resp.Sku.Name != "E0" {
		configServer, err := configServersClient.Get(ctx, id.ResourceGroup, id.SpringName)
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

		outboundPublicIPAddresses := flattenOutboundPublicIPAddresses(props.NetworkProfile)
		if err := d.Set("outbound_public_ip_addresses", outboundPublicIPAddresses); err != nil {
			return fmt.Errorf("setting `outbound_public_ip_addresses`: %+v", err)
		}

		if err := d.Set("required_network_traffic_rules", flattenRequiredTraffic(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `required_network_traffic_rules`: %+v", err)
		}

		if err := d.Set("marketplace", flattenSpringCloudMarketplaceResource(props.MarketplaceResource)); err != nil {
			return fmt.Errorf("setting `marketplace`: %+v", err)
		}

		if vnetAddons := props.VnetAddons; vnetAddons != nil {
			if err := d.Set("log_stream_public_endpoint_enabled", utils.Bool(*vnetAddons.LogStreamPublicEndpoint)); err != nil {
				return fmt.Errorf("setting `log_stream_public_endpoint_enabled`: %+v", err)
			}
		}

		if managedEnvironmentID := props.ManagedEnvironmentID; managedEnvironmentID != nil {
			if err := d.Set("managed_environment_id", utils.String(*props.ManagedEnvironmentID)); err != nil {
				return fmt.Errorf("setting `managed_environment_id`: %+v", err)
			}
		}

		d.Set("zone_redundant", props.ZoneRedundant)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSpringCloudServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
		}
	}

	return nil
}

func updateConfigServerSettings(ctx context.Context, client *appplatform.ConfigServersClient, id parse.SpringCloudServiceId, gitProperty *appplatform.ConfigServerGitProperty) error {
	log.Printf("[DEBUG] Updating Config Server Settings for %s..", id)
	configServer := appplatform.ConfigServerResource{
		Properties: &appplatform.ConfigServerProperties{
			ConfigServer: &appplatform.ConfigServerSettings{
				GitProperty: gitProperty,
			},
		},
	}
	updateFuture, err := client.UpdatePut(ctx, id.ResourceGroup, id.SpringName, configServer)
	if err != nil {
		return fmt.Errorf("updating config server for %s: %+v", id, err)
	}
	if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of config server for %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Retrieving Config Server Settings for %s..", id)
	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return fmt.Errorf("retrieving config server for %s: %+v", id, err)
	}
	if resp.Properties != nil && resp.Properties.Error != nil {
		if err := resp.Properties.Error; err != nil {
			return fmt.Errorf("setting config server for %s: %+v", id, err)
		}
	}
	log.Printf("[DEBUG] Updated Config Server Settings for %s.", id)
	return nil
}

func applyContainerRegistries(ctx context.Context, client *appplatform.ContainerRegistriesClient, springId parse.SpringCloudServiceId, old []appplatform.ContainerRegistryResource, new []appplatform.ContainerRegistryResource) error {
	containerRegistriesToRemove := make(map[string]bool)
	for _, oldRegistry := range old {
		containerRegistriesToRemove[*oldRegistry.Name] = true
	}
	for _, newRegistry := range new {
		containerRegistriesToRemove[*newRegistry.Name] = false
	}
	for name := range containerRegistriesToRemove {
		if containerRegistriesToRemove[name] {
			id := parse.NewSpringCloudContainerRegistryID(springId.SubscriptionId, springId.ResourceGroup, springId.SpringName, name)
			future, err := client.Delete(ctx, springId.ResourceGroup, springId.SpringName, name)
			if err != nil {
				return fmt.Errorf("removing %s: %+v", id, err)
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", id, err)
			}
		}
	}
	for _, newRegistry := range new {
		id := parse.NewSpringCloudContainerRegistryID(springId.SubscriptionId, springId.ResourceGroup, springId.SpringName, *newRegistry.Name)
		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, *newRegistry.Name, newRegistry)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation of %s: %+v", id, err)
		}
	}
	return nil
}

func expandSpringCloudNetwork(input []interface{}) *appplatform.NetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	cidrRanges := utils.ExpandStringSlice(v["cidr_ranges"].([]interface{}))
	network := &appplatform.NetworkProfile{
		ServiceRuntimeSubnetID: utils.String(v["service_runtime_subnet_id"].(string)),
		AppSubnetID:            utils.String(v["app_subnet_id"].(string)),
		ServiceCidr:            utils.String(strings.Join(*cidrRanges, ",")),
		OutboundType:           utils.String(v["outbound_type"].(string)),
	}
	if readTimeoutInSeconds := v["read_timeout_seconds"].(int); readTimeoutInSeconds != 0 {
		network.IngressConfig = &appplatform.IngressConfig{
			ReadTimeoutInSeconds: utils.Int32(int32(readTimeoutInSeconds)),
		}
	}
	if serviceRuntimeNetworkResourceGroup := v["service_runtime_network_resource_group"].(string); serviceRuntimeNetworkResourceGroup != "" {
		network.ServiceRuntimeNetworkResourceGroup = utils.String(serviceRuntimeNetworkResourceGroup)
	}
	if appNetworkResourceGroup := v["app_network_resource_group"].(string); appNetworkResourceGroup != "" {
		network.AppNetworkResourceGroup = utils.String(appNetworkResourceGroup)
	}
	return network
}

func expandSpringCloudConfigServerGitProperty(input []interface{}) (*appplatform.ConfigServerGitProperty, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	if v == nil {
		return nil, nil
	}

	result := appplatform.ConfigServerGitProperty{
		URI: utils.String(v["uri"].(string)),
	}

	if label := v["label"].(string); label != "" {
		result.Label = utils.String(label)
	}
	if searchPaths := v["search_paths"].([]interface{}); len(searchPaths) > 0 {
		result.SearchPaths = utils.ExpandStringSlice(searchPaths)
	}

	httpBasicAuth := v["http_basic_auth"].([]interface{})
	sshAuth := v["ssh_auth"].([]interface{})
	if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
		return nil, fmt.Errorf("can not set both `http_basic_auth` and `ssh_auth`")
	}
	if len(httpBasicAuth) > 0 {
		v := httpBasicAuth[0].(map[string]interface{})
		result.Username = utils.String(v["username"].(string))
		result.Password = utils.String(v["password"].(string))
	}
	if len(sshAuth) > 0 {
		v := sshAuth[0].(map[string]interface{})
		result.PrivateKey = utils.String(v["private_key"].(string))
		result.StrictHostKeyChecking = utils.Bool(v["strict_host_key_checking_enabled"].(bool))

		if hostKey := v["host_key"].(string); hostKey != "" {
			result.HostKey = utils.String(hostKey)
		}
		if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
			result.HostKeyAlgorithm = utils.String(hostKeyAlgorithm)
		}
	}

	if v, ok := v["repository"]; ok {
		repositories, err := expandSpringCloudGitPatternRepository(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		result.Repositories = repositories
	}

	return &result, nil
}

func expandSpringCloudGitPatternRepository(input []interface{}) (*[]appplatform.GitPatternRepository, error) {
	results := make([]appplatform.GitPatternRepository, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		result := appplatform.GitPatternRepository{
			Name: utils.String(v["name"].(string)),
			URI:  utils.String(v["uri"].(string)),
		}

		if label := v["label"].(string); len(label) > 0 {
			result.Label = utils.String(label)
		}
		if pattern := v["pattern"].([]interface{}); len(pattern) > 0 {
			result.Pattern = utils.ExpandStringSlice(pattern)
		}
		if searchPaths := v["search_paths"].([]interface{}); len(searchPaths) > 0 {
			result.SearchPaths = utils.ExpandStringSlice(searchPaths)
		}

		httpBasicAuth := v["http_basic_auth"].([]interface{})
		sshAuth := v["ssh_auth"].([]interface{})
		if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
			return nil, fmt.Errorf("can not set both `http_basic_auth` and `ssh_auth` for the same repository")
		}
		if len(httpBasicAuth) > 0 {
			v := httpBasicAuth[0].(map[string]interface{})
			result.Username = utils.String(v["username"].(string))
			result.Password = utils.String(v["password"].(string))
		}
		if len(sshAuth) > 0 {
			v := sshAuth[0].(map[string]interface{})
			result.PrivateKey = utils.String(v["private_key"].(string))
			result.StrictHostKeyChecking = utils.Bool(v["strict_host_key_checking_enabled"].(bool))

			if hostKey := v["host_key"].(string); hostKey != "" {
				result.HostKey = utils.String(hostKey)
			}
			if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
				result.HostKeyAlgorithm = utils.String(hostKeyAlgorithm)
			}
		}

		results = append(results, result)
	}
	return &results, nil
}

func expandSpringCloudTrace(input []interface{}) *appplatform.MonitoringSettingProperties {
	if len(input) == 0 || input[0] == nil {
		return &appplatform.MonitoringSettingProperties{
			TraceEnabled: utils.Bool(false),
		}
	}

	v := input[0].(map[string]interface{})
	return &appplatform.MonitoringSettingProperties{
		TraceEnabled:                  utils.Bool(true),
		AppInsightsInstrumentationKey: utils.String(v["connection_string"].(string)),
		AppInsightsSamplingRate:       utils.Float(v["sample_rate"].(float64)),
	}
}

func expandSpringCloudContainerRegistries(input []interface{}) []appplatform.ContainerRegistryResource {
	if len(input) == 0 {
		return nil
	}
	out := make([]appplatform.ContainerRegistryResource, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		out = append(out, appplatform.ContainerRegistryResource{
			Name: utils.String(v["name"].(string)),
			Properties: &appplatform.ContainerRegistryProperties{
				Credentials: &appplatform.ContainerRegistryBasicCredentials{
					Username: utils.String(v["username"].(string)),
					Password: utils.String(v["password"].(string)),
					Server:   utils.String(v["server"].(string)),
				},
			},
		})
	}
	return out
}

func expandSpringCloudBuildService(input []interface{}, springId parse.SpringCloudServiceId) appplatform.BuildServiceProperties {
	if len(input) == 0 || input[0] == nil {
		return appplatform.BuildServiceProperties{}
	}
	v := input[0].(map[string]interface{})
	out := appplatform.BuildServiceProperties{}
	if value := v["container_registry_name"].(string); value != "" {
		out.ContainerRegistry = utils.String(parse.NewSpringCloudContainerRegistryID(springId.SubscriptionId, springId.ResourceGroup, springId.SpringName, value).ID())
	}
	return out
}

func expandSpringCloudMarketplaceResource(input []interface{}) *appplatform.MarketplaceResource {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.MarketplaceResource{
		Plan:      utils.String(v["plan"].(string)),
		Publisher: utils.String(v["publisher"].(string)),
		Product:   utils.String(v["product"].(string)),
	}
}

func flattenSpringCloudConfigServerGitProperty(input *appplatform.ConfigServerProperties, d *pluginsdk.ResourceData) []interface{} {
	if input == nil || input.ConfigServer == nil || input.ConfigServer.GitProperty == nil {
		return []interface{}{}
	}

	gitProperty := input.ConfigServer.GitProperty

	// prepare old state to find sensitive props not returned by API.
	oldGitSetting := make(map[string]interface{})
	if oldGitSettings := d.Get("config_server_git_setting").([]interface{}); len(oldGitSettings) > 0 {
		oldGitSetting = oldGitSettings[0].(map[string]interface{})
	}

	uri := ""
	if gitProperty.URI != nil {
		uri = *gitProperty.URI
	}

	label := ""
	if gitProperty.Label != nil {
		label = *gitProperty.Label
	}

	searchPaths := utils.FlattenStringSlice(gitProperty.SearchPaths)

	httpBasicAuth := make([]interface{}, 0)
	if gitProperty.Username != nil && gitProperty.Password != nil {
		// username and password returned by API are *
		// to avoid state diff, we get the props from old state
		username := ""
		password := ""
		if v, ok := oldGitSetting["http_basic_auth"]; ok {
			oldHTTPBasicAuth := v.([]interface{})
			if len(oldHTTPBasicAuth) > 0 {
				oldItem := oldHTTPBasicAuth[0].(map[string]interface{})
				username = oldItem["username"].(string)
				password = oldItem["password"].(string)
			}
		}

		httpBasicAuth = []interface{}{
			map[string]interface{}{
				"username": username,
				"password": password,
			},
		}
	}

	sshAuth := []interface{}{}
	if gitProperty.PrivateKey != nil {
		// private_key, host_key and host_key_algorithm returned by API are *
		// to avoid state diff, we get the props from old state
		privateKey := ""
		hostKey := ""
		hostKeyAlgorithm := ""
		if v, ok := oldGitSetting["ssh_auth"]; ok {
			sshAuth := v.([]interface{})
			if len(sshAuth) > 0 {
				oldItem := sshAuth[0].(map[string]interface{})
				privateKey = oldItem["private_key"].(string)
				hostKey = oldItem["host_key"].(string)
				hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
			}
		}

		strictHostKeyChecking := false
		if gitProperty.StrictHostKeyChecking != nil {
			strictHostKeyChecking = *gitProperty.StrictHostKeyChecking
		}

		sshAuth = []interface{}{
			map[string]interface{}{
				"private_key":                      privateKey,
				"host_key":                         hostKey,
				"host_key_algorithm":               hostKeyAlgorithm,
				"strict_host_key_checking_enabled": strictHostKeyChecking,
			},
		}
	}

	return []interface{}{
		map[string]interface{}{
			"uri":             uri,
			"label":           label,
			"search_paths":    searchPaths,
			"http_basic_auth": httpBasicAuth,
			"ssh_auth":        sshAuth,
			"repository":      flattenSpringCloudGitPatternRepository(gitProperty.Repositories, d),
		},
	}
}

func flattenSpringCloudGitPatternRepository(input *[]appplatform.GitPatternRepository, d *pluginsdk.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// prepare old state to find sensitive props not returned by API.
	oldGitPatternRepositories := map[string]interface{}{}
	if oldGitSettings := d.Get("config_server_git_setting").([]interface{}); len(oldGitSettings) > 0 {
		oldGitSetting := oldGitSettings[0].(map[string]interface{})
		for _, r := range oldGitSetting["repository"].([]interface{}) {
			repo := r.(map[string]interface{})
			if name, ok := repo["name"]; ok {
				oldGitPatternRepositories[name.(string)] = r
			}
		}
	}

	for _, item := range *input {
		name := ""
		if item.Name != nil {
			name = *item.Name
		}

		uri := ""
		if item.URI != nil {
			uri = *item.URI
		}

		label := ""
		if item.Label != nil {
			label = *item.Label
		}

		// prepare old state to find sensitive props not returned by API.
		oldGitPatternRepository := make(map[string]interface{})
		if gpr, ok := oldGitPatternRepositories[name]; ok {
			oldGitPatternRepository = gpr.(map[string]interface{})
		}

		pattern := utils.FlattenStringSlice(item.Pattern)
		searchPaths := utils.FlattenStringSlice(item.SearchPaths)

		httpBasicAuth := []interface{}{}
		if item.Username != nil && item.Password != nil {
			// username and password returned by API are *
			// to avoid state diff, we get the props from old state
			username := ""
			password := ""
			if v, ok := oldGitPatternRepository["http_basic_auth"]; ok {
				oldHTTPBasicAuth := v.([]interface{})
				if len(oldHTTPBasicAuth) > 0 {
					oldItem := oldHTTPBasicAuth[0].(map[string]interface{})
					username = oldItem["username"].(string)
					password = oldItem["password"].(string)
				}
			}

			httpBasicAuth = []interface{}{
				map[string]interface{}{
					"username": username,
					"password": password,
				},
			}
		}

		sshAuth := []interface{}{}
		if item.PrivateKey != nil {
			// private_key, host_key and host_key_algorithm returned by API are *
			// to avoid state diff, we get the props from old state
			privateKey := ""
			hostKey := ""
			hostKeyAlgorithm := ""
			if v, ok := oldGitPatternRepository["ssh_auth"]; ok {
				sshAuth := v.([]interface{})
				if len(sshAuth) > 0 {
					oldItem := sshAuth[0].(map[string]interface{})
					privateKey = oldItem["private_key"].(string)
					hostKey = oldItem["host_key"].(string)
					hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
				}
			}

			strictHostKeyChecking := false
			if item.StrictHostKeyChecking != nil {
				strictHostKeyChecking = *item.StrictHostKeyChecking
			}

			sshAuth = []interface{}{
				map[string]interface{}{
					"private_key":                      privateKey,
					"host_key":                         hostKey,
					"host_key_algorithm":               hostKeyAlgorithm,
					"strict_host_key_checking_enabled": strictHostKeyChecking,
				},
			}
		}

		results = append(results, map[string]interface{}{
			"name":            name,
			"uri":             uri,
			"label":           label,
			"pattern":         pattern,
			"search_paths":    searchPaths,
			"http_basic_auth": httpBasicAuth,
			"ssh_auth":        sshAuth,
		})
	}

	return results
}

func flattenSpringCloudTrace(input *appplatform.MonitoringSettingProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enabled := false
	connectionString := ""
	samplingRate := 0.0
	if input.TraceEnabled != nil {
		enabled = *input.TraceEnabled
	}
	if input.AppInsightsInstrumentationKey != nil {
		connectionString = *input.AppInsightsInstrumentationKey
	}
	if input.AppInsightsSamplingRate != nil {
		samplingRate = *input.AppInsightsSamplingRate
	}

	if !enabled {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"connection_string": connectionString,
			"sample_rate":       samplingRate,
		},
	}
}

func flattenSpringCloudNetwork(input *appplatform.NetworkProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var serviceRuntimeSubnetID, appSubnetID, serviceRuntimeNetworkResourceGroup, appNetworkResourceGroup string
	var readTimeoutInSeconds int32
	var cidrRanges []interface{}
	if input.ServiceRuntimeSubnetID != nil {
		serviceRuntimeSubnetID = *input.ServiceRuntimeSubnetID
	}
	if input.AppSubnetID != nil {
		appSubnetID = *input.AppSubnetID
	}
	if input.ServiceCidr != nil {
		cidrs := strings.Split(*input.ServiceCidr, ",")
		cidrRanges = utils.FlattenStringSlice(&cidrs)
	}
	if input.ServiceRuntimeNetworkResourceGroup != nil {
		serviceRuntimeNetworkResourceGroup = *input.ServiceRuntimeNetworkResourceGroup
	}
	if input.AppNetworkResourceGroup != nil {
		appNetworkResourceGroup = *input.AppNetworkResourceGroup
	}

	if ingressConfig := input.IngressConfig; ingressConfig != nil {
		if ingressConfig.ReadTimeoutInSeconds != nil {
			readTimeoutInSeconds = *ingressConfig.ReadTimeoutInSeconds
		}
	}

	outboundType := "loadBalancer"
	if input.OutboundType != nil {
		outboundType = *input.OutboundType
	}

	if serviceRuntimeSubnetID == "" && appSubnetID == "" && serviceRuntimeNetworkResourceGroup == "" && appNetworkResourceGroup == "" && len(cidrRanges) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"app_subnet_id":                          appSubnetID,
			"service_runtime_subnet_id":              serviceRuntimeSubnetID,
			"cidr_ranges":                            cidrRanges,
			"app_network_resource_group":             appNetworkResourceGroup,
			"read_timeout_seconds":                   readTimeoutInSeconds,
			"service_runtime_network_resource_group": serviceRuntimeNetworkResourceGroup,
			"outbound_type":                          outboundType,
		},
	}
}

func flattenOutboundPublicIPAddresses(input *appplatform.NetworkProfile) []interface{} {
	if input == nil || input.OutboundIPs == nil {
		return []interface{}{}
	}

	return utils.FlattenStringSlice(input.OutboundIPs.PublicIPs)
}

func flattenRequiredTraffic(input *appplatform.NetworkProfile) []interface{} {
	if input == nil || input.RequiredTraffics == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input.RequiredTraffics {
		protocol := ""
		if v.Protocol != nil {
			protocol = *v.Protocol
		}

		port := 0
		if v.Port != nil {
			port = int(*v.Port)
		}

		result = append(result, map[string]interface{}{
			"protocol":     protocol,
			"port":         port,
			"ip_addresses": utils.FlattenStringSlice(v.Ips),
			"fqdns":        utils.FlattenStringSlice(v.Fqdns),
			"direction":    string(v.Direction),
		})
	}
	return result
}

func flattenSpringCloudContainerRegistries(state []appplatform.ContainerRegistryResource, input []appplatform.ContainerRegistryResource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	statePasswordMap := make(map[string]string)
	for _, v := range state {
		if v.Name == nil || v.Properties == nil || v.Properties.Credentials == nil {
			continue
		}
		if basicCredential, ok := v.Properties.Credentials.AsContainerRegistryBasicCredentials(); ok && basicCredential.Password != nil {
			statePasswordMap[*v.Name] = *basicCredential.Password
		}
	}
	result := make([]interface{}, 0)
	for _, v := range input {
		name := ""
		username := ""
		password := ""
		server := ""
		if v.Name != nil {
			name = *v.Name
			password = statePasswordMap[name]
		}
		if v.Properties != nil && v.Properties.Credentials != nil {
			if basicCredential, ok := v.Properties.Credentials.AsContainerRegistryBasicCredentials(); ok {
				if basicCredential.Username != nil {
					username = *basicCredential.Username
				}
				if basicCredential.Server != nil {
					server = *basicCredential.Server
				}
			}
		}

		result = append(result, map[string]interface{}{
			"name":     name,
			"username": username,
			"password": password,
			"server":   server,
		})
	}
	return result
}

func flattenSpringCloudBuildService(input *appplatform.BuildServiceProperties) []interface{} {
	if input == nil || input.ContainerRegistry == nil {
		return []interface{}{}
	}
	id, err := parse.SpringCloudContainerRegistryIDInsensitively(*input.ContainerRegistry)
	if err == nil {
		return []interface{}{
			map[string]interface{}{
				"container_registry_name": id.ContainerRegistryName,
			},
		}
	}
	return []interface{}{}
}

func flattenSpringCloudMarketplaceResource(input *appplatform.MarketplaceResource) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	plan := ""
	publisher := ""
	product := ""
	if input.Plan != nil {
		plan = *input.Plan
	}
	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	if input.Product != nil {
		product = *input.Product
	}
	return []interface{}{
		map[string]interface{}{
			"plan":      plan,
			"publisher": publisher,
			"product":   product,
		},
	}
}
