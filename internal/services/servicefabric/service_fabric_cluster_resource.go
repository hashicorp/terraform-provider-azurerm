// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicefabric

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	serviceFabricValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabric/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceFabricCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceFabricClusterCreateUpdate,
		Read:   resourceServiceFabricClusterRead,
		Update: resourceServiceFabricClusterCreateUpdate,
		Delete: resourceServiceFabricClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cluster.ParseClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"reliability_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cluster.ReliabilityLevelNone),
					string(cluster.ReliabilityLevelBronze),
					string(cluster.ReliabilityLevelSilver),
					string(cluster.ReliabilityLevelGold),
					string(cluster.ReliabilityLevelPlatinum),
				}, false),
			},

			"upgrade_mode": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cluster.UpgradeModeAutomatic),
					string(cluster.UpgradeModeManual),
				}, false),
			},

			"service_fabric_zonal_upgrade_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cluster.SfZonalUpgradeModeHierarchical),
					string(cluster.SfZonalUpgradeModeParallel),
				}, false),
			},

			"vmss_zonal_upgrade_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cluster.VMSSZonalUpgradeModeHierarchical),
					string(cluster.VMSSZonalUpgradeModeParallel),
				}, false),
			},

			"cluster_code_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"management_endpoint": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vm_image": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"add_on_features": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"azure_active_directory": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"cluster_application_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"client_application_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"certificate": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"certificate_common_names"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"thumbprint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"thumbprint_secondary": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"x509_store_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"certificate_common_names": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"certificate"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"common_names": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"certificate_common_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"certificate_issuer_thumbprint": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"x509_store_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"reverse_proxy_certificate": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"reverse_proxy_certificate_common_names"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"thumbprint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"thumbprint_secondary": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"x509_store_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"reverse_proxy_certificate_common_names": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"reverse_proxy_certificate"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"common_names": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"certificate_common_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"certificate_issuer_thumbprint": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"x509_store_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"client_certificate_thumbprint": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"thumbprint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"is_admin": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					},
				},
			},

			"client_certificate_common_name": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"common_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"issuer_thumbprint": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							// todo remove this when https://github.com/Azure/azure-sdk-for-go/issues/17744 is fixed
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"is_admin": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					},
				},
			},

			"diagnostics_config": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"storage_account_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"protected_account_key_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"blob_endpoint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"queue_endpoint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"table_endpoint": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"upgrade_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"force_restart_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"health_check_retry_timeout": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "00:45:00",
							ValidateFunc: serviceFabricValidate.UpgradeTimeout,
						},
						"health_check_stable_duration": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "00:01:00",
							ValidateFunc: serviceFabricValidate.UpgradeTimeout,
						},
						"health_check_wait_duration": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "00:00:30",
							ValidateFunc: serviceFabricValidate.UpgradeTimeout,
						},
						"upgrade_domain_timeout": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "02:00:00",
							ValidateFunc: serviceFabricValidate.UpgradeTimeout,
						},
						"upgrade_replica_set_check_timeout": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "10675199.02:48:05.4775807",
							ValidateFunc: serviceFabricValidate.UpgradeTimeout,
						},
						"upgrade_timeout": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "12:00:00",
							ValidateFunc: serviceFabricValidate.UpgradeTimeout,
						},
						"health_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"max_unhealthy_applications_percent": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"max_unhealthy_nodes_percent": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
						"delta_health_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"max_delta_unhealthy_applications_percent": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"max_delta_unhealthy_nodes_percent": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"max_upgrade_domain_delta_unhealthy_nodes_percent": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
					},
				},
			},

			"fabric_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"parameters": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"node_type": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"placement_properties": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"capacities": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"instance_count": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"is_primary": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"is_stateless": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"multiple_availability_zones": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"client_endpoint_port": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"http_endpoint_port": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"reverse_proxy_endpoint_port": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validate.PortNumber,
						},
						"durability_level": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cluster.DurabilityLevelBronze),
							ValidateFunc: validation.StringInSlice([]string{
								string(cluster.DurabilityLevelBronze),
								string(cluster.DurabilityLevelSilver),
								string(cluster.DurabilityLevelGold),
							}, false),
						},

						"application_ports": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start_port": {
										Type:     pluginsdk.TypeInt,
										Required: true,
									},
									"end_port": {
										Type:     pluginsdk.TypeInt,
										Required: true,
									},
								},
							},
						},

						"ephemeral_ports": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start_port": {
										Type:     pluginsdk.TypeInt,
										Required: true,
									},
									"end_port": {
										Type:     pluginsdk.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"cluster_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServiceFabricClusterCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabric.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cluster.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_service_fabric_cluster", id.ID())
		}
	}

	addOnFeaturesRaw := d.Get("add_on_features").(*pluginsdk.Set).List()
	addOnFeatures := expandServiceFabricClusterAddOnFeatures(addOnFeaturesRaw)

	azureActiveDirectoryRaw := d.Get("azure_active_directory").([]interface{})
	azureActiveDirectory := expandServiceFabricClusterAzureActiveDirectory(azureActiveDirectoryRaw)

	diagnosticsRaw := d.Get("diagnostics_config").([]interface{})
	diagnostics := expandServiceFabricClusterDiagnosticsConfig(diagnosticsRaw)

	upgradePolicyRaw := d.Get("upgrade_policy").([]interface{})
	upgradePolicy := expandServiceFabricClusterUpgradePolicy(upgradePolicyRaw)

	fabricSettingsRaw := d.Get("fabric_settings").([]interface{})
	fabricSettings := expandServiceFabricClusterFabricSettings(fabricSettingsRaw)

	nodeTypesRaw := d.Get("node_type").([]interface{})
	nodeTypes := expandServiceFabricClusterNodeTypes(nodeTypesRaw)

	location := d.Get("location").(string)
	reliabilityLevel := cluster.ReliabilityLevel(d.Get("reliability_level").(string))
	managementEndpoint := d.Get("management_endpoint").(string)
	upgradeMode := cluster.UpgradeMode(d.Get("upgrade_mode").(string))
	clusterCodeVersion := d.Get("cluster_code_version").(string)
	vmImage := d.Get("vm_image").(string)
	t := d.Get("tags").(map[string]interface{})

	clusterModel := cluster.Cluster{
		Location: location,
		Tags:     tags.Expand(t),
		Properties: &cluster.ClusterProperties{
			AddOnFeatures:                      addOnFeatures,
			AzureActiveDirectory:               azureActiveDirectory,
			CertificateCommonNames:             expandServiceFabricClusterCertificateCommonNames(d),
			ReverseProxyCertificateCommonNames: expandServiceFabricClusterReverseProxyCertificateCommonNames(d),
			DiagnosticsStorageAccountConfig:    diagnostics,
			FabricSettings:                     fabricSettings,
			ManagementEndpoint:                 managementEndpoint,
			NodeTypes:                          nodeTypes,
			ReliabilityLevel:                   &reliabilityLevel,
			UpgradeDescription:                 upgradePolicy,
			UpgradeMode:                        &upgradeMode,
			VmImage:                            utils.String(vmImage),
		},
	}

	if sfZonalUpgradeMode, ok := d.GetOk("service_fabric_zonal_upgrade_mode"); ok {
		mode := cluster.SfZonalUpgradeMode(sfZonalUpgradeMode.(string))
		clusterModel.Properties.SfZonalUpgradeMode = &mode
	}

	if vmssZonalUpgradeMode, ok := d.GetOk("vmss_zonal_upgrade_mode"); ok {
		mode := cluster.VMSSZonalUpgradeMode(vmssZonalUpgradeMode.(string))
		clusterModel.Properties.VMSSZonalUpgradeMode = &mode
	}

	if certificateRaw, ok := d.GetOk("certificate"); ok {
		certificate := expandServiceFabricClusterCertificate(certificateRaw.([]interface{}))
		clusterModel.Properties.Certificate = certificate
	}

	if reverseProxyCertificateRaw, ok := d.GetOk("reverse_proxy_certificate"); ok {
		reverseProxyCertificate := expandServiceFabricClusterReverseProxyCertificate(reverseProxyCertificateRaw.([]interface{}))
		clusterModel.Properties.ReverseProxyCertificate = reverseProxyCertificate
	}

	if clientCertificateThumbprintRaw, ok := d.GetOk("client_certificate_thumbprint"); ok {
		clientCertificateThumbprints := expandServiceFabricClusterClientCertificateThumbprints(clientCertificateThumbprintRaw.([]interface{}))
		clusterModel.Properties.ClientCertificateThumbprints = clientCertificateThumbprints
	}

	if clientCertificateCommonNamesRaw, ok := d.GetOk("client_certificate_common_name"); ok {
		clientCertificateCommonNames := expandServiceFabricClusterClientCertificateCommonNames(clientCertificateCommonNamesRaw.([]interface{}))
		clusterModel.Properties.ClientCertificateCommonNames = clientCertificateCommonNames
	}

	if clusterCodeVersion != "" {
		clusterModel.Properties.ClusterCodeVersion = utils.String(clusterCodeVersion)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, clusterModel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceServiceFabricClusterRead(d, meta)
}

func resourceServiceFabricClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabric.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cluster.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state!", id.ID())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("cluster_code_version", props.ClusterCodeVersion)
			d.Set("cluster_endpoint", props.ClusterEndpoint)
			d.Set("management_endpoint", props.ManagementEndpoint)
			d.Set("vm_image", props.VmImage)

			reliabilityLevel := ""
			if props.ReliabilityLevel != nil {
				reliabilityLevel = string(*props.ReliabilityLevel)
			}
			d.Set("reliability_level", reliabilityLevel)

			upgradeMode := ""
			if props.UpgradeMode != nil {
				upgradeMode = string(*props.UpgradeMode)
			}
			d.Set("upgrade_mode", upgradeMode)

			sfZonalMode := ""
			if props.SfZonalUpgradeMode != nil {
				sfZonalMode = string(*props.SfZonalUpgradeMode)
			}
			d.Set("service_fabric_zonal_upgrade_mode", sfZonalMode)

			vmssZonalMode := ""
			if props.VMSSZonalUpgradeMode != nil {
				vmssZonalMode = string(*props.VMSSZonalUpgradeMode)
			}
			d.Set("vmss_zonal_upgrade_mode", vmssZonalMode)

			addOnFeatures := flattenServiceFabricClusterAddOnFeatures(props.AddOnFeatures)
			if err := d.Set("add_on_features", pluginsdk.NewSet(pluginsdk.HashString, addOnFeatures)); err != nil {
				return fmt.Errorf("setting `add_on_features`: %+v", err)
			}

			azureActiveDirectory := flattenServiceFabricClusterAzureActiveDirectory(props.AzureActiveDirectory)
			if err := d.Set("azure_active_directory", azureActiveDirectory); err != nil {
				return fmt.Errorf("setting `azure_active_directory`: %+v", err)
			}

			certificate := flattenServiceFabricClusterCertificate(props.Certificate)
			if err := d.Set("certificate", certificate); err != nil {
				return fmt.Errorf("setting `certificate`: %+v", err)
			}

			certificateCommonNames := flattenServiceFabricClusterCertificateCommonNames(props.CertificateCommonNames)
			if err := d.Set("certificate_common_names", certificateCommonNames); err != nil {
				return fmt.Errorf("setting `certificate_common_names`: %+v", err)
			}

			reverseProxyCertificate := flattenServiceFabricClusterReverseProxyCertificate(props.ReverseProxyCertificate)
			if err := d.Set("reverse_proxy_certificate", reverseProxyCertificate); err != nil {
				return fmt.Errorf("setting `reverse_proxy_certificate`: %+v", err)
			}

			reverseProxyCertificateCommonNames := flattenServiceFabricClusterCertificateCommonNames(props.ReverseProxyCertificateCommonNames)
			if err := d.Set("reverse_proxy_certificate_common_names", reverseProxyCertificateCommonNames); err != nil {
				return fmt.Errorf("setting `reverse_proxy_certificate_common_names`: %+v", err)
			}

			clientCertificateThumbprints := flattenServiceFabricClusterClientCertificateThumbprints(props.ClientCertificateThumbprints)
			if err := d.Set("client_certificate_thumbprint", clientCertificateThumbprints); err != nil {
				return fmt.Errorf("setting `client_certificate_thumbprint`: %+v", err)
			}

			clientCertificateCommonNames := flattenServiceFabricClusterClientCertificateCommonNames(props.ClientCertificateCommonNames)
			if err := d.Set("client_certificate_common_name", clientCertificateCommonNames); err != nil {
				return fmt.Errorf("setting `client_certificate_common_name`: %+v", err)
			}

			diagnostics := flattenServiceFabricClusterDiagnosticsConfig(props.DiagnosticsStorageAccountConfig)
			if err := d.Set("diagnostics_config", diagnostics); err != nil {
				return fmt.Errorf("setting `diagnostics_config`: %+v", err)
			}

			upgradePolicy := flattenServiceFabricClusterUpgradePolicy(props.UpgradeDescription)
			if err := d.Set("upgrade_policy", upgradePolicy); err != nil {
				return fmt.Errorf("setting `upgrade_policy`: %+v", err)
			}

			fabricSettings := flattenServiceFabricClusterFabricSettings(props.FabricSettings)
			if err := d.Set("fabric_settings", fabricSettings); err != nil {
				return fmt.Errorf("setting `fabric_settings`: %+v", err)
			}

			nodeTypes := flattenServiceFabricClusterNodeTypes(props.NodeTypes)
			if err := d.Set("node_type", nodeTypes); err != nil {
				return fmt.Errorf("setting `node_type`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceServiceFabricClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabric.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cluster.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", id.ID())

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id.ID(), err)
		}
	}

	return nil
}

func expandServiceFabricClusterAddOnFeatures(input []interface{}) *[]cluster.AddOnFeatures {
	output := make([]cluster.AddOnFeatures, 0)

	for _, v := range input {
		output = append(output, cluster.AddOnFeatures(v.(string)))
	}

	return &output
}

func expandServiceFabricClusterAzureActiveDirectory(input []interface{}) *cluster.AzureActiveDirectory {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	tenantId := v["tenant_id"].(string)
	clusterApplication := v["cluster_application_id"].(string)
	clientApplication := v["client_application_id"].(string)

	config := cluster.AzureActiveDirectory{
		TenantId:           utils.String(tenantId),
		ClusterApplication: utils.String(clusterApplication),
		ClientApplication:  utils.String(clientApplication),
	}
	return &config
}

func flattenServiceFabricClusterAzureActiveDirectory(input *cluster.AzureActiveDirectory) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if name := v.TenantId; name != nil {
			output["tenant_id"] = *name
		}

		if name := v.ClusterApplication; name != nil {
			output["cluster_application_id"] = *name
		}

		if endpoint := v.ClientApplication; endpoint != nil {
			output["client_application_id"] = *endpoint
		}

		results = append(results, output)
	}

	return results
}

func flattenServiceFabricClusterAddOnFeatures(input *[]cluster.AddOnFeatures) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, string(v))
		}
	}

	return output
}

func expandServiceFabricClusterCertificate(input []interface{}) *cluster.CertificateDescription {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	thumbprint := v["thumbprint"].(string)
	x509StoreName := cluster.X509StoreName(v["x509_store_name"].(string))

	result := cluster.CertificateDescription{
		Thumbprint:    thumbprint,
		X509StoreName: &x509StoreName,
	}

	if thumb, ok := v["thumbprint_secondary"]; ok {
		result.ThumbprintSecondary = utils.String(thumb.(string))
	}

	return &result
}

func flattenServiceFabricClusterCertificate(input *cluster.CertificateDescription) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		output["thumbprint"] = input.Thumbprint

		if thumbprint := input.ThumbprintSecondary; thumbprint != nil {
			output["thumbprint_secondary"] = *thumbprint
		}

		output["x509_store_name"] = input.X509StoreName
		results = append(results, output)
	}

	return results
}

func expandServiceFabricClusterCertificateCommonNames(d *pluginsdk.ResourceData) *cluster.ServerCertificateCommonNames {
	i := d.Get("certificate_common_names").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	commonNamesRaw := input["common_names"].(*pluginsdk.Set).List()
	commonNames := make([]cluster.ServerCertificateCommonName, 0)

	for _, commonName := range commonNamesRaw {
		commonNameDetails := commonName.(map[string]interface{})

		commonName := cluster.ServerCertificateCommonName{
			CertificateCommonName:       commonNameDetails["certificate_common_name"].(string),
			CertificateIssuerThumbprint: commonNameDetails["certificate_issuer_thumbprint"].(string),
		}

		commonNames = append(commonNames, commonName)
	}

	x509StoreName := cluster.X509StoreName(input["x509_store_name"].(string))

	output := cluster.ServerCertificateCommonNames{
		CommonNames:   &commonNames,
		X509StoreName: &x509StoreName,
	}

	return &output
}

func expandServiceFabricClusterReverseProxyCertificateCommonNames(d *pluginsdk.ResourceData) *cluster.ServerCertificateCommonNames {
	i := d.Get("reverse_proxy_certificate_common_names").([]interface{})
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})

	commonNamesRaw := input["common_names"].(*pluginsdk.Set).List()
	commonNames := make([]cluster.ServerCertificateCommonName, 0)

	for _, commonName := range commonNamesRaw {
		commonNameDetails := commonName.(map[string]interface{})

		commonName := cluster.ServerCertificateCommonName{
			CertificateCommonName:       commonNameDetails["certificate_common_name"].(string),
			CertificateIssuerThumbprint: commonNameDetails["certificate_issuer_thumbprint"].(string),
		}

		commonNames = append(commonNames, commonName)
	}

	x509StoreName := cluster.X509StoreName(input["x509_store_name"].(string))

	output := cluster.ServerCertificateCommonNames{
		CommonNames:   &commonNames,
		X509StoreName: &x509StoreName,
	}

	return &output
}

func flattenServiceFabricClusterCertificateCommonNames(in *cluster.ServerCertificateCommonNames) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if commonNames := in.CommonNames; commonNames != nil {
		common_names := make([]map[string]interface{}, 0)
		for _, i := range *commonNames {
			commonName := make(map[string]interface{})

			commonName["certificate_common_name"] = i.CertificateCommonName
			commonName["certificate_issuer_thumbprint"] = i.CertificateIssuerThumbprint

			common_names = append(common_names, commonName)
		}

		output["common_names"] = common_names
	}

	output["x509_store_name"] = in.X509StoreName

	return []interface{}{output}
}

func expandServiceFabricClusterReverseProxyCertificate(input []interface{}) *cluster.CertificateDescription {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	x509StoreName := cluster.X509StoreName(v["x509_store_name"].(string))

	result := cluster.CertificateDescription{
		Thumbprint:    v["thumbprint"].(string),
		X509StoreName: &x509StoreName,
	}

	if thumb, ok := v["thumbprint_secondary"]; ok {
		result.ThumbprintSecondary = utils.String(thumb.(string))
	}

	return &result
}

func flattenServiceFabricClusterReverseProxyCertificate(input *cluster.CertificateDescription) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		output["thumbprint"] = input.Thumbprint

		if thumbprint := input.ThumbprintSecondary; thumbprint != nil {
			output["thumbprint_secondary"] = *thumbprint
		}

		output["x509_store_name"] = input.X509StoreName
		results = append(results, output)
	}

	return results
}

func expandServiceFabricClusterClientCertificateThumbprints(input []interface{}) *[]cluster.ClientCertificateThumbprint {
	results := make([]cluster.ClientCertificateThumbprint, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		result := cluster.ClientCertificateThumbprint{
			CertificateThumbprint: val["thumbprint"].(string),
			IsAdmin:               val["is_admin"].(bool),
		}
		results = append(results, result)
	}

	return &results
}

func flattenServiceFabricClusterClientCertificateThumbprints(input *[]cluster.ClientCertificateThumbprint) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		result["thumbprint"] = v.CertificateThumbprint
		result["is_admin"] = v.IsAdmin

		results = append(results, result)
	}

	return results
}

func expandServiceFabricClusterClientCertificateCommonNames(input []interface{}) *[]cluster.ClientCertificateCommonName {
	results := make([]cluster.ClientCertificateCommonName, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		certificateCommonName := val["common_name"].(string)
		certificateIssuerThumbprint := val["issuer_thumbprint"].(string)
		isAdmin := val["is_admin"].(bool)

		result := cluster.ClientCertificateCommonName{
			CertificateCommonName:       certificateCommonName,
			CertificateIssuerThumbprint: certificateIssuerThumbprint,
			IsAdmin:                     isAdmin,
		}
		results = append(results, result)
	}

	return &results
}

func flattenServiceFabricClusterClientCertificateCommonNames(input *[]cluster.ClientCertificateCommonName) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		result["common_name"] = v.CertificateCommonName
		result["issuer_thumbprint"] = v.CertificateIssuerThumbprint
		result["is_admin"] = v.IsAdmin

		results = append(results, result)
	}

	return results
}

func expandServiceFabricClusterDiagnosticsConfig(input []interface{}) *cluster.DiagnosticsStorageAccountConfig {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	config := cluster.DiagnosticsStorageAccountConfig{
		StorageAccountName:      v["storage_account_name"].(string),
		ProtectedAccountKeyName: v["protected_account_key_name"].(string),
		BlobEndpoint:            v["blob_endpoint"].(string),
		QueueEndpoint:           v["queue_endpoint"].(string),
		TableEndpoint:           v["table_endpoint"].(string),
	}
	return &config
}

func flattenServiceFabricClusterDiagnosticsConfig(input *cluster.DiagnosticsStorageAccountConfig) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		output["storage_account_name"] = v.StorageAccountName
		output["protected_account_key_name"] = v.ProtectedAccountKeyName
		output["blob_endpoint"] = v.BlobEndpoint
		output["queue_endpoint"] = v.QueueEndpoint
		output["table_endpoint"] = v.TableEndpoint

		results = append(results, output)
	}

	return results
}

func expandServiceFabricClusterUpgradePolicyDeltaHealthPolicy(input []interface{}) *cluster.ClusterUpgradeDeltaHealthPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	deltaHealthPolicy := &cluster.ClusterUpgradeDeltaHealthPolicy{}
	v := input[0].(map[string]interface{})
	deltaHealthPolicy.MaxPercentDeltaUnhealthyNodes = int64(v["max_delta_unhealthy_nodes_percent"].(int))
	deltaHealthPolicy.MaxPercentUpgradeDomainDeltaUnhealthyNodes = int64(v["max_upgrade_domain_delta_unhealthy_nodes_percent"].(int))
	deltaHealthPolicy.MaxPercentDeltaUnhealthyApplications = int64(v["max_delta_unhealthy_applications_percent"].(int))

	return deltaHealthPolicy
}

func expandServiceFabricClusterUpgradePolicyHealthPolicy(input []interface{}) cluster.ClusterHealthPolicy {
	healthPolicy := cluster.ClusterHealthPolicy{}
	if len(input) == 0 || input[0] == nil {
		return healthPolicy
	}

	v := input[0].(map[string]interface{})
	healthPolicy.MaxPercentUnhealthyApplications = utils.Int64(int64(v["max_unhealthy_applications_percent"].(int)))
	healthPolicy.MaxPercentUnhealthyNodes = utils.Int64(int64(v["max_unhealthy_nodes_percent"].(int)))

	return healthPolicy
}

func expandServiceFabricClusterUpgradePolicy(input []interface{}) *cluster.ClusterUpgradePolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	policy := &cluster.ClusterUpgradePolicy{}
	v := input[0].(map[string]interface{})

	policy.ForceRestart = utils.Bool(v["force_restart_enabled"].(bool))
	policy.HealthCheckStableDuration = v["health_check_stable_duration"].(string)
	policy.UpgradeDomainTimeout = v["upgrade_domain_timeout"].(string)
	policy.UpgradeReplicaSetCheckTimeout = v["upgrade_replica_set_check_timeout"].(string)
	policy.UpgradeTimeout = v["upgrade_timeout"].(string)
	policy.HealthCheckRetryTimeout = v["health_check_retry_timeout"].(string)
	policy.HealthCheckWaitDuration = v["health_check_wait_duration"].(string)

	if v["health_policy"] != nil {
		policy.HealthPolicy = expandServiceFabricClusterUpgradePolicyHealthPolicy(v["health_policy"].([]interface{}))
	}
	if v["delta_health_policy"] != nil {
		policy.DeltaHealthPolicy = expandServiceFabricClusterUpgradePolicyDeltaHealthPolicy(v["delta_health_policy"].([]interface{}))
	}

	return policy
}

func flattenServiceFabricClusterUpgradePolicy(input *cluster.ClusterUpgradePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if forceRestart := input.ForceRestart; forceRestart != nil {
		output["force_restart_enabled"] = *forceRestart
	}

	output["health_check_retry_timeout"] = input.HealthCheckRetryTimeout
	output["health_check_stable_duration"] = input.HealthCheckStableDuration
	output["health_check_wait_duration"] = input.HealthCheckWaitDuration
	output["upgrade_domain_timeout"] = input.UpgradeDomainTimeout
	output["upgrade_replica_set_check_timeout"] = input.UpgradeReplicaSetCheckTimeout
	output["upgrade_timeout"] = input.UpgradeTimeout

	output["health_policy"] = flattenServiceFabricClusterUpgradePolicyHealthPolicy(input.HealthPolicy)
	output["delta_health_policy"] = flattenServiceFabricClusterUpgradePolicyDeltaHealthPolicy(input.DeltaHealthPolicy)

	return []interface{}{output}
}

func flattenServiceFabricClusterUpgradePolicyHealthPolicy(input cluster.ClusterHealthPolicy) []interface{} {
	output := make(map[string]interface{})

	if input.MaxPercentUnhealthyApplications != nil {
		output["max_unhealthy_applications_percent"] = *input.MaxPercentUnhealthyApplications
	}

	if input.MaxPercentUnhealthyNodes != nil {
		output["max_unhealthy_nodes_percent"] = *input.MaxPercentUnhealthyNodes
	}

	return []interface{}{output}
}

func flattenServiceFabricClusterUpgradePolicyDeltaHealthPolicy(input *cluster.ClusterUpgradeDeltaHealthPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	output["max_delta_unhealthy_applications_percent"] = input.MaxPercentDeltaUnhealthyApplications
	output["max_delta_unhealthy_nodes_percent"] = input.MaxPercentDeltaUnhealthyNodes
	output["max_upgrade_domain_delta_unhealthy_nodes_percent"] = input.MaxPercentUpgradeDomainDeltaUnhealthyNodes

	return []interface{}{output}
}

func expandServiceFabricClusterFabricSettings(input []interface{}) *[]cluster.SettingsSectionDescription {
	results := make([]cluster.SettingsSectionDescription, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		params := make([]cluster.SettingsParameterDescription, 0)
		paramsRaw := val["parameters"].(map[string]interface{})
		for k, v := range paramsRaw {
			param := cluster.SettingsParameterDescription{
				Name:  k,
				Value: v.(string),
			}
			params = append(params, param)
		}

		result := cluster.SettingsSectionDescription{
			Name:       name,
			Parameters: params,
		}
		results = append(results, result)
	}

	return &results
}

func flattenServiceFabricClusterFabricSettings(input *[]cluster.SettingsSectionDescription) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		result := make(map[string]interface{})

		result["name"] = v.Name

		parameters := make(map[string]interface{})
		if paramsRaw := v.Parameters; paramsRaw != nil {
			for _, p := range paramsRaw {
				parameters[p.Name] = p.Value
			}
		}
		result["parameters"] = parameters
		results = append(results, result)
	}

	return results
}

func expandServiceFabricClusterNodeTypes(input []interface{}) []cluster.NodeTypeDescription {
	results := make([]cluster.NodeTypeDescription, 0)

	for _, v := range input {
		node := v.(map[string]interface{})

		durabilityLevel := cluster.DurabilityLevel(node["durability_level"].(string))

		result := cluster.NodeTypeDescription{
			Name:                         node["name"].(string),
			VMInstanceCount:              int64(node["instance_count"].(int)),
			IsPrimary:                    node["is_primary"].(bool),
			ClientConnectionEndpointPort: int64(node["client_endpoint_port"].(int)),
			HTTPGatewayEndpointPort:      int64(node["http_endpoint_port"].(int)),
			DurabilityLevel:              &durabilityLevel,
		}

		if isStateless, ok := node["is_stateless"]; ok {
			result.IsStateless = utils.Bool(isStateless.(bool))
		}

		if multipleAvailabilityZones, ok := node["multiple_availability_zones"]; ok {
			result.MultipleAvailabilityZones = utils.Bool(multipleAvailabilityZones.(bool))
		}

		if props, ok := node["placement_properties"]; ok {
			placementProperties := make(map[string]string)
			for key, value := range props.(map[string]interface{}) {
				placementProperties[key] = value.(string)
			}

			result.PlacementProperties = &placementProperties
		}

		if caps, ok := node["capacities"]; ok {
			capacities := make(map[string]string)
			for key, value := range caps.(map[string]interface{}) {
				capacities[key] = value.(string)
			}

			result.Capacities = &capacities
		}

		if v := int64(node["reverse_proxy_endpoint_port"].(int)); v != 0 {
			result.ReverseProxyEndpointPort = utils.Int64(v)
		}

		applicationPortsRaw := node["application_ports"].([]interface{})
		if len(applicationPortsRaw) > 0 {
			portsRaw := applicationPortsRaw[0].(map[string]interface{})

			result.ApplicationPorts = &cluster.EndpointRangeDescription{
				StartPort: int64(portsRaw["start_port"].(int)),
				EndPort:   int64(portsRaw["end_port"].(int)),
			}
		}

		ephemeralPortsRaw := node["ephemeral_ports"].([]interface{})
		if len(ephemeralPortsRaw) > 0 {
			portsRaw := ephemeralPortsRaw[0].(map[string]interface{})

			result.EphemeralPorts = &cluster.EndpointRangeDescription{
				StartPort: int64(portsRaw["start_port"].(int)),
				EndPort:   int64(portsRaw["end_port"].(int)),
			}
		}

		results = append(results, result)
	}

	return results
}

func flattenServiceFabricClusterNodeTypes(input []cluster.NodeTypeDescription) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, v := range input {
		output := make(map[string]interface{})

		output["name"] = v.Name
		output["instance_count"] = v.VMInstanceCount
		output["is_primary"] = v.IsPrimary
		output["client_endpoint_port"] = v.ClientConnectionEndpointPort
		output["http_endpoint_port"] = v.HTTPGatewayEndpointPort

		if placementProperties := v.PlacementProperties; placementProperties != nil {
			output["placement_properties"] = *placementProperties
		}

		if capacities := v.Capacities; capacities != nil {
			output["capacities"] = *capacities
		}

		if port := v.ReverseProxyEndpointPort; port != nil {
			output["reverse_proxy_endpoint_port"] = *port
		}

		if isStateless := v.IsStateless; isStateless != nil {
			output["is_stateless"] = *isStateless
		}

		if multipleAvailabilityZones := v.MultipleAvailabilityZones; multipleAvailabilityZones != nil {
			output["multiple_availability_zones"] = *multipleAvailabilityZones
		}

		if durabilityLevel := v.DurabilityLevel; durabilityLevel != nil {
			output["durability_level"] = string(*v.DurabilityLevel)
		}

		applicationPorts := make([]interface{}, 0)
		if ports := v.ApplicationPorts; ports != nil {
			r := make(map[string]interface{})
			r["start_port"] = int(ports.StartPort)
			r["end_port"] = int(ports.EndPort)
			applicationPorts = append(applicationPorts, r)
		}
		output["application_ports"] = applicationPorts

		ephemeralPorts := make([]interface{}, 0)
		if ports := v.EphemeralPorts; ports != nil {
			r := make(map[string]interface{})
			r["start_port"] = int(ports.StartPort)
			r["end_port"] = int(ports.EndPort)
			ephemeralPorts = append(ephemeralPorts, r)
		}
		output["ephemeral_ports"] = ephemeralPorts

		results = append(results, output)
	}

	return results
}
