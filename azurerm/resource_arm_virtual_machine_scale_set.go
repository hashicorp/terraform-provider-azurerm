package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualMachineScaleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineScaleSetCreate,
		Read:   resourceArmVirtualMachineScaleSetRead,
		Update: resourceArmVirtualMachineScaleSetCreate,
		Delete: resourceArmVirtualMachineScaleSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"zones": zonesSchema(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, true),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sku": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"tier": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetSkuHash,
			},

			"upgrade_policy_mode": {
				Type:     schema.TypeString,
				Required: true,
			},

			"overprovision": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"single_placement_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"os_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"computer_name_prefix": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
						},

						"admin_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"custom_data": {
							Type:             schema.TypeString,
							Optional:         true,
							StateFunc:        userDataStateFunc,
							DiffSuppressFunc: userDataDiffSuppressFunc,
						},
					},
				},
			},

			"os_profile_secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_vault_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vault_certificates": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"certificate_store": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"os_profile_windows_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provision_vm_agent": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_automatic_upgrades": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"winrm": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
									"certificate_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"additional_unattend_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pass": {
										Type:     schema.TypeString,
										Required: true,
									},
									"component": {
										Type:     schema.TypeString,
										Required: true,
									},
									"setting_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetOsProfileWindowsConfigHash,
			},

			"os_profile_linux_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_password_authentication": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"ssh_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"key_data": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash,
			},

			"network_profile": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"primary": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"accelerated_networking": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"network_security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"ip_configuration": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"subnet_id": {
										Type:     schema.TypeString,
										Required: true,
									},

									"application_gateway_backend_address_pool_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"load_balancer_backend_address_pool_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"load_balancer_inbound_nat_rules_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"primary": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"public_ip_address_configuration": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},

												"idle_timeout": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(4, 32),
												},

												"domain_name_label": {
													Type:     schema.TypeString,
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
				Set: resourceArmVirtualMachineScaleSetNetworkConfigurationHash,
			},

			"boot_diagnostics": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"storage_uri": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"storage_profile_os_disk": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"vhd_containers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"managed_disk_type": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_profile_os_disk.vhd_containers"},
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.PremiumLRS),
								string(compute.StandardLRS),
							}, true),
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"create_option": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash,
			},

			"storage_profile_data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lun": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"create_option": {
							Type:     schema.TypeString,
							Required: true,
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"disk_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateDiskSizeGB,
						},

						"managed_disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.PremiumLRS),
								string(compute.StandardLRS),
							}, true),
						},
					},
				},
			},

			"storage_profile_image_reference": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"offer": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"sku": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash,
			},

			"plan": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},

						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"extension": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"type_handler_version": {
							Type:     schema.TypeString,
							Required: true,
						},

						"auto_upgrade_minor_version": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"settings": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.ValidateJsonString,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},

						"protected_settings": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateFunc:     validation.ValidateJsonString,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetExtensionHash,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualMachineScaleSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine Scale Set creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	zones := expandZones(d.Get("zones").([]interface{}))

	sku, err := expandVirtualMachineScaleSetSku(d)
	if err != nil {
		return err
	}

	storageProfile := compute.VirtualMachineScaleSetStorageProfile{}
	osDisk, err := expandAzureRMVirtualMachineScaleSetsStorageProfileOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile.OsDisk = osDisk

	if _, ok := d.GetOk("storage_profile_data_disk"); ok {
		dataDisks, err := expandAzureRMVirtualMachineScaleSetsStorageProfileDataDisk(d)
		if err != nil {
			return err
		}
		storageProfile.DataDisks = &dataDisks
	}

	if _, ok := d.GetOk("storage_profile_image_reference"); ok {
		imageRef, err := expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d)
		if err != nil {
			return err
		}
		storageProfile.ImageReference = imageRef
	}

	osProfile, err := expandAzureRMVirtualMachineScaleSetsOsProfile(d)
	if err != nil {
		return err
	}

	extensions, err := expandAzureRMVirtualMachineScaleSetExtensions(d)
	if err != nil {
		return err
	}

	updatePolicy := d.Get("upgrade_policy_mode").(string)
	overprovision := d.Get("overprovision").(bool)
	singlePlacementGroup := d.Get("single_placement_group").(bool)

	scaleSetProps := compute.VirtualMachineScaleSetProperties{
		UpgradePolicy: &compute.UpgradePolicy{
			Mode: compute.UpgradeMode(updatePolicy),
		},
		VirtualMachineProfile: &compute.VirtualMachineScaleSetVMProfile{
			NetworkProfile:   expandAzureRmVirtualMachineScaleSetNetworkProfile(d),
			StorageProfile:   &storageProfile,
			OsProfile:        osProfile,
			ExtensionProfile: extensions,
		},
		Overprovision:        &overprovision,
		SinglePlacementGroup: &singlePlacementGroup,
	}

	if _, ok := d.GetOk("boot_diagnostics"); ok {
		diagnosticProfile := expandAzureRMVirtualMachineScaleSetsDiagnosticProfile(d)
		scaleSetProps.VirtualMachineProfile.DiagnosticsProfile = &diagnosticProfile
	}

	scaleSetParams := compute.VirtualMachineScaleSet{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		Sku:      sku,
		VirtualMachineScaleSetProperties: &scaleSetProps,
		Zones: zones,
	}

	if _, ok := d.GetOk("identity"); ok {
		scaleSetParams.Identity = expandAzureRmVirtualMachineScaleSetIdentity(d)
	}

	if _, ok := d.GetOk("plan"); ok {
		plan, err := expandAzureRmVirtualMachineScaleSetPlan(d)
		if err != nil {
			return err
		}

		scaleSetParams.Plan = plan
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, scaleSetParams)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Machine Scale Set %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualMachineScaleSetRead(d, meta)
}

func resourceArmVirtualMachineScaleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachineScaleSets"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] AzureRM Virtual Machine Scale Set (%s) Not Found. Removing from State", name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Virtual Machine Scale Set %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("zones", resp.Zones)

	if err := d.Set("sku", flattenAzureRmVirtualMachineScaleSetSku(resp.Sku)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Sku error: %#v", err)
	}

	if err := d.Set("identity", flattenAzureRmVirtualMachineScaleSetIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("[DEBUG] Error flattening `identity`: %+v", err)
	}

	properties := resp.VirtualMachineScaleSetProperties

	d.Set("upgrade_policy_mode", properties.UpgradePolicy.Mode)
	d.Set("overprovision", properties.Overprovision)
	d.Set("single_placement_group", properties.SinglePlacementGroup)

	osProfile, err := flattenAzureRMVirtualMachineScaleSetOsProfile(d, properties.VirtualMachineProfile.OsProfile)
	if err != nil {
		return fmt.Errorf("[DEBUG] Error flattening Virtual Machine Scale Set OS Profile. Error: %#v", err)
	}

	if err := d.Set("os_profile", osProfile); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set OS Profile error: %#v", err)
	}

	if properties.VirtualMachineProfile.OsProfile.Secrets != nil {
		if err := d.Set("os_profile_secrets", flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(properties.VirtualMachineProfile.OsProfile.Secrets)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set OS Profile Secrets error: %#v", err)
		}
	}

	if properties.VirtualMachineProfile.DiagnosticsProfile != nil {
		if err := d.Set("boot_diagnostics", flattenAzureRmVirtualMachineScaleSetBootDiagnostics(properties.VirtualMachineProfile.DiagnosticsProfile.BootDiagnostics)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virutal Machine Scale Set Boot Diagnostics: %#v", err)
		}
	}

	if properties.VirtualMachineProfile.OsProfile.WindowsConfiguration != nil {
		if err := d.Set("os_profile_windows_config", flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(properties.VirtualMachineProfile.OsProfile.WindowsConfiguration)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set OS Profile Windows config error: %#v", err)
		}
	}

	if properties.VirtualMachineProfile.OsProfile.LinuxConfiguration != nil {
		if err := d.Set("os_profile_linux_config", flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(properties.VirtualMachineProfile.OsProfile.LinuxConfiguration)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set OS Profile Linux config error: %#v", err)
		}
	}

	if err := d.Set("network_profile", flattenAzureRmVirtualMachineScaleSetNetworkProfile(properties.VirtualMachineProfile.NetworkProfile)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Network Profile error: %#v", err)
	}

	if properties.VirtualMachineProfile.StorageProfile.ImageReference != nil {
		if err := d.Set("storage_profile_image_reference", flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(properties.VirtualMachineProfile.StorageProfile.ImageReference)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Storage Profile Image Reference error: %#v", err)
		}
	}

	if err := d.Set("storage_profile_os_disk", flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(properties.VirtualMachineProfile.StorageProfile.OsDisk)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Storage Profile OS Disk error: %#v", err)
	}

	if resp.VirtualMachineProfile.StorageProfile.DataDisks != nil {
		if err := d.Set("storage_profile_data_disk", flattenAzureRmVirtualMachineScaleSetStorageProfileDataDisk(properties.VirtualMachineProfile.StorageProfile.DataDisks)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Storage Profile Data Disk error: %#v", err)
		}
	}

	if properties.VirtualMachineProfile.ExtensionProfile != nil {
		extension, err := flattenAzureRmVirtualMachineScaleSetExtensionProfile(properties.VirtualMachineProfile.ExtensionProfile)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Extension Profile error: %#v", err)
		}
		d.Set("extension", extension)
	}

	if resp.Plan != nil {
		if err := d.Set("plan", flattenAzureRmVirtualMachineScaleSetPlan(resp.Plan)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Plan. Error: %#v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualMachineScaleSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachineScaleSets"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}

func flattenAzureRmVirtualMachineScaleSetIdentity(identity *compute.VirtualMachineScaleSetIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(config *compute.LinuxConfiguration) []interface{} {
	result := make(map[string]interface{})
	result["disable_password_authentication"] = *config.DisablePasswordAuthentication

	if config.SSH != nil && len(*config.SSH.PublicKeys) > 0 {
		ssh_keys := make([]map[string]interface{}, 0, len(*config.SSH.PublicKeys))
		for _, i := range *config.SSH.PublicKeys {
			key := make(map[string]interface{})
			key["path"] = *i.Path

			if i.KeyData != nil {
				key["key_data"] = *i.KeyData
			}

			ssh_keys = append(ssh_keys, key)
		}

		result["ssh_keys"] = ssh_keys
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
	b := map[string]interface{}{
		"enabled":     *bootDiagnostic.Enabled,
		"storage_uri": *bootDiagnostic.StorageURI,
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

		if v := netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.NetworkSecurityGroup; v != nil {
			s["network_security_group_id"] = *v.ID
		}

		if netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations != nil {
			ipConfigs := make([]map[string]interface{}, 0, len(*netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations))
			for _, ipConfig := range *netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations {
				config := make(map[string]interface{})
				config["name"] = *ipConfig.Name

				properties := ipConfig.VirtualMachineScaleSetIPConfigurationProperties

				if ipConfig.VirtualMachineScaleSetIPConfigurationProperties.Subnet != nil {
					config["subnet_id"] = *properties.Subnet.ID
				}

				addressPools := make([]interface{}, 0)
				if properties.ApplicationGatewayBackendAddressPools != nil {
					for _, pool := range *properties.ApplicationGatewayBackendAddressPools {
						addressPools = append(addressPools, *pool.ID)
					}
				}
				config["application_gateway_backend_address_pool_ids"] = schema.NewSet(schema.HashString, addressPools)

				if properties.LoadBalancerBackendAddressPools != nil {
					addressPools := make([]interface{}, 0, len(*properties.LoadBalancerBackendAddressPools))
					for _, pool := range *properties.LoadBalancerBackendAddressPools {
						addressPools = append(addressPools, *pool.ID)
					}
					config["load_balancer_backend_address_pool_ids"] = schema.NewSet(schema.HashString, addressPools)
				}

				if properties.LoadBalancerInboundNatPools != nil {
					inboundNatPools := make([]interface{}, 0, len(*properties.LoadBalancerInboundNatPools))
					for _, rule := range *properties.LoadBalancerInboundNatPools {
						inboundNatPools = append(inboundNatPools, *rule.ID)
					}
					config["load_balancer_inbound_nat_rules_ids"] = schema.NewSet(schema.HashString, inboundNatPools)
				}

				if properties.Primary != nil {
					config["primary"] = *properties.Primary
				}

				if properties.PublicIPAddressConfiguration != nil {
					publicIpInfo := properties.PublicIPAddressConfiguration
					publicIpConfigs := make([]map[string]interface{}, 0, 1)
					publicIpConfig := make(map[string]interface{})
					publicIpConfig["name"] = *publicIpInfo.Name
					publicIpConfig["domain_name_label"] = *publicIpInfo.VirtualMachineScaleSetPublicIPAddressConfigurationProperties.DNSSettings
					publicIpConfig["idle_timeout"] = *publicIpInfo.VirtualMachineScaleSetPublicIPAddressConfigurationProperties.IdleTimeoutInMinutes
					config["public_ip_address_configuration"] = publicIpConfigs
				}

				ipConfigs = append(ipConfigs, config)
			}

			s["ip_configuration"] = ipConfigs
		}

		result = append(result, s)
	}

	return result
}

func flattenAzureRMVirtualMachineScaleSetOsProfile(d *schema.ResourceData, profile *compute.VirtualMachineScaleSetOSProfile) ([]interface{}, error) {
	result := make(map[string]interface{})

	result["computer_name_prefix"] = *profile.ComputerNamePrefix
	result["admin_username"] = *profile.AdminUsername

	// admin password isn't returned, so let's look it up
	if v, ok := d.GetOk("os_profile.0.admin_password"); ok {
		password := v.(string)
		result["admin_password"] = password
	}

	if profile.CustomData != nil {
		result["custom_data"] = *profile.CustomData
	} else {
		// look up the current custom data
		value := d.Get("os_profile.0.custom_data").(string)
		if !isBase64Encoded(value) {
			value = base64Encode(value)
		}
		result["custom_data"] = value
	}

	return []interface{}{result}, nil
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(profile *compute.VirtualMachineScaleSetOSDisk) []interface{} {
	result := make(map[string]interface{})

	if profile.Name != nil {
		result["name"] = *profile.Name
	}

	if profile.Image != nil {
		result["image"] = *profile.Image.URI
	}

	if profile.VhdContainers != nil {
		containers := make([]interface{}, 0, len(*profile.VhdContainers))
		for _, container := range *profile.VhdContainers {
			containers = append(containers, container)
		}
		result["vhd_containers"] = schema.NewSet(schema.HashString, containers)
	}

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
		l["lun"] = *disk.Lun

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
	result["name"] = *sku.Name
	result["capacity"] = *sku.Capacity

	if *sku.Tier != "" {
		result["tier"] = *sku.Tier
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

			if properties.Settings != nil {
				settings, err := structure.FlattenJsonToString(*properties.Settings)
				if err != nil {
					return nil, err
				}
				e["settings"] = settings
			}
		}

		result = append(result, e)
	}

	return result, nil
}

func resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["publisher"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
	}
	if m["offer"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["offer"].(string)))
	}
	if m["sku"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["sku"].(string)))
	}
	if m["version"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["version"].(string)))
	}
	if m["id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["id"].(string)))
	}
	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetSkuHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	if m["tier"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["tier"].(string))))
	}
	buf.WriteString(fmt.Sprintf("%d-", m["capacity"].(int)))

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))

	if m["vhd_containers"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["vhd_containers"].(*schema.Set).List()))
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetNetworkConfigurationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%t-", m["primary"].(bool)))
	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetOsProfileWindowsConfigHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["provision_vm_agent"] != nil {
		buf.WriteString(fmt.Sprintf("%t-", m["provision_vm_agent"].(bool)))
	}
	if m["enable_automatic_upgrades"] != nil {
		buf.WriteString(fmt.Sprintf("%t-", m["enable_automatic_upgrades"].(bool)))
	}
	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetExtensionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["type_handler_version"].(string)))
	if m["auto_upgrade_minor_version"] != nil {
		buf.WriteString(fmt.Sprintf("%t-", m["auto_upgrade_minor_version"].(bool)))
	}

	// we need to ensure the whitespace is consistent
	settings := m["settings"].(string)
	if settings != "" {
		expandedSettings, err := structure.ExpandJsonFromString(settings)
		if err == nil {
			serialisedSettings, err := structure.FlattenJsonToString(expandedSettings)
			if err == nil {
				buf.WriteString(fmt.Sprintf("%s-", serialisedSettings))
			}
		}
	}

	return hashcode.String(buf.String())
}

func expandVirtualMachineScaleSetSku(d *schema.ResourceData) (*compute.Sku, error) {
	skuConfig := d.Get("sku").(*schema.Set).List()

	config := skuConfig[0].(map[string]interface{})

	name := config["name"].(string)
	tier := config["tier"].(string)
	capacity := int64(config["capacity"].(int))

	sku := &compute.Sku{
		Name:     &name,
		Capacity: &capacity,
	}

	if tier != "" {
		sku.Tier = &tier
	}

	return sku, nil
}

func expandAzureRmVirtualMachineScaleSetNetworkProfile(d *schema.ResourceData) *compute.VirtualMachineScaleSetNetworkProfile {
	scaleSetNetworkProfileConfigs := d.Get("network_profile").(*schema.Set).List()
	networkProfileConfig := make([]compute.VirtualMachineScaleSetNetworkConfiguration, 0, len(scaleSetNetworkProfileConfigs))

	for _, npProfileConfig := range scaleSetNetworkProfileConfigs {
		config := npProfileConfig.(map[string]interface{})

		name := config["name"].(string)
		primary := config["primary"].(bool)
		acceleratedNetworking := config["accelerated_networking"].(bool)

		ipConfigurationConfigs := config["ip_configuration"].([]interface{})
		ipConfigurations := make([]compute.VirtualMachineScaleSetIPConfiguration, 0, len(ipConfigurationConfigs))
		for _, ipConfigConfig := range ipConfigurationConfigs {
			ipconfig := ipConfigConfig.(map[string]interface{})
			name := ipconfig["name"].(string)
			subnetId := ipconfig["subnet_id"].(string)

			ipConfiguration := compute.VirtualMachineScaleSetIPConfiguration{
				Name: &name,
				VirtualMachineScaleSetIPConfigurationProperties: &compute.VirtualMachineScaleSetIPConfigurationProperties{
					Subnet: &compute.APIEntityReference{
						ID: &subnetId,
					},
				},
			}

			if v := ipconfig["application_gateway_backend_address_pool_ids"]; v != nil {
				pools := v.(*schema.Set).List()
				resources := make([]compute.SubResource, 0, len(pools))
				for _, p := range pools {
					id := p.(string)
					resources = append(resources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.ApplicationGatewayBackendAddressPools = &resources
			}

			if v := ipconfig["load_balancer_backend_address_pool_ids"]; v != nil {
				pools := v.(*schema.Set).List()
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
				rules := v.(*schema.Set).List()
				rulesResources := make([]compute.SubResource, 0, len(rules))
				for _, m := range rules {
					id := m.(string)
					rulesResources = append(rulesResources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.LoadBalancerInboundNatPools = &rulesResources
			}

			if v := ipconfig["primary"]; v != nil {
				primary := v.(bool)
				ipConfiguration.Primary = &primary
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

func expandAzureRMVirtualMachineScaleSetsOsProfile(d *schema.ResourceData) (*compute.VirtualMachineScaleSetOSProfile, error) {
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
		customData = base64Encode(customData)
		osProfile.CustomData = &customData
	}

	if _, ok := d.GetOk("os_profile_secrets"); ok {
		secrets := expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d)
		if secrets != nil {
			osProfile.Secrets = secrets
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		linuxConfig, err := expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d)
		if err != nil {
			return nil, err
		}
		osProfile.LinuxConfiguration = linuxConfig
	}

	if _, ok := d.GetOk("os_profile_windows_config"); ok {
		winConfig, err := expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d)
		if err != nil {
			return nil, err
		}
		if winConfig != nil {
			osProfile.WindowsConfiguration = winConfig
		}
	}

	return osProfile, nil
}

func expandAzureRMVirtualMachineScaleSetsDiagnosticProfile(d *schema.ResourceData) compute.DiagnosticsProfile {
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

func expandAzureRmVirtualMachineScaleSetIdentity(d *schema.ResourceData) *compute.VirtualMachineScaleSetIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := identity["type"].(string)
	return &compute.VirtualMachineScaleSetIdentity{
		Type: compute.ResourceIdentityType(identityType),
	}
}

func expandAzureRMVirtualMachineScaleSetsStorageProfileOsDisk(d *schema.ResourceData) (*compute.VirtualMachineScaleSetOSDisk, error) {
	osDiskConfigs := d.Get("storage_profile_os_disk").(*schema.Set).List()

	osDiskConfig := osDiskConfigs[0].(map[string]interface{})
	name := osDiskConfig["name"].(string)
	image := osDiskConfig["image"].(string)
	vhd_containers := osDiskConfig["vhd_containers"].(*schema.Set).List()
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

	//BEGIN: code to be removed after GH-13016 is merged
	if image != "" && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `image` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}

	if len(vhd_containers) > 0 && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_containers` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}
	//END: code to be removed after GH-13016 is merged

	return osDisk, nil
}

func expandAzureRMVirtualMachineScaleSetsStorageProfileDataDisk(d *schema.ResourceData) ([]compute.VirtualMachineScaleSetDataDisk, error) {
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
			managedDiskVMSS.StorageAccountType = compute.StorageAccountTypes(compute.StandardLRS)
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

	return dataDisks, nil
}

func expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d *schema.ResourceData) (*compute.ImageReference, error) {
	storageImageRefs := d.Get("storage_profile_image_reference").(*schema.Set).List()

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

func expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d *schema.ResourceData) (*compute.LinuxConfiguration, error) {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*schema.Set).List()

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

	return config, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d *schema.ResourceData) (*compute.WindowsConfiguration, error) {
	osProfilesWindowsConfig := d.Get("os_profile_windows_config").(*schema.Set).List()

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
	return config, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d *schema.ResourceData) *[]compute.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").(*schema.Set).List()
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

func expandAzureRMVirtualMachineScaleSetExtensions(d *schema.ResourceData) (*compute.VirtualMachineScaleSetExtensionProfile, error) {
	extensions := d.Get("extension").(*schema.Set).List()
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

		if s := config["settings"].(string); s != "" {
			settings, err := structure.ExpandJsonFromString(s)
			if err != nil {
				return nil, fmt.Errorf("unable to parse settings: %+v", err)
			}
			extension.VirtualMachineScaleSetExtensionProperties.Settings = &settings
		}

		if s := config["protected_settings"].(string); s != "" {
			protectedSettings, err := structure.ExpandJsonFromString(s)
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

func expandAzureRmVirtualMachineScaleSetPlan(d *schema.ResourceData) (*compute.Plan, error) {
	planConfigs := d.Get("plan").(*schema.Set).List()

	planConfig := planConfigs[0].(map[string]interface{})

	publisher := planConfig["publisher"].(string)
	name := planConfig["name"].(string)
	product := planConfig["product"].(string)

	return &compute.Plan{
		Publisher: &publisher,
		Name:      &name,
		Product:   &product,
	}, nil
}

func flattenAzureRmVirtualMachineScaleSetPlan(plan *compute.Plan) []interface{} {
	result := make(map[string]interface{})

	result["name"] = *plan.Name
	result["publisher"] = *plan.Publisher
	result["product"] = *plan.Product

	return []interface{}{result}
}
