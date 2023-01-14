package legacy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2022-08-01/compute"
)

// NOTE: the `azurerm_virtual_machine_scale_set` resource has been superseded by the
//
//	`azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources
//	and as such this resource is feature-frozen and new functionality will be added to these new resources instead.
func resourceVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetCreateUpdate,
		Read:   resourceVirtualMachineScaleSetRead,
		Update: resourceVirtualMachineScaleSetCreateUpdate,
		Delete: resourceVirtualMachineScaleSetDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LegacyVMSSV0ToV1{},
		}),

		// NOTE: @tombuildsstuff - don't remove with 3.0
		DeprecationMessage: `The 'azurerm_virtual_machine_scale_set' resource has been superseded by the 'azurerm_linux_virtual_machine_scale_set' and 'azurerm_windows_virtual_machine_scale_set' resources. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.`,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualMachineScaleSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"zones": {
				// @tombuildsstuff: since this is the legacy VMSS resource this is intentionally not using commonschema for consistency
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.ResourceIdentityTypeSystemAssigned),
								string(compute.ResourceIdentityTypeUserAssigned),
								string(compute.ResourceIdentityTypeSystemAssignedUserAssigned),
							}, false),
						},
						"identity_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: commonids.ValidateUserAssignedIdentityID,
							},
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"tier": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Windows_Client",
					"Windows_Server",
				}, false),
			},

			"upgrade_policy_mode": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.UpgradeModeAutomatic),
					string(compute.UpgradeModeManual),
					string(compute.UpgradeModeRolling),
				}, false),
			},

			"health_probe_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"automatic_os_upgrade": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"rolling_upgrade_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"max_batch_instance_percent": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      20,
							ValidateFunc: validation.IntBetween(5, 100),
						},

						"max_unhealthy_instance_percent": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      20,
							ValidateFunc: validation.IntBetween(5, 100),
						},

						"max_unhealthy_upgraded_instance_percent": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      20,
							ValidateFunc: validation.IntBetween(5, 100),
						},

						"pause_time_between_batches": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "PT0S",
							ValidateFunc: validate.ISO8601Duration,
						},
					},
				},
				DiffSuppressFunc: azureRmVirtualMachineScaleSetSuppressRollingUpgradePolicyDiff,
			},

			"overprovision": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"single_placement_group": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"priority": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.VirtualMachinePriorityTypesLow),
					string(compute.VirtualMachinePriorityTypesRegular),
				}, false),
			},

			"eviction_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.VirtualMachineEvictionPolicyTypesDeallocate),
					string(compute.VirtualMachineEvictionPolicyTypesDelete),
				}, false),
			},

			"os_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"computer_name_prefix": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},

						"admin_username": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"admin_password": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"custom_data": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							StateFunc:        userDataStateFunc,
							DiffSuppressFunc: userDataDiffSuppressFunc,
						},
					},
				},
			},

			"os_profile_secrets": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"source_vault_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"vault_certificates": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"certificate_url": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"certificate_store": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			//lintignore:S018
			"os_profile_windows_config": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"provision_vm_agent": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"enable_automatic_upgrades": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"winrm": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"protocol": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"certificate_url": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
								},
							},
						},
						"additional_unattend_config": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"pass": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"component": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"setting_name": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"content": {
										Type:      pluginsdk.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
				Set: resourceVirtualMachineScaleSetOsProfileWindowsConfigHash,
			},

			//lintignore:S018
			"os_profile_linux_config": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disable_password_authentication": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"ssh_keys": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"path": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"key_data": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
				Set: resourceVirtualMachineScaleSetOsProfileLinuxConfigHash,
			},

			//lintignore:S018
			"network_profile": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"primary": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},

						"accelerated_networking": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"ip_forwarding": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"network_security_group_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"dns_settings": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"dns_servers": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},

						"ip_configuration": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"subnet_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},

									"application_gateway_backend_address_pool_ids": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"application_security_group_ids": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: azure.ValidateResourceID,
										},
										Set:      pluginsdk.HashString,
										MaxItems: 20,
									},

									"load_balancer_backend_address_pool_ids": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"load_balancer_inbound_nat_rules_ids": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"primary": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},

									"public_ip_address_configuration": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:     pluginsdk.TypeString,
													Required: true,
												},

												"idle_timeout": {
													Type:         pluginsdk.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(4, 32),
												},

												"domain_name_label": {
													Type:     pluginsdk.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Set: resourceVirtualMachineScaleSetNetworkConfigurationHash,
			},

			"boot_diagnostics": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"storage_uri": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			//lintignore:S018
			"storage_profile_os_disk": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"image": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"vhd_containers": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"managed_disk_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StorageAccountTypesPremiumLRS),
								string(compute.StorageAccountTypesStandardLRS),
								string(compute.StorageAccountTypesStandardSSDLRS),
							}, false),
						},

						"caching": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"os_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"create_option": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceVirtualMachineScaleSetStorageProfileOsDiskHash,
			},

			"storage_profile_data_disk": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"lun": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"create_option": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"caching": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate2.DiskSizeGB,
						},

						"managed_disk_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StorageAccountTypesPremiumLRS),
								string(compute.StorageAccountTypesStandardLRS),
								string(compute.StorageAccountTypesStandardSSDLRS),
							}, false),
						},
					},
				},
			},

			//lintignore:S018
			"storage_profile_image_reference": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"publisher": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"offer": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"sku": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceVirtualMachineScaleSetStorageProfileImageReferenceHash,
			},

			//lintignore:S018
			"plan": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"publisher": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"product": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			//lintignore:S018
			"extension": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"publisher": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"type_handler_version": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"auto_upgrade_minor_version": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"provision_after_extensions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							Set: pluginsdk.HashString,
						},

						"settings": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
						},

						"protected_settings": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
						},
					},
				},
				Set: resourceVirtualMachineScaleSetExtensionHash,
			},

			"proximity_placement_group_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,

				// We have to ignore case due to incorrect capitalisation of resource group name in
				// proximity placement group ID in the response we get from the API request
				//
				// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(azureRmVirtualMachineScaleSetCustomizeDiff),
	}
}

func resourceVirtualMachineScaleSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Legacy.VMScaleSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine Scale Set creation.")

	id := parse.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	zones := expandZones(d.Get("zones").([]interface{}))

	sku := expandVirtualMachineScaleSetSku(d)

	storageProfile := compute.VirtualMachineScaleSetStorageProfile{}
	osDisk, err := expandAzureRMVirtualMachineScaleSetsStorageProfileOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile.OsDisk = osDisk

	if _, ok := d.GetOk("storage_profile_data_disk"); ok {
		storageProfile.DataDisks = expandAzureRMVirtualMachineScaleSetsStorageProfileDataDisk(d)
	}

	if _, ok := d.GetOk("storage_profile_image_reference"); ok {
		imageRef, err2 := expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d)
		if err2 != nil {
			return err2
		}
		storageProfile.ImageReference = imageRef
	}

	osProfile := expandAzureRMVirtualMachineScaleSetsOsProfile(d)
	if err != nil {
		return err
	}

	extensions, err := expandAzureRMVirtualMachineScaleSetExtensions(d)
	if err != nil {
		return err
	}

	upgradePolicy := d.Get("upgrade_policy_mode").(string)
	automaticOsUpgrade := d.Get("automatic_os_upgrade").(bool)
	overprovision := d.Get("overprovision").(bool)
	singlePlacementGroup := d.Get("single_placement_group").(bool)
	priority := d.Get("priority").(string)
	evictionPolicy := d.Get("eviction_policy").(string)

	scaleSetProps := compute.VirtualMachineScaleSetProperties{
		UpgradePolicy: &compute.UpgradePolicy{
			Mode: compute.UpgradeMode(upgradePolicy),
			AutomaticOSUpgradePolicy: &compute.AutomaticOSUpgradePolicy{
				EnableAutomaticOSUpgrade: utils.Bool(automaticOsUpgrade),
			},
			RollingUpgradePolicy: expandAzureRmRollingUpgradePolicy(d),
		},
		VirtualMachineProfile: &compute.VirtualMachineScaleSetVMProfile{
			NetworkProfile:   expandAzureRmVirtualMachineScaleSetNetworkProfile(d),
			StorageProfile:   &storageProfile,
			OsProfile:        osProfile,
			ExtensionProfile: extensions,
			Priority:         compute.VirtualMachinePriorityTypes(priority),
		},
		// OrchestrationMode needs to be hardcoded to Uniform, for the
		// standard VMSS resource, since virtualMachineProfile is now supported
		// in both VMSS and Orchestrated VMSS...
		OrchestrationMode:    compute.OrchestrationModeUniform,
		Overprovision:        &overprovision,
		SinglePlacementGroup: &singlePlacementGroup,
	}

	if strings.EqualFold(priority, string(compute.VirtualMachinePriorityTypesLow)) {
		scaleSetProps.VirtualMachineProfile.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(evictionPolicy)
	}

	if _, ok := d.GetOk("boot_diagnostics"); ok {
		diagnosticProfile := expandAzureRMVirtualMachineScaleSetsDiagnosticProfile(d)
		scaleSetProps.VirtualMachineProfile.DiagnosticsProfile = &diagnosticProfile
	}

	if v, ok := d.GetOk("health_probe_id"); ok {
		scaleSetProps.VirtualMachineProfile.NetworkProfile.HealthProbe = &compute.APIEntityReference{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		scaleSetProps.ProximityPlacementGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	properties := compute.VirtualMachineScaleSet{
		Name:                             &id.Name,
		Location:                         &location,
		Tags:                             tags.Expand(t),
		Sku:                              sku,
		VirtualMachineScaleSetProperties: &scaleSetProps,
		Zones:                            zones,
	}

	if _, ok := d.GetOk("identity"); ok {
		properties.Identity = expandAzureRmVirtualMachineScaleSetIdentity(d)
	}

	if v, ok := d.GetOk("license_type"); ok {
		properties.VirtualMachineProfile.LicenseType = utils.String(v.(string))
	}

	if _, ok := d.GetOk("plan"); ok {
		properties.Plan = expandAzureRmVirtualMachineScaleSetPlan(d)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, properties)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetRead(d, meta)
}

func resourceVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Legacy.VMScaleSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] AzureRM Virtual Machine Scale Set (%s) Not Found. Removing from State", id.Name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Virtual Machine Scale Set %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("zones", resp.Zones)

	if err := d.Set("sku", flattenAzureRmVirtualMachineScaleSetSku(resp.Sku)); err != nil {
		return fmt.Errorf("[DEBUG] setting `sku`: %#v", err)
	}

	flattenedIdentity, err := flattenAzureRmVirtualMachineScaleSetIdentity(resp.Identity)
	if err != nil {
		return err
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("[DEBUG] setting `identity`: %+v", err)
	}

	if properties := resp.VirtualMachineScaleSetProperties; properties != nil {
		if upgradePolicy := properties.UpgradePolicy; upgradePolicy != nil {
			d.Set("upgrade_policy_mode", upgradePolicy.Mode)
			if policy := upgradePolicy.AutomaticOSUpgradePolicy; policy != nil {
				d.Set("automatic_os_upgrade", policy.EnableAutomaticOSUpgrade)
			}

			if rollingUpgradePolicy := upgradePolicy.RollingUpgradePolicy; rollingUpgradePolicy != nil {
				if err := d.Set("rolling_upgrade_policy", flattenAzureRmVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicy)); err != nil {
					return fmt.Errorf("[DEBUG] setting Virtual Machine Scale Set Rolling Upgrade Policy error: %#v", err)
				}
			}

			if proximityPlacementGroup := properties.ProximityPlacementGroup; proximityPlacementGroup != nil {
				d.Set("proximity_placement_group_id", proximityPlacementGroup.ID)
			}
		}
		d.Set("overprovision", properties.Overprovision)
		d.Set("single_placement_group", properties.SinglePlacementGroup)

		if profile := properties.VirtualMachineProfile; profile != nil {
			d.Set("license_type", profile.LicenseType)
			d.Set("priority", string(profile.Priority))
			d.Set("eviction_policy", string(profile.EvictionPolicy))

			osProfile := flattenAzureRMVirtualMachineScaleSetOsProfile(d, profile.OsProfile)
			if err := d.Set("os_profile", osProfile); err != nil {
				return fmt.Errorf("[DEBUG] setting `os_profile`: %#v", err)
			}

			if osProfile := profile.OsProfile; osProfile != nil {
				if linuxConfiguration := osProfile.LinuxConfiguration; linuxConfiguration != nil {
					flattenedLinuxConfiguration := flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(linuxConfiguration)
					if err := d.Set("os_profile_linux_config", flattenedLinuxConfiguration); err != nil {
						return fmt.Errorf("[DEBUG] setting `os_profile_linux_config`: %#v", err)
					}
				}

				if secrets := osProfile.Secrets; secrets != nil {
					flattenedSecrets := flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(secrets)
					if err := d.Set("os_profile_secrets", flattenedSecrets); err != nil {
						return fmt.Errorf("[DEBUG] setting `os_profile_secrets`: %#v", err)
					}
				}

				if windowsConfiguration := osProfile.WindowsConfiguration; windowsConfiguration != nil {
					flattenedWindowsConfiguration := flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(windowsConfiguration)
					if err := d.Set("os_profile_windows_config", flattenedWindowsConfiguration); err != nil {
						return fmt.Errorf("[DEBUG] setting `os_profile_windows_config`: %#v", err)
					}
				}
			}

			if diagnosticsProfile := profile.DiagnosticsProfile; diagnosticsProfile != nil {
				if bootDiagnostics := diagnosticsProfile.BootDiagnostics; bootDiagnostics != nil {
					flattenedDiagnostics := flattenAzureRmVirtualMachineScaleSetBootDiagnostics(bootDiagnostics)
					if err := d.Set("boot_diagnostics", flattenedDiagnostics); err != nil {
						return fmt.Errorf("[DEBUG] setting `boot_diagnostics`: %#v", err)
					}
				}
			}

			if networkProfile := profile.NetworkProfile; networkProfile != nil {
				if hp := networkProfile.HealthProbe; hp != nil {
					if id := hp.ID; id != nil {
						d.Set("health_probe_id", id)
					}
				}

				flattenedNetworkProfile := flattenAzureRmVirtualMachineScaleSetNetworkProfile(networkProfile)
				if err := d.Set("network_profile", flattenedNetworkProfile); err != nil {
					return fmt.Errorf("[DEBUG] setting `network_profile`: %#v", err)
				}
			}

			if storageProfile := profile.StorageProfile; storageProfile != nil {
				if dataDisks := resp.VirtualMachineProfile.StorageProfile.DataDisks; dataDisks != nil {
					flattenedDataDisks := flattenAzureRmVirtualMachineScaleSetStorageProfileDataDisk(dataDisks)
					if err := d.Set("storage_profile_data_disk", flattenedDataDisks); err != nil {
						return fmt.Errorf("[DEBUG] setting `storage_profile_data_disk`: %#v", err)
					}
				}

				if imageRef := storageProfile.ImageReference; imageRef != nil {
					flattenedImageRef := flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(imageRef)
					if err := d.Set("storage_profile_image_reference", flattenedImageRef); err != nil {
						return fmt.Errorf("[DEBUG] setting `storage_profile_image_reference`: %#v", err)
					}
				}

				if osDisk := storageProfile.OsDisk; osDisk != nil {
					flattenedOSDisk := flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(osDisk)
					if err := d.Set("storage_profile_os_disk", flattenedOSDisk); err != nil {
						return fmt.Errorf("[DEBUG] setting `storage_profile_os_disk`: %#v", err)
					}
				}
			}

			if extensionProfile := properties.VirtualMachineProfile.ExtensionProfile; extensionProfile != nil {
				extension, err := flattenAzureRmVirtualMachineScaleSetExtensionProfile(extensionProfile)
				if err != nil {
					return fmt.Errorf("[DEBUG] setting Virtual Machine Scale Set Extension Profile error: %#v", err)
				}
				if err := d.Set("extension", extension); err != nil {
					return fmt.Errorf("[DEBUG] setting `extension`: %#v", err)
				}
			}
		}
	}

	if plan := resp.Plan; plan != nil {
		flattenedPlan := flattenAzureRmVirtualMachineScaleSetPlan(plan)
		if err := d.Set("plan", flattenedPlan); err != nil {
			return fmt.Errorf("[DEBUG] setting `plan`: %#v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Legacy.VMScaleSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// @ArcturusZhang (mimicking from virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	var forceDeletion *bool = nil
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, forceDeletion)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	return nil
}

func flattenAzureRmVirtualMachineScaleSetIdentity(identity *compute.VirtualMachineScaleSetIdentity) ([]interface{}, error) {
	if identity == nil {
		return make([]interface{}, 0), nil
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}

	identityIds := make([]string, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}
	result["identity_ids"] = identityIds

	return []interface{}{result}, nil
}

func flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(config *compute.LinuxConfiguration) []interface{} {
	result := make(map[string]interface{})

	if v := config.DisablePasswordAuthentication; v != nil {
		result["disable_password_authentication"] = *v
	}

	if ssh := config.SSH; ssh != nil {
		if keys := ssh.PublicKeys; keys != nil {
			ssh_keys := make([]map[string]interface{}, 0, len(*keys))
			for _, i := range *keys {
				key := make(map[string]interface{})

				if i.Path != nil {
					key["path"] = *i.Path
				}

				if i.KeyData != nil {
					key["key_data"] = *i.KeyData
				}

				ssh_keys = append(ssh_keys, key)
			}

			result["ssh_keys"] = ssh_keys
		}
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(config *compute.WindowsConfiguration) []interface{} {
	result := make(map[string]interface{})

	if config.ProvisionVMAgent != nil {
		result["provision_vm_agent"] = *config.ProvisionVMAgent
	}

	if config.EnableAutomaticUpdates != nil {
		result["enable_automatic_upgrades"] = *config.EnableAutomaticUpdates
	}

	if config.WinRM != nil {
		listeners := make([]map[string]interface{}, 0, len(*config.WinRM.Listeners))
		for _, i := range *config.WinRM.Listeners {
			listener := make(map[string]interface{})
			listener["protocol"] = i.Protocol

			if i.CertificateURL != nil {
				listener["certificate_url"] = *i.CertificateURL
			}

			listeners = append(listeners, listener)
		}

		result["winrm"] = listeners
	}

	if config.AdditionalUnattendContent != nil {
		content := make([]map[string]interface{}, 0, len(*config.AdditionalUnattendContent))
		for _, i := range *config.AdditionalUnattendContent {
			c := make(map[string]interface{})
			c["pass"] = i.PassName
			c["component"] = i.ComponentName
			c["setting_name"] = i.SettingName

			if i.Content != nil {
				c["content"] = *i.Content
			}

			content = append(content, c)
		}

		result["additional_unattend_config"] = content
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(secrets *[]compute.VaultSecretGroup) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(*secrets))
	for _, secret := range *secrets {
		s := map[string]interface{}{
			"source_vault_id": *secret.SourceVault.ID,
		}

		if secret.VaultCertificates != nil {
			certs := make([]map[string]interface{}, 0, len(*secret.VaultCertificates))
			for _, cert := range *secret.VaultCertificates {
				vaultCert := make(map[string]interface{})
				vaultCert["certificate_url"] = *cert.CertificateURL

				if cert.CertificateStore != nil {
					vaultCert["certificate_store"] = *cert.CertificateStore
				}

				certs = append(certs, vaultCert)
			}

			s["vault_certificates"] = certs
		}

		result = append(result, s)
	}
	return result
}

func flattenAzureRmVirtualMachineScaleSetBootDiagnostics(bootDiagnostic *compute.BootDiagnostics) []interface{} {
	b := make(map[string]interface{})

	if bootDiagnostic.Enabled != nil {
		b["enabled"] = *bootDiagnostic.Enabled
	}

	if bootDiagnostic.StorageURI != nil {
		b["storage_uri"] = *bootDiagnostic.StorageURI
	}

	return []interface{}{b}
}

func flattenAzureRmVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicy *compute.RollingUpgradePolicy) []interface{} {
	b := make(map[string]interface{})

	if v := rollingUpgradePolicy.MaxBatchInstancePercent; v != nil {
		b["max_batch_instance_percent"] = *v
	}
	if v := rollingUpgradePolicy.MaxUnhealthyInstancePercent; v != nil {
		b["max_unhealthy_instance_percent"] = *v
	}
	if v := rollingUpgradePolicy.MaxUnhealthyUpgradedInstancePercent; v != nil {
		b["max_unhealthy_upgraded_instance_percent"] = *v
	}
	if v := rollingUpgradePolicy.PauseTimeBetweenBatches; v != nil {
		b["pause_time_between_batches"] = *v
	}

	return []interface{}{b}
}

func flattenAzureRmVirtualMachineScaleSetNetworkProfile(profile *compute.VirtualMachineScaleSetNetworkProfile) []map[string]interface{} {
	networkConfigurations := profile.NetworkInterfaceConfigurations
	result := make([]map[string]interface{}, 0, len(*networkConfigurations))
	for _, netConfig := range *networkConfigurations {
		s := map[string]interface{}{
			"name":    *netConfig.Name,
			"primary": *netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.Primary,
		}

		if v := netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.EnableAcceleratedNetworking; v != nil {
			s["accelerated_networking"] = *v
		}

		if v := netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.EnableIPForwarding; v != nil {
			s["ip_forwarding"] = *v
		}

		if v := netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.NetworkSecurityGroup; v != nil {
			s["network_security_group_id"] = *v.ID
		}

		if dnsSettings := netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.DNSSettings; dnsSettings != nil {
			dnsServers := make([]string, 0)
			if s := dnsSettings.DNSServers; s != nil {
				dnsServers = *s
			}

			s["dns_settings"] = []interface{}{map[string]interface{}{
				"dns_servers": dnsServers,
			}}
		}

		if netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations != nil {
			ipConfigs := make([]map[string]interface{}, 0, len(*netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations))
			for _, ipConfig := range *netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations {
				config := make(map[string]interface{})
				config["name"] = *ipConfig.Name

				if properties := ipConfig.VirtualMachineScaleSetIPConfigurationProperties; properties != nil {
					if properties.Subnet != nil {
						config["subnet_id"] = *properties.Subnet.ID
					}

					addressPools := make([]interface{}, 0)
					if properties.ApplicationGatewayBackendAddressPools != nil {
						for _, pool := range *properties.ApplicationGatewayBackendAddressPools {
							if v := pool.ID; v != nil {
								addressPools = append(addressPools, *v)
							}
						}
					}
					config["application_gateway_backend_address_pool_ids"] = pluginsdk.NewSet(pluginsdk.HashString, addressPools)

					applicationSecurityGroups := make([]interface{}, 0)
					if properties.ApplicationSecurityGroups != nil {
						for _, asg := range *properties.ApplicationSecurityGroups {
							if v := asg.ID; v != nil {
								applicationSecurityGroups = append(applicationSecurityGroups, *v)
							}
						}
					}
					config["application_security_group_ids"] = pluginsdk.NewSet(pluginsdk.HashString, applicationSecurityGroups)

					if properties.LoadBalancerBackendAddressPools != nil {
						addressPools := make([]interface{}, 0, len(*properties.LoadBalancerBackendAddressPools))
						for _, pool := range *properties.LoadBalancerBackendAddressPools {
							if v := pool.ID; v != nil {
								addressPools = append(addressPools, *v)
							}
						}
						config["load_balancer_backend_address_pool_ids"] = pluginsdk.NewSet(pluginsdk.HashString, addressPools)
					}

					if properties.LoadBalancerInboundNatPools != nil {
						inboundNatPools := make([]interface{}, 0, len(*properties.LoadBalancerInboundNatPools))
						for _, rule := range *properties.LoadBalancerInboundNatPools {
							if v := rule.ID; v != nil {
								inboundNatPools = append(inboundNatPools, *v)
							}
						}
						config["load_balancer_inbound_nat_rules_ids"] = pluginsdk.NewSet(pluginsdk.HashString, inboundNatPools)
					}

					if properties.Primary != nil {
						config["primary"] = *properties.Primary
					}

					if publicIpInfo := properties.PublicIPAddressConfiguration; publicIpInfo != nil {
						publicIpConfigs := make([]map[string]interface{}, 0, 1)
						publicIpConfig := make(map[string]interface{})
						if publicIpName := publicIpInfo.Name; publicIpName != nil {
							publicIpConfig["name"] = *publicIpName
						}
						if publicIpProperties := publicIpInfo.VirtualMachineScaleSetPublicIPAddressConfigurationProperties; publicIpProperties != nil {
							if dns := publicIpProperties.DNSSettings; dns != nil {
								publicIpConfig["domain_name_label"] = *dns.DomainNameLabel
							}
							if timeout := publicIpProperties.IdleTimeoutInMinutes; timeout != nil {
								publicIpConfig["idle_timeout"] = *timeout
							}
							publicIpConfigs = append(publicIpConfigs, publicIpConfig)
						}
						config["public_ip_address_configuration"] = publicIpConfigs
					}

					ipConfigs = append(ipConfigs, config)
				}
			}

			s["ip_configuration"] = ipConfigs
		}

		result = append(result, s)
	}

	return result
}

func flattenAzureRMVirtualMachineScaleSetOsProfile(d *pluginsdk.ResourceData, profile *compute.VirtualMachineScaleSetOSProfile) []interface{} {
	result := make(map[string]interface{})

	if profile != nil {

		if profile.ComputerNamePrefix != nil {
			result["computer_name_prefix"] = *profile.ComputerNamePrefix
		} else {
			result["computer_name_prefix"] = ""
		}
		if profile.AdminUsername != nil {
			result["admin_username"] = *profile.AdminUsername
		} else {
			result["admin_username"] = ""
		}

		// admin password isn't returned, so let's look it up
		if v, ok := d.GetOk("os_profile.0.admin_password"); ok {
			password := v.(string)
			result["admin_password"] = password
		}

		if profile.CustomData != nil {
			result["custom_data"] = *profile.CustomData
		} else {
			// look up the current custom data
			result["custom_data"] = utils.Base64EncodeIfNot(d.Get("os_profile.0.custom_data").(string))
		}
	} else {
		result["computer_name_prefix"] = ""
		result["admin_username"] = ""
		result["admin_password"] = ""
		result["custom_data"] = ""
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(profile *compute.VirtualMachineScaleSetOSDisk) []interface{} {
	result := make(map[string]interface{})

	if profile.Name != nil {
		result["name"] = *profile.Name
	}

	if profile.Image != nil {
		result["image"] = *profile.Image.URI
	}

	containers := make([]interface{}, 0)
	if profile.VhdContainers != nil {
		for _, container := range *profile.VhdContainers {
			containers = append(containers, container)
		}
	}
	result["vhd_containers"] = pluginsdk.NewSet(pluginsdk.HashString, containers)

	if profile.ManagedDisk != nil {
		result["managed_disk_type"] = string(profile.ManagedDisk.StorageAccountType)
	}

	result["caching"] = profile.Caching
	result["create_option"] = profile.CreateOption
	result["os_type"] = profile.OsType

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileDataDisk(disks *[]compute.VirtualMachineScaleSetDataDisk) interface{} {
	result := make([]interface{}, len(*disks))
	for i, disk := range *disks {
		l := make(map[string]interface{})
		if disk.ManagedDisk != nil {
			l["managed_disk_type"] = string(disk.ManagedDisk.StorageAccountType)
		}

		l["create_option"] = disk.CreateOption
		l["caching"] = string(disk.Caching)
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}
		if v := disk.Lun; v != nil {
			l["lun"] = *v
		}

		result[i] = l
	}
	return result
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(image *compute.ImageReference) []interface{} {
	result := make(map[string]interface{})
	if image.Publisher != nil {
		result["publisher"] = *image.Publisher
	}
	if image.Offer != nil {
		result["offer"] = *image.Offer
	}
	if image.Sku != nil {
		result["sku"] = *image.Sku
	}
	if image.Version != nil {
		result["version"] = *image.Version
	}
	if image.ID != nil {
		result["id"] = *image.ID
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetSku(sku *compute.Sku) []interface{} {
	result := make(map[string]interface{})
	if sku != nil {
		if sku.Name != nil {
			result["name"] = *sku.Name
		}

		if sku.Capacity != nil {
			result["capacity"] = *sku.Capacity
		}

		if sku.Tier != nil && *sku.Tier != "" {
			result["tier"] = *sku.Tier
		}
	}
	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetExtensionProfile(profile *compute.VirtualMachineScaleSetExtensionProfile) ([]map[string]interface{}, error) {
	if profile.Extensions == nil {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, len(*profile.Extensions))
	for _, extension := range *profile.Extensions {
		e := make(map[string]interface{})
		e["name"] = *extension.Name
		properties := extension.VirtualMachineScaleSetExtensionProperties
		if properties != nil {
			e["publisher"] = *properties.Publisher
			e["type"] = *properties.Type
			e["type_handler_version"] = *properties.TypeHandlerVersion
			if properties.AutoUpgradeMinorVersion != nil {
				e["auto_upgrade_minor_version"] = *properties.AutoUpgradeMinorVersion
			}

			provisionAfterExtensions := make([]interface{}, 0)
			if properties.ProvisionAfterExtensions != nil {
				for _, provisionAfterExtension := range *properties.ProvisionAfterExtensions {
					provisionAfterExtensions = append(provisionAfterExtensions, provisionAfterExtension)
				}
			}
			e["provision_after_extensions"] = pluginsdk.NewSet(pluginsdk.HashString, provisionAfterExtensions)

			if settings := properties.Settings; settings != nil {
				settingsVal := settings.(map[string]interface{})
				settingsJson, err := pluginsdk.FlattenJsonToString(settingsVal)
				if err != nil {
					return nil, err
				}
				e["settings"] = settingsJson
			}
		}

		result = append(result, e)
	}

	return result, nil
}

func resourceVirtualMachineScaleSetStorageProfileImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["publisher"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["offer"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["sku"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["version"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetStorageProfileOsDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))

		if v, ok := m["vhd_containers"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(*pluginsdk.Set).List()))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetNetworkConfigurationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["primary"].(bool)))

		if v, ok := m["accelerated_networking"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
		if v, ok := m["ip_forwarding"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
		if v, ok := m["network_security_group_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["dns_settings"].(map[string]interface{}); ok {
			if k, ok := v["dns_servers"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", k))
			}
		}
		if ipConfig, ok := m["ip_configuration"].([]interface{}); ok {
			for _, it := range ipConfig {
				config := it.(map[string]interface{})
				if name, ok := config["name"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", name.(string)))
				}
				if subnetid, ok := config["subnet_id"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", subnetid.(string)))
				}
				if appPoolId, ok := config["application_gateway_backend_address_pool_ids"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", appPoolId.(*pluginsdk.Set).List()))
				}
				if appSecGroup, ok := config["application_security_group_ids"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", appSecGroup.(*pluginsdk.Set).List()))
				}
				if lbPoolIds, ok := config["load_balancer_backend_address_pool_ids"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", lbPoolIds.(*pluginsdk.Set).List()))
				}
				if lbInNatRules, ok := config["load_balancer_inbound_nat_rules_ids"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", lbInNatRules.(*pluginsdk.Set).List()))
				}
				if primary, ok := config["primary"]; ok {
					buf.WriteString(fmt.Sprintf("%t-", primary.(bool)))
				}
				if publicIPConfig, ok := config["public_ip_address_configuration"].([]interface{}); ok {
					for _, publicIPIt := range publicIPConfig {
						publicip := publicIPIt.(map[string]interface{})
						if publicIPConfigName, ok := publicip["name"]; ok {
							buf.WriteString(fmt.Sprintf("%s-", publicIPConfigName.(string)))
						}
						if idle_timeout, ok := publicip["idle_timeout"]; ok {
							buf.WriteString(fmt.Sprintf("%d-", idle_timeout.(int)))
						}
						if dnsLabel, ok := publicip["domain_name_label"]; ok {
							buf.WriteString(fmt.Sprintf("%s-", dnsLabel.(string)))
						}
					}
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))

		if sshKeys, ok := m["ssh_keys"].([]interface{}); ok {
			for _, item := range sshKeys {
				k := item.(map[string]interface{})
				if path, ok := k["path"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", path.(string)))
				}
				if data, ok := k["key_data"]; ok {
					buf.WriteString(fmt.Sprintf("%s-", data.(string)))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetOsProfileWindowsConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["provision_vm_agent"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
		if v, ok := m["enable_automatic_upgrades"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetExtensionHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type_handler_version"].(string)))

		if v, ok := m["auto_upgrade_minor_version"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}

		if v, ok := m["provision_after_extensions"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(*pluginsdk.Set).List()))
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
			if err == nil {
				serializedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
				if err == nil {
					buf.WriteString(fmt.Sprintf("%s-", serializedSettings))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func expandVirtualMachineScaleSetSku(d *pluginsdk.ResourceData) *compute.Sku {
	skuConfig := d.Get("sku").([]interface{})
	config := skuConfig[0].(map[string]interface{})

	sku := &compute.Sku{
		Name:     utils.String(config["name"].(string)),
		Capacity: utils.Int64(int64(config["capacity"].(int))),
	}

	if tier, ok := config["tier"].(string); ok && tier != "" {
		sku.Tier = &tier
	}

	return sku
}

func expandAzureRmRollingUpgradePolicy(d *pluginsdk.ResourceData) *compute.RollingUpgradePolicy {
	if config, ok := d.GetOk("rolling_upgrade_policy.0"); ok {
		policy := config.(map[string]interface{})
		return &compute.RollingUpgradePolicy{
			MaxBatchInstancePercent:             utils.Int32(int32(policy["max_batch_instance_percent"].(int))),
			MaxUnhealthyInstancePercent:         utils.Int32(int32(policy["max_unhealthy_instance_percent"].(int))),
			MaxUnhealthyUpgradedInstancePercent: utils.Int32(int32(policy["max_unhealthy_upgraded_instance_percent"].(int))),
			PauseTimeBetweenBatches:             utils.String(policy["pause_time_between_batches"].(string)),
		}
	}
	return nil
}

func expandAzureRmVirtualMachineScaleSetNetworkProfile(d *pluginsdk.ResourceData) *compute.VirtualMachineScaleSetNetworkProfile {
	scaleSetNetworkProfileConfigs := d.Get("network_profile").(*pluginsdk.Set).List()
	networkProfileConfig := make([]compute.VirtualMachineScaleSetNetworkConfiguration, 0, len(scaleSetNetworkProfileConfigs))

	for _, npProfileConfig := range scaleSetNetworkProfileConfigs {
		config := npProfileConfig.(map[string]interface{})

		name := config["name"].(string)
		primary := config["primary"].(bool)
		acceleratedNetworking := config["accelerated_networking"].(bool)
		ipForwarding := config["ip_forwarding"].(bool)

		dnsSettingsConfigs := config["dns_settings"].([]interface{})
		dnsSettings := compute.VirtualMachineScaleSetNetworkConfigurationDNSSettings{}
		for _, dnsSettingsConfig := range dnsSettingsConfigs {
			dns_settings := dnsSettingsConfig.(map[string]interface{})

			if v := dns_settings["dns_servers"]; v != nil {
				dns_servers := dns_settings["dns_servers"].([]interface{})
				if len(dns_servers) > 0 {
					var dnsServers []string
					for _, v := range dns_servers {
						str := v.(string)
						dnsServers = append(dnsServers, str)
					}
					dnsSettings.DNSServers = &dnsServers
				}
			}
		}
		ipConfigurationConfigs := config["ip_configuration"].([]interface{})
		ipConfigurations := make([]compute.VirtualMachineScaleSetIPConfiguration, 0, len(ipConfigurationConfigs))
		for _, ipConfigConfig := range ipConfigurationConfigs {
			ipconfig := ipConfigConfig.(map[string]interface{})
			name := ipconfig["name"].(string)
			primary := ipconfig["primary"].(bool)
			subnetId := ipconfig["subnet_id"].(string)

			ipConfiguration := compute.VirtualMachineScaleSetIPConfiguration{
				Name: &name,
				VirtualMachineScaleSetIPConfigurationProperties: &compute.VirtualMachineScaleSetIPConfigurationProperties{
					Subnet: &compute.APIEntityReference{
						ID: &subnetId,
					},
				},
			}

			ipConfiguration.Primary = &primary

			if v := ipconfig["application_gateway_backend_address_pool_ids"]; v != nil {
				pools := v.(*pluginsdk.Set).List()
				resources := make([]compute.SubResource, 0, len(pools))
				for _, p := range pools {
					id := p.(string)
					resources = append(resources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.ApplicationGatewayBackendAddressPools = &resources
			}

			if v := ipconfig["application_security_group_ids"]; v != nil {
				asgs := v.(*pluginsdk.Set).List()
				resources := make([]compute.SubResource, 0, len(asgs))
				for _, p := range asgs {
					id := p.(string)
					resources = append(resources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.ApplicationSecurityGroups = &resources
			}

			if v := ipconfig["load_balancer_backend_address_pool_ids"]; v != nil {
				pools := v.(*pluginsdk.Set).List()
				resources := make([]compute.SubResource, 0, len(pools))
				for _, p := range pools {
					id := p.(string)
					resources = append(resources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.LoadBalancerBackendAddressPools = &resources
			}

			if v := ipconfig["load_balancer_inbound_nat_rules_ids"]; v != nil {
				rules := v.(*pluginsdk.Set).List()
				rulesResources := make([]compute.SubResource, 0, len(rules))
				for _, m := range rules {
					id := m.(string)
					rulesResources = append(rulesResources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.LoadBalancerInboundNatPools = &rulesResources
			}

			if v := ipconfig["public_ip_address_configuration"]; v != nil {
				publicIpConfigs := v.([]interface{})
				for _, publicIpConfigConfig := range publicIpConfigs {
					publicIpConfig := publicIpConfigConfig.(map[string]interface{})

					domainNameLabel := publicIpConfig["domain_name_label"].(string)
					dnsSettings := compute.VirtualMachineScaleSetPublicIPAddressConfigurationDNSSettings{
						DomainNameLabel: &domainNameLabel,
					}

					idleTimeout := int32(publicIpConfig["idle_timeout"].(int))
					prop := compute.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
						DNSSettings:          &dnsSettings,
						IdleTimeoutInMinutes: &idleTimeout,
					}

					publicIPConfigName := publicIpConfig["name"].(string)
					config := compute.VirtualMachineScaleSetPublicIPAddressConfiguration{
						Name: &publicIPConfigName,
						VirtualMachineScaleSetPublicIPAddressConfigurationProperties: &prop,
					}
					ipConfiguration.PublicIPAddressConfiguration = &config
				}
			}

			ipConfigurations = append(ipConfigurations, ipConfiguration)
		}

		nProfile := compute.VirtualMachineScaleSetNetworkConfiguration{
			Name: &name,
			VirtualMachineScaleSetNetworkConfigurationProperties: &compute.VirtualMachineScaleSetNetworkConfigurationProperties{
				Primary:                     &primary,
				IPConfigurations:            &ipConfigurations,
				EnableAcceleratedNetworking: &acceleratedNetworking,
				EnableIPForwarding:          &ipForwarding,
				DNSSettings:                 &dnsSettings,
			},
		}

		if v := config["network_security_group_id"].(string); v != "" {
			networkSecurityGroupId := compute.SubResource{
				ID: &v,
			}
			nProfile.VirtualMachineScaleSetNetworkConfigurationProperties.NetworkSecurityGroup = &networkSecurityGroupId
		}

		networkProfileConfig = append(networkProfileConfig, nProfile)
	}

	return &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: &networkProfileConfig,
	}
}

func expandAzureRMVirtualMachineScaleSetsOsProfile(d *pluginsdk.ResourceData) *compute.VirtualMachineScaleSetOSProfile {
	osProfileConfigs := d.Get("os_profile").([]interface{})

	osProfileConfig := osProfileConfigs[0].(map[string]interface{})
	namePrefix := osProfileConfig["computer_name_prefix"].(string)
	username := osProfileConfig["admin_username"].(string)
	password := osProfileConfig["admin_password"].(string)
	customData := osProfileConfig["custom_data"].(string)

	osProfile := &compute.VirtualMachineScaleSetOSProfile{
		ComputerNamePrefix: &namePrefix,
		AdminUsername:      &username,
	}

	if password != "" {
		osProfile.AdminPassword = &password
	}

	if customData != "" {
		customData = utils.Base64EncodeIfNot(customData)
		osProfile.CustomData = &customData
	}

	if _, ok := d.GetOk("os_profile_secrets"); ok {
		secrets := expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d)
		if secrets != nil {
			osProfile.Secrets = secrets
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		osProfile.LinuxConfiguration = expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d)
	}

	if _, ok := d.GetOk("os_profile_windows_config"); ok {
		winConfig := expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d)
		if winConfig != nil {
			osProfile.WindowsConfiguration = winConfig
		}
	}

	return osProfile
}

func expandAzureRMVirtualMachineScaleSetsDiagnosticProfile(d *pluginsdk.ResourceData) compute.DiagnosticsProfile {
	bootDiagnosticConfigs := d.Get("boot_diagnostics").([]interface{})
	bootDiagnosticConfig := bootDiagnosticConfigs[0].(map[string]interface{})

	enabled := bootDiagnosticConfig["enabled"].(bool)
	storageURI := bootDiagnosticConfig["storage_uri"].(string)

	bootDiagnostic := &compute.BootDiagnostics{
		Enabled:    &enabled,
		StorageURI: &storageURI,
	}

	diagnosticsProfile := compute.DiagnosticsProfile{
		BootDiagnostics: bootDiagnostic,
	}

	return diagnosticsProfile
}

func expandAzureRmVirtualMachineScaleSetIdentity(d *pluginsdk.ResourceData) *compute.VirtualMachineScaleSetIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := compute.ResourceIdentityType(identity["type"].(string))

	identityIds := make(map[string]*compute.UserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &compute.UserAssignedIdentitiesValue{}
	}

	vmssIdentity := compute.VirtualMachineScaleSetIdentity{
		Type: identityType,
	}

	if vmssIdentity.Type == compute.ResourceIdentityTypeUserAssigned || vmssIdentity.Type == compute.ResourceIdentityTypeSystemAssignedUserAssigned {
		vmssIdentity.UserAssignedIdentities = identityIds
	}

	return &vmssIdentity
}

func expandAzureRMVirtualMachineScaleSetsStorageProfileOsDisk(d *pluginsdk.ResourceData) (*compute.VirtualMachineScaleSetOSDisk, error) {
	osDiskConfigs := d.Get("storage_profile_os_disk").(*pluginsdk.Set).List()

	osDiskConfig := osDiskConfigs[0].(map[string]interface{})
	name := osDiskConfig["name"].(string)
	image := osDiskConfig["image"].(string)
	vhd_containers := osDiskConfig["vhd_containers"].(*pluginsdk.Set).List()
	caching := osDiskConfig["caching"].(string)
	osType := osDiskConfig["os_type"].(string)
	createOption := osDiskConfig["create_option"].(string)
	managedDiskType := osDiskConfig["managed_disk_type"].(string)

	if managedDiskType == "" && name == "" {
		return nil, fmt.Errorf("[ERROR] `name` must be set in `storage_profile_os_disk` for unmanaged disk")
	}

	osDisk := &compute.VirtualMachineScaleSetOSDisk{
		Name:         &name,
		Caching:      compute.CachingTypes(caching),
		OsType:       compute.OperatingSystemTypes(osType),
		CreateOption: compute.DiskCreateOptionTypes(createOption),
	}

	if image != "" {
		osDisk.Image = &compute.VirtualHardDisk{
			URI: &image,
		}
	}

	if len(vhd_containers) > 0 {
		var vhdContainers []string
		for _, v := range vhd_containers {
			str := v.(string)
			vhdContainers = append(vhdContainers, str)
		}
		osDisk.VhdContainers = &vhdContainers
	}

	managedDisk := &compute.VirtualMachineScaleSetManagedDiskParameters{}

	if managedDiskType != "" {
		if name != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `name` and `managed_disk_type` on `storage_profile_os_disk` (please remove name or set it to blank)")
		}

		osDisk.Name = nil
		managedDisk.StorageAccountType = compute.StorageAccountTypes(managedDiskType)
		osDisk.ManagedDisk = managedDisk
	}

	// BEGIN: code to be removed after GH-13016 is merged
	if image != "" && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `image` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}

	if len(vhd_containers) > 0 && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_containers` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}
	// END: code to be removed after GH-13016 is merged

	return osDisk, nil
}

func expandAzureRMVirtualMachineScaleSetsStorageProfileDataDisk(d *pluginsdk.ResourceData) *[]compute.VirtualMachineScaleSetDataDisk {
	disks := d.Get("storage_profile_data_disk").([]interface{})
	dataDisks := make([]compute.VirtualMachineScaleSetDataDisk, 0, len(disks))
	for _, diskConfig := range disks {
		config := diskConfig.(map[string]interface{})

		createOption := config["create_option"].(string)
		managedDiskType := config["managed_disk_type"].(string)
		lun := int32(config["lun"].(int))

		dataDisk := compute.VirtualMachineScaleSetDataDisk{
			Lun:          &lun,
			CreateOption: compute.DiskCreateOptionTypes(createOption),
		}

		managedDiskVMSS := &compute.VirtualMachineScaleSetManagedDiskParameters{}

		if managedDiskType != "" {
			managedDiskVMSS.StorageAccountType = compute.StorageAccountTypes(managedDiskType)
		} else {
			managedDiskVMSS.StorageAccountType = compute.StorageAccountTypes(compute.StorageAccountTypeStandardLRS)
		}

		// assume that data disks in VMSS can only be Managed Disks
		dataDisk.ManagedDisk = managedDiskVMSS
		if v := config["caching"].(string); v != "" {
			dataDisk.Caching = compute.CachingTypes(v)
		}

		if v := config["disk_size_gb"]; v != nil {
			diskSize := int32(config["disk_size_gb"].(int))
			dataDisk.DiskSizeGB = &diskSize
		}

		dataDisks = append(dataDisks, dataDisk)
	}

	return &dataDisks
}

func expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d *pluginsdk.ResourceData) (*compute.ImageReference, error) {
	storageImageRefs := d.Get("storage_profile_image_reference").(*pluginsdk.Set).List()

	storageImageRef := storageImageRefs[0].(map[string]interface{})

	imageID := storageImageRef["id"].(string)
	publisher := storageImageRef["publisher"].(string)

	imageReference := compute.ImageReference{}

	if imageID != "" && publisher != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `id` and `publisher` (only one or the other can be used)")
	}

	if imageID != "" {
		imageReference.ID = utils.String(storageImageRef["id"].(string))
	} else {
		offer := storageImageRef["offer"].(string)
		sku := storageImageRef["sku"].(string)
		version := storageImageRef["version"].(string)

		imageReference.Publisher = utils.String(publisher)
		imageReference.Offer = utils.String(offer)
		imageReference.Sku = utils.String(sku)
		imageReference.Version = utils.String(version)
	}

	return &imageReference, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d *pluginsdk.ResourceData) *compute.LinuxConfiguration {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*pluginsdk.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]interface{})
	disablePasswordAuth := linuxConfig["disable_password_authentication"].(bool)

	linuxKeys := linuxConfig["ssh_keys"].([]interface{})
	sshPublicKeys := make([]compute.SSHPublicKey, 0, len(linuxKeys))
	for _, key := range linuxKeys {
		if key == nil {
			continue
		}
		sshKey := key.(map[string]interface{})
		path := sshKey["path"].(string)
		keyData := sshKey["key_data"].(string)

		sshPublicKey := compute.SSHPublicKey{
			Path:    &path,
			KeyData: &keyData,
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	config := &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &disablePasswordAuth,
		SSH: &compute.SSHConfiguration{
			PublicKeys: &sshPublicKeys,
		},
	}

	return config
}

func expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d *pluginsdk.ResourceData) *compute.WindowsConfiguration {
	osProfilesWindowsConfig := d.Get("os_profile_windows_config").(*pluginsdk.Set).List()

	osProfileConfig := osProfilesWindowsConfig[0].(map[string]interface{})
	config := &compute.WindowsConfiguration{}

	if v := osProfileConfig["provision_vm_agent"]; v != nil {
		provision := v.(bool)
		config.ProvisionVMAgent = &provision
	}

	if v := osProfileConfig["enable_automatic_upgrades"]; v != nil {
		update := v.(bool)
		config.EnableAutomaticUpdates = &update
	}

	if v := osProfileConfig["winrm"]; v != nil {
		winRm := v.([]interface{})
		if len(winRm) > 0 {
			winRmListeners := make([]compute.WinRMListener, 0, len(winRm))
			for _, winRmConfig := range winRm {
				config := winRmConfig.(map[string]interface{})

				protocol := config["protocol"].(string)
				winRmListener := compute.WinRMListener{
					Protocol: compute.ProtocolTypes(protocol),
				}
				if v := config["certificate_url"].(string); v != "" {
					winRmListener.CertificateURL = &v
				}

				winRmListeners = append(winRmListeners, winRmListener)
			}
			config.WinRM = &compute.WinRMConfiguration{
				Listeners: &winRmListeners,
			}
		}
	}
	if v := osProfileConfig["additional_unattend_config"]; v != nil {
		additionalConfig := v.([]interface{})
		if len(additionalConfig) > 0 {
			additionalConfigContent := make([]compute.AdditionalUnattendContent, 0, len(additionalConfig))
			for _, addConfig := range additionalConfig {
				config := addConfig.(map[string]interface{})
				pass := config["pass"].(string)
				component := config["component"].(string)
				settingName := config["setting_name"].(string)
				content := config["content"].(string)

				addContent := compute.AdditionalUnattendContent{
					PassName:      compute.PassNames(pass),
					ComponentName: compute.ComponentNames(component),
					SettingName:   compute.SettingNames(settingName),
				}

				if content != "" {
					addContent.Content = &content
				}

				additionalConfigContent = append(additionalConfigContent, addContent)
			}
			config.AdditionalUnattendContent = &additionalConfigContent
		}
	}
	return config
}

func expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d *pluginsdk.ResourceData) *[]compute.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").(*pluginsdk.Set).List()
	secrets := make([]compute.VaultSecretGroup, 0, len(secretsConfig))

	for _, secretConfig := range secretsConfig {
		config := secretConfig.(map[string]interface{})
		sourceVaultId := config["source_vault_id"].(string)

		vaultSecretGroup := compute.VaultSecretGroup{
			SourceVault: &compute.SubResource{
				ID: &sourceVaultId,
			},
		}

		if v := config["vault_certificates"]; v != nil {
			certsConfig := v.([]interface{})
			certs := make([]compute.VaultCertificate, 0, len(certsConfig))
			for _, certConfig := range certsConfig {
				config := certConfig.(map[string]interface{})

				certUrl := config["certificate_url"].(string)
				cert := compute.VaultCertificate{
					CertificateURL: &certUrl,
				}
				if v := config["certificate_store"].(string); v != "" {
					cert.CertificateStore = &v
				}

				certs = append(certs, cert)
			}
			vaultSecretGroup.VaultCertificates = &certs
		}

		secrets = append(secrets, vaultSecretGroup)
	}

	return &secrets
}

func expandAzureRMVirtualMachineScaleSetExtensions(d *pluginsdk.ResourceData) (*compute.VirtualMachineScaleSetExtensionProfile, error) {
	extensions := d.Get("extension").(*pluginsdk.Set).List()
	resources := make([]compute.VirtualMachineScaleSetExtension, 0, len(extensions))
	for _, e := range extensions {
		config := e.(map[string]interface{})
		name := config["name"].(string)
		publisher := config["publisher"].(string)
		t := config["type"].(string)
		version := config["type_handler_version"].(string)

		extension := compute.VirtualMachineScaleSetExtension{
			Name: &name,
			VirtualMachineScaleSetExtensionProperties: &compute.VirtualMachineScaleSetExtensionProperties{
				Publisher:          &publisher,
				Type:               &t,
				TypeHandlerVersion: &version,
			},
		}

		if u := config["auto_upgrade_minor_version"]; u != nil {
			upgrade := u.(bool)
			extension.VirtualMachineScaleSetExtensionProperties.AutoUpgradeMinorVersion = &upgrade
		}

		if a := config["provision_after_extensions"]; a != nil {
			provision_after_extensions := config["provision_after_extensions"].(*pluginsdk.Set).List()
			if len(provision_after_extensions) > 0 {
				var provisionAfterExtensions []string
				for _, a := range provision_after_extensions {
					str := a.(string)
					provisionAfterExtensions = append(provisionAfterExtensions, str)
				}
				extension.VirtualMachineScaleSetExtensionProperties.ProvisionAfterExtensions = &provisionAfterExtensions
			}
		}

		if s := config["settings"].(string); s != "" {
			settings, err := pluginsdk.ExpandJsonFromString(s)
			if err != nil {
				return nil, fmt.Errorf("unable to parse settings: %+v", err)
			}
			extension.VirtualMachineScaleSetExtensionProperties.Settings = &settings
		}

		if s := config["protected_settings"].(string); s != "" {
			protectedSettings, err := pluginsdk.ExpandJsonFromString(s)
			if err != nil {
				return nil, fmt.Errorf("unable to parse protected_settings: %+v", err)
			}
			extension.VirtualMachineScaleSetExtensionProperties.ProtectedSettings = &protectedSettings
		}

		resources = append(resources, extension)
	}

	return &compute.VirtualMachineScaleSetExtensionProfile{
		Extensions: &resources,
	}, nil
}

func expandAzureRmVirtualMachineScaleSetPlan(d *pluginsdk.ResourceData) *compute.Plan {
	planConfigs := d.Get("plan").(*pluginsdk.Set).List()

	planConfig := planConfigs[0].(map[string]interface{})

	publisher := planConfig["publisher"].(string)
	name := planConfig["name"].(string)
	product := planConfig["product"].(string)

	return &compute.Plan{
		Publisher: &publisher,
		Name:      &name,
		Product:   &product,
	}
}

func flattenAzureRmVirtualMachineScaleSetPlan(plan *compute.Plan) []interface{} {
	result := make(map[string]interface{})

	result["name"] = *plan.Name
	result["publisher"] = *plan.Publisher
	result["product"] = *plan.Product

	return []interface{}{result}
}

// When upgrade_policy_mode is not Rolling, we will just ignore rolling_upgrade_policy (returns true).
func azureRmVirtualMachineScaleSetSuppressRollingUpgradePolicyDiff(k, _, new string, d *pluginsdk.ResourceData) bool {
	if k == "rolling_upgrade_policy.#" && new == "0" {
		return strings.ToLower(d.Get("upgrade_policy_mode").(string)) != "rolling"
	}
	return false
}

// Make sure rolling_upgrade_policy is default value when upgrade_policy_mode is not Rolling.
func azureRmVirtualMachineScaleSetCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	mode := d.Get("upgrade_policy_mode").(string)
	if strings.ToLower(mode) != "rolling" {
		if policyRaw, ok := d.GetOk("rolling_upgrade_policy.0"); ok {
			policy := policyRaw.(map[string]interface{})
			isDefault := (policy["max_batch_instance_percent"].(int) == 20) &&
				(policy["max_unhealthy_instance_percent"].(int) == 20) &&
				(policy["max_unhealthy_upgraded_instance_percent"].(int) == 20) &&
				(policy["pause_time_between_batches"] == "PT0S")
			if !isDefault {
				return fmt.Errorf("If `upgrade_policy_mode` is `%s`, `rolling_upgrade_policy` must be removed or set to default values", mode)
			}
		}
	}
	return nil
}
