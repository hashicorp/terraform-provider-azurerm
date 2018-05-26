package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineCreate,
		Read:   resourceArmVirtualMachineRead,
		Update: resourceArmVirtualMachineCreate,
		Delete: resourceArmVirtualMachineDelete,
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

			"zones": singleZonesSchema(),

			"plan": {
				Type:     schema.TypeList,
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

			"availability_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				StateFunc: func(id interface{}) string {
					return strings.ToLower(id.(string))
				},
			},

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

			"license_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					"Windows_Client",
					"Windows_Server",
				}, true),
			},

			"vm_size": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"storage_image_reference": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"offer": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"sku": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceArmVirtualMachineStorageImageReferenceHash,
			},

			"storage_os_disk": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.Linux),
								string(compute.Windows),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vhd_uri": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ConflictsWith: []string{
								"storage_os_disk.0.managed_disk_id",
								"storage_os_disk.0.managed_disk_type",
							},
						},

						"managed_disk_id": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_os_disk.0.vhd_uri"},
						},

						"managed_disk_type": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_os_disk.0.vhd_uri"},
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.PremiumLRS),
								string(compute.StandardLRS),
							}, true),
						},

						"image_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"create_option": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"disk_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateDiskSizeGB,
						},

						"write_accelerator_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"delete_os_disk_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"storage_data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vhd_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"managed_disk_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

						"create_option": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

						"lun": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"write_accelerator_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"delete_data_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"boot_diagnostics": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"storage_uri": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"os_profile": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"computer_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},

						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
						},

						"admin_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"custom_data": {
							Type:      schema.TypeString,
							ForceNew:  true,
							Optional:  true,
							Computed:  true,
							StateFunc: userDataStateFunc,
						},
					},
				},
				Set: resourceArmVirtualMachineStorageOsProfileHash,
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
							Default:  false,
						},
						"enable_automatic_upgrades": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
				Set:           resourceArmVirtualMachineStorageOsProfileWindowsConfigHash,
				ConflictsWith: []string{"os_profile_linux_config"},
			},

			"os_profile_linux_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_password_authentication": {
							Type:     schema.TypeBool,
							Required: true,
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
				Set:           resourceArmVirtualMachineStorageOsProfileLinuxConfigHash,
				ConflictsWith: []string{"os_profile_windows_config"},
			},

			"os_profile_secrets": {
				Type:     schema.TypeList,
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

			"network_interface_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"primary_network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)
	zones := expandZones(d.Get("zones").([]interface{}))

	osDisk, err := expandAzureRmVirtualMachineOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile := compute.StorageProfile{
		OsDisk: osDisk,
	}

	if _, ok := d.GetOk("storage_image_reference"); ok {
		imageRef, err := expandAzureRmVirtualMachineImageReference(d)
		if err != nil {
			return err
		}
		storageProfile.ImageReference = imageRef
	}

	if _, ok := d.GetOk("storage_data_disk"); ok {
		dataDisks, err := expandAzureRmVirtualMachineDataDisk(d)
		if err != nil {
			return err
		}
		storageProfile.DataDisks = &dataDisks
	}

	networkProfile := expandAzureRmVirtualMachineNetworkProfile(d)
	vmSize := d.Get("vm_size").(string)
	properties := compute.VirtualMachineProperties{
		NetworkProfile: &networkProfile,
		HardwareProfile: &compute.HardwareProfile{
			VMSize: compute.VirtualMachineSizeTypes(vmSize),
		},
		StorageProfile: &storageProfile,
	}

	if v, ok := d.GetOk("license_type"); ok {
		license := v.(string)
		properties.LicenseType = &license
	}

	if _, ok := d.GetOk("boot_diagnostics"); ok {
		diagnosticsProfile := expandAzureRmVirtualMachineDiagnosticsProfile(d)
		if diagnosticsProfile != nil {
			properties.DiagnosticsProfile = diagnosticsProfile
		}
	}

	if _, ok := d.GetOk("os_profile"); ok {
		osProfile, err := expandAzureRmVirtualMachineOsProfile(d)
		if err != nil {
			return err
		}
		properties.OsProfile = osProfile
	}

	if v, ok := d.GetOk("availability_set_id"); ok {
		availabilitySet := v.(string)
		availSet := compute.SubResource{
			ID: &availabilitySet,
		}

		properties.AvailabilitySet = &availSet
	}

	vm := compute.VirtualMachine{
		Name:                     &name,
		Location:                 &location,
		VirtualMachineProperties: &properties,
		Tags:  expandedTags,
		Zones: zones,
	}

	if _, ok := d.GetOk("identity"); ok {
		vmIdentity := expandAzureRmVirtualMachineIdentity(d)
		vm.Identity = vmIdentity
	}

	if _, ok := d.GetOk("plan"); ok {
		plan, err := expandAzureRmVirtualMachinePlan(d)
		if err != nil {
			return err
		}

		vm.Plan = plan
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, vm)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Machine %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualMachineRead(d, meta)
}

func resourceArmVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	vmClient := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachines"]

	resp, err := vmClient.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Virtual Machine %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("zones", resp.Zones)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if resp.Plan != nil {
		if err := d.Set("plan", flattenAzureRmVirtualMachinePlan(resp.Plan)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Plan error: %#v", err)
		}
	}

	d.Set("identity", flattenAzureRmVirtualMachineIdentity(resp.Identity))

	if resp.VirtualMachineProperties.AvailabilitySet != nil {
		d.Set("availability_set_id", strings.ToLower(*resp.VirtualMachineProperties.AvailabilitySet.ID))
	}

	d.Set("vm_size", resp.VirtualMachineProperties.HardwareProfile.VMSize)

	if resp.VirtualMachineProperties.StorageProfile.ImageReference != nil {
		if err := d.Set("storage_image_reference", schema.NewSet(resourceArmVirtualMachineStorageImageReferenceHash, flattenAzureRmVirtualMachineImageReference(resp.VirtualMachineProperties.StorageProfile.ImageReference))); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Image Reference error: %#v", err)
		}
	}

	if osDisk := resp.VirtualMachineProperties.StorageProfile.OsDisk; osDisk != nil {
		diskInfo, err := resourceArmVirtualMachineGetManagedDiskInfo(osDisk.ManagedDisk, meta)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error getting managed OS disk detailed information: %#v", err)
		}
		if err := d.Set("storage_os_disk", flattenAzureRmVirtualMachineOsDisk(osDisk, diskInfo)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Disk error: %#v", err)
		}
	}

	if dataDisks := resp.VirtualMachineProperties.StorageProfile.DataDisks; dataDisks != nil {
		disksInfo := make([]*compute.Disk, len(*dataDisks))
		for i, dataDisk := range *dataDisks {
			diskInfo, err := resourceArmVirtualMachineGetManagedDiskInfo(dataDisk.ManagedDisk, meta)
			if err != nil {
				return fmt.Errorf("[DEBUG] Error getting managed data disk detailed information: %#v", err)
			}
			disksInfo[i] = diskInfo
		}
		if err := d.Set("storage_data_disk", flattenAzureRmVirtualMachineDataDisk(dataDisks, disksInfo)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Data Disks error: %#v", err)
		}
	}

	if resp.VirtualMachineProperties.OsProfile != nil {
		if err := d.Set("os_profile", schema.NewSet(resourceArmVirtualMachineStorageOsProfileHash, flattenAzureRmVirtualMachineOsProfile(resp.VirtualMachineProperties.OsProfile))); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile: %#v", err)
		}

		if resp.VirtualMachineProperties.OsProfile.WindowsConfiguration != nil {
			if err := d.Set("os_profile_windows_config", schema.NewSet(resourceArmVirtualMachineStorageOsProfileWindowsConfigHash, flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(resp.VirtualMachineProperties.OsProfile.WindowsConfiguration))); err != nil {
				return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile Windows Configuration: %#v", err)
			}
		}

		if resp.VirtualMachineProperties.OsProfile.LinuxConfiguration != nil {
			if err := d.Set("os_profile_linux_config", schema.NewSet(resourceArmVirtualMachineStorageOsProfileLinuxConfigHash, flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(resp.VirtualMachineProperties.OsProfile.LinuxConfiguration))); err != nil {
				return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile Linux Configuration: %#v", err)
			}
		}

		if resp.VirtualMachineProperties.OsProfile.Secrets != nil {
			if err := d.Set("os_profile_secrets", flattenAzureRmVirtualMachineOsProfileSecrets(resp.VirtualMachineProperties.OsProfile.Secrets)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile Secrets: %#v", err)
			}
		}
	}

	if resp.VirtualMachineProperties.DiagnosticsProfile != nil && resp.VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics != nil {
		if err := d.Set("boot_diagnostics", flattenAzureRmVirtualMachineDiagnosticsProfile(resp.VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Diagnostics Profile: %#v", err)
		}
	}

	if resp.VirtualMachineProperties.NetworkProfile != nil {
		if err := d.Set("network_interface_ids", flattenAzureRmVirtualMachineNetworkInterfaces(resp.VirtualMachineProperties.NetworkProfile)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Network Interfaces: %#v", err)
		}

		if resp.VirtualMachineProperties.NetworkProfile.NetworkInterfaces != nil {
			for _, nic := range *resp.VirtualMachineProperties.NetworkProfile.NetworkInterfaces {
				if props := nic.NetworkInterfaceReferenceProperties; props != nil {
					if props.Primary != nil && *props.Primary {
						d.Set("primary_network_interface_id", nic.ID)
						break
					}
				}
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachines"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	// delete OS Disk if opted in
	if deleteOsDisk := d.Get("delete_os_disk_on_termination").(bool); deleteOsDisk {
		log.Printf("[INFO] delete_os_disk_on_termination is enabled, deleting disk from %s", name)

		osDisk, err := expandAzureRmVirtualMachineOsDisk(d)
		if err != nil {
			return fmt.Errorf("Error expanding OS Disk: %s", err)
		}

		if osDisk.Vhd != nil {
			if err = resourceArmVirtualMachineDeleteVhd(*osDisk.Vhd.URI, meta); err != nil {
				return fmt.Errorf("Error deleting OS Disk VHD: %+v", err)
			}
		} else if osDisk.ManagedDisk != nil {
			if err = resourceArmVirtualMachineDeleteManagedDisk(*osDisk.ManagedDisk.ID, meta); err != nil {
				return fmt.Errorf("Error deleting OS Managed Disk: %+v", err)
			}
		} else {
			return fmt.Errorf("Unable to locate OS managed disk properties from %s", name)
		}
	}

	// delete Data disks if opted in
	if deleteDataDisks := d.Get("delete_data_disks_on_termination").(bool); deleteDataDisks {
		log.Printf("[INFO] delete_data_disks_on_termination is enabled, deleting each data disk from %s", name)

		disks, err := expandAzureRmVirtualMachineDataDisk(d)
		if err != nil {
			return fmt.Errorf("Error expanding Data Disks: %s", err)
		}

		for _, disk := range disks {
			if disk.Vhd != nil {
				if err = resourceArmVirtualMachineDeleteVhd(*disk.Vhd.URI, meta); err != nil {
					return fmt.Errorf("Error deleting Data Disk VHD: %+v", err)
				}
			} else if disk.ManagedDisk != nil {
				if err = resourceArmVirtualMachineDeleteManagedDisk(*disk.ManagedDisk.ID, meta); err != nil {
					return fmt.Errorf("Error deleting Data Managed Disk: %+v", err)
				}
			} else {
				return fmt.Errorf("Unable to locate data managed disk properties from %s", name)
			}
		}
	}

	return nil
}

func resourceArmVirtualMachineDeleteVhd(uri string, meta interface{}) error {
	vhdURL, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("Cannot parse Disk VHD URI: %s", err)
	}

	// VHD URI is in the form: https://storageAccountName.blob.core.windows.net/containerName/blobName
	storageAccountName := strings.Split(vhdURL.Host, ".")[0]
	path := strings.Split(strings.TrimPrefix(vhdURL.Path, "/"), "/")
	containerName := path[0]
	blobName := path[1]

	storageAccountResourceGroupName, err := findStorageAccountResourceGroup(meta, storageAccountName)
	if err != nil {
		return fmt.Errorf("Error finding resource group for storage account %s: %+v", storageAccountName, err)
	}

	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	blobClient, saExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, storageAccountResourceGroupName, storageAccountName)
	if err != nil {
		return fmt.Errorf("Error creating blob store client for VHD deletion: %+v", err)
	}

	if !saExists {
		log.Printf("[INFO] Storage Account %q in resource group %q doesn't exist so the VHD blob won't exist", storageAccountName, storageAccountResourceGroupName)
		return nil
	}

	log.Printf("[INFO] Deleting VHD blob %s", blobName)
	container := blobClient.GetContainerReference(containerName)
	blob := container.GetBlobReference(blobName)
	options := &storage.DeleteBlobOptions{}
	err = blob.Delete(options)
	if err != nil {
		return fmt.Errorf("Error deleting VHD blob: %+v", err)
	}

	return nil
}

func resourceArmVirtualMachineDeleteManagedDisk(managedDiskID string, meta interface{}) error {
	client := meta.(*ArmClient).diskClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(managedDiskID)
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Disk (%s %s) %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Disk (%s %s) %+v", name, resGroup, err)
	}

	return nil
}

func flattenAzureRmVirtualMachinePlan(plan *compute.Plan) []interface{} {
	result := make(map[string]interface{})
	result["name"] = *plan.Name
	result["publisher"] = *plan.Publisher
	result["product"] = *plan.Product

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineImageReference(image *compute.ImageReference) []interface{} {
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

func flattenAzureRmVirtualMachineIdentity(identity *compute.VirtualMachineIdentity) []interface{} {
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

func flattenAzureRmVirtualMachineDiagnosticsProfile(profile *compute.BootDiagnostics) []interface{} {
	result := make(map[string]interface{})

	result["enabled"] = *profile.Enabled
	if profile.StorageURI != nil {
		result["storage_uri"] = *profile.StorageURI
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineNetworkInterfaces(profile *compute.NetworkProfile) []interface{} {
	result := make([]interface{}, 0)
	for _, nic := range *profile.NetworkInterfaces {
		result = append(result, *nic.ID)
	}
	return result
}

func flattenAzureRmVirtualMachineOsProfileSecrets(secrets *[]compute.VaultSecretGroup) []interface{} {
	result := make([]interface{}, 0)
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

func flattenAzureRmVirtualMachineDataDisk(disks *[]compute.DataDisk, disksInfo []*compute.Disk) interface{} {
	result := make([]interface{}, len(*disks))
	for i, disk := range *disks {
		l := make(map[string]interface{})
		l["name"] = *disk.Name
		if disk.Vhd != nil {
			l["vhd_uri"] = *disk.Vhd.URI
		}
		if disk.ManagedDisk != nil {
			l["managed_disk_type"] = string(disk.ManagedDisk.StorageAccountType)
			l["managed_disk_id"] = *disk.ManagedDisk.ID
		}
		l["create_option"] = disk.CreateOption
		l["caching"] = string(disk.Caching)
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}
		l["lun"] = *disk.Lun

		if v := disk.WriteAcceleratorEnabled; v != nil {
			l["write_accelerator_enabled"] = *disk.WriteAcceleratorEnabled
		}

		flattenAzureRmVirtualMachineReviseDiskInfo(l, disksInfo[i])

		result[i] = l
	}
	return result
}

func flattenAzureRmVirtualMachineOsProfile(input *compute.OSProfile) []interface{} {
	result := make(map[string]interface{})
	result["computer_name"] = *input.ComputerName
	result["admin_username"] = *input.AdminUsername
	if input.CustomData != nil {
		result["custom_data"] = *input.CustomData
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(config *compute.WindowsConfiguration) []interface{} {
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

func flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(config *compute.LinuxConfiguration) []interface{} {

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

func flattenAzureRmVirtualMachineOsDisk(disk *compute.OSDisk, diskInfo *compute.Disk) []interface{} {
	result := make(map[string]interface{})
	result["name"] = *disk.Name
	if disk.Vhd != nil {
		result["vhd_uri"] = *disk.Vhd.URI
	}
	if disk.ManagedDisk != nil {
		result["managed_disk_type"] = string(disk.ManagedDisk.StorageAccountType)
		if disk.ManagedDisk.ID != nil {
			result["managed_disk_id"] = *disk.ManagedDisk.ID
		}
	}
	result["create_option"] = disk.CreateOption
	result["caching"] = disk.Caching
	if disk.DiskSizeGB != nil {
		result["disk_size_gb"] = *disk.DiskSizeGB
	}
	result["os_type"] = string(disk.OsType)

	if v := disk.WriteAcceleratorEnabled; v != nil {
		result["write_accelerator_enabled"] = *disk.WriteAcceleratorEnabled
	}

	flattenAzureRmVirtualMachineReviseDiskInfo(result, diskInfo)

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineReviseDiskInfo(result map[string]interface{}, diskInfo *compute.Disk) {
	if diskInfo != nil {
		if diskInfo.Sku != nil {
			result["managed_disk_type"] = string(diskInfo.Sku.Name)
		}
		if diskInfo.DiskProperties != nil && diskInfo.DiskProperties.DiskSizeGB != nil {
			result["disk_size_gb"] = *diskInfo.DiskProperties.DiskSizeGB
		}
	}
}

func expandAzureRmVirtualMachinePlan(d *schema.ResourceData) (*compute.Plan, error) {
	planConfigs := d.Get("plan").([]interface{})

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

func expandAzureRmVirtualMachineIdentity(d *schema.ResourceData) *compute.VirtualMachineIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := identity["type"].(string)
	return &compute.VirtualMachineIdentity{
		Type: compute.ResourceIdentityType(identityType),
	}
}

func expandAzureRmVirtualMachineOsProfile(d *schema.ResourceData) (*compute.OSProfile, error) {
	osProfiles := d.Get("os_profile").(*schema.Set).List()

	osProfile := osProfiles[0].(map[string]interface{})

	adminUsername := osProfile["admin_username"].(string)
	adminPassword := osProfile["admin_password"].(string)
	computerName := osProfile["computer_name"].(string)

	profile := &compute.OSProfile{
		AdminUsername: &adminUsername,
		ComputerName:  &computerName,
	}

	if adminPassword != "" {
		profile.AdminPassword = &adminPassword
	}

	if _, ok := d.GetOk("os_profile_windows_config"); ok {
		winConfig, err := expandAzureRmVirtualMachineOsProfileWindowsConfig(d)
		if err != nil {
			return nil, err
		}
		if winConfig != nil {
			profile.WindowsConfiguration = winConfig
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		linuxConfig, err := expandAzureRmVirtualMachineOsProfileLinuxConfig(d)
		if err != nil {
			return nil, err
		}
		if linuxConfig != nil {
			profile.LinuxConfiguration = linuxConfig
		}
	}

	if profile.LinuxConfiguration == nil && profile.WindowsConfiguration == nil {
		return nil, fmt.Errorf("Error: either a `os_profile_linux_config` or a `os_profile_windows_config` must be specified.")
	}

	if _, ok := d.GetOk("os_profile_secrets"); ok {
		secrets := expandAzureRmVirtualMachineOsProfileSecrets(d)
		if secrets != nil {
			profile.Secrets = secrets
		}
	}

	if v := osProfile["custom_data"].(string); v != "" {
		v = base64Encode(v)
		profile.CustomData = &v
	}

	return profile, nil
}

func expandAzureRmVirtualMachineOsProfileSecrets(d *schema.ResourceData) *[]compute.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").([]interface{})
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

func expandAzureRmVirtualMachineOsProfileLinuxConfig(d *schema.ResourceData) (*compute.LinuxConfiguration, error) {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*schema.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]interface{})
	disablePasswordAuth := linuxConfig["disable_password_authentication"].(bool)

	config := &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &disablePasswordAuth,
	}

	linuxKeys := linuxConfig["ssh_keys"].([]interface{})
	sshPublicKeys := []compute.SSHPublicKey{}
	for _, key := range linuxKeys {

		sshKey, ok := key.(map[string]interface{})
		if !ok {
			continue
		}
		path := sshKey["path"].(string)
		keyData := sshKey["key_data"].(string)

		sshPublicKey := compute.SSHPublicKey{
			Path:    &path,
			KeyData: &keyData,
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	if len(sshPublicKeys) > 0 {
		config.SSH = &compute.SSHConfiguration{
			PublicKeys: &sshPublicKeys,
		}
	}

	return config, nil
}

func expandAzureRmVirtualMachineOsProfileWindowsConfig(d *schema.ResourceData) (*compute.WindowsConfiguration, error) {
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

func expandAzureRmVirtualMachineDataDisk(d *schema.ResourceData) ([]compute.DataDisk, error) {
	disks := d.Get("storage_data_disk").([]interface{})
	data_disks := make([]compute.DataDisk, 0, len(disks))
	for _, disk_config := range disks {
		config := disk_config.(map[string]interface{})

		name := config["name"].(string)
		createOption := config["create_option"].(string)
		vhdURI := config["vhd_uri"].(string)
		managedDiskType := config["managed_disk_type"].(string)
		managedDiskID := config["managed_disk_id"].(string)
		lun := int32(config["lun"].(int))

		data_disk := compute.DataDisk{
			Name:         &name,
			Lun:          &lun,
			CreateOption: compute.DiskCreateOptionTypes(createOption),
		}

		if vhdURI != "" {
			data_disk.Vhd = &compute.VirtualHardDisk{
				URI: &vhdURI,
			}
		}

		managedDisk := &compute.ManagedDiskParameters{}

		if managedDiskType != "" {
			managedDisk.StorageAccountType = compute.StorageAccountTypes(managedDiskType)
			data_disk.ManagedDisk = managedDisk
		}

		if managedDiskID != "" {
			managedDisk.ID = &managedDiskID
			data_disk.ManagedDisk = managedDisk
		}

		if vhdURI != "" && managedDiskID != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
		}
		if vhdURI != "" && managedDiskType != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
		}
		if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(data_disk.CreateOption), string(compute.Attach)) {
			return nil, fmt.Errorf("[ERROR] Must specify `vhd_uri` or `managed_disk_id` to attach")
		}

		if v := config["caching"].(string); v != "" {
			data_disk.Caching = compute.CachingTypes(v)
		}

		if v, ok := config["disk_size_gb"].(int); ok {
			data_disk.DiskSizeGB = utils.Int32(int32(v))
		}

		if v, ok := config["write_accelerator_enabled"].(bool); ok {
			data_disk.WriteAcceleratorEnabled = utils.Bool(v)
		}

		data_disks = append(data_disks, data_disk)
	}

	return data_disks, nil
}

func expandAzureRmVirtualMachineDiagnosticsProfile(d *schema.ResourceData) *compute.DiagnosticsProfile {
	bootDiagnostics := d.Get("boot_diagnostics").([]interface{})

	diagnosticsProfile := &compute.DiagnosticsProfile{}
	if len(bootDiagnostics) > 0 {
		bootDiagnostic := bootDiagnostics[0].(map[string]interface{})

		diagnostic := &compute.BootDiagnostics{
			Enabled:    utils.Bool(bootDiagnostic["enabled"].(bool)),
			StorageURI: utils.String(bootDiagnostic["storage_uri"].(string)),
		}

		diagnosticsProfile.BootDiagnostics = diagnostic

		return diagnosticsProfile
	}

	return nil
}

func expandAzureRmVirtualMachineImageReference(d *schema.ResourceData) (*compute.ImageReference, error) {
	storageImageRefs := d.Get("storage_image_reference").(*schema.Set).List()

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

		imageReference = compute.ImageReference{
			Publisher: &publisher,
			Offer:     &offer,
			Sku:       &sku,
			Version:   &version,
		}
	}

	return &imageReference, nil
}

func expandAzureRmVirtualMachineNetworkProfile(d *schema.ResourceData) compute.NetworkProfile {
	nicIds := d.Get("network_interface_ids").([]interface{})
	primaryNicId := d.Get("primary_network_interface_id").(string)
	network_interfaces := make([]compute.NetworkInterfaceReference, 0, len(nicIds))

	network_profile := compute.NetworkProfile{}

	for _, nic := range nicIds {
		id := nic.(string)
		primary := id == primaryNicId

		network_interface := compute.NetworkInterfaceReference{
			ID: &id,
			NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
				Primary: &primary,
			},
		}
		network_interfaces = append(network_interfaces, network_interface)
	}

	network_profile.NetworkInterfaces = &network_interfaces

	return network_profile
}

func expandAzureRmVirtualMachineOsDisk(d *schema.ResourceData) (*compute.OSDisk, error) {
	disks := d.Get("storage_os_disk").([]interface{})

	config := disks[0].(map[string]interface{})

	name := config["name"].(string)
	imageURI := config["image_uri"].(string)
	createOption := config["create_option"].(string)
	vhdURI := config["vhd_uri"].(string)
	managedDiskType := config["managed_disk_type"].(string)
	managedDiskID := config["managed_disk_id"].(string)

	osDisk := &compute.OSDisk{
		Name:         &name,
		CreateOption: compute.DiskCreateOptionTypes(createOption),
	}

	if vhdURI != "" {
		osDisk.Vhd = &compute.VirtualHardDisk{
			URI: &vhdURI,
		}
	}

	managedDisk := &compute.ManagedDiskParameters{}

	if managedDiskType != "" {
		managedDisk.StorageAccountType = compute.StorageAccountTypes(managedDiskType)
		osDisk.ManagedDisk = managedDisk
	}

	if managedDiskID != "" {
		managedDisk.ID = &managedDiskID
		osDisk.ManagedDisk = managedDisk
	}

	//BEGIN: code to be removed after GH-13016 is merged
	if vhdURI != "" && managedDiskID != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
	}
	if vhdURI != "" && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
	}
	//END: code to be removed after GH-13016 is merged
	if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(osDisk.CreateOption), string(compute.Attach)) {
		return nil, fmt.Errorf("[ERROR] Must specify `vhd_uri` or `managed_disk_id` to attach")
	}

	if v := config["image_uri"].(string); v != "" {
		osDisk.Image = &compute.VirtualHardDisk{
			URI: &imageURI,
		}
	}

	if v := config["os_type"].(string); v != "" {
		osDisk.OsType = compute.OperatingSystemTypes(v)
	}

	if v := config["caching"].(string); v != "" {
		osDisk.Caching = compute.CachingTypes(v)
	}

	if v := config["disk_size_gb"].(int); v != 0 {
		osDisk.DiskSizeGB = utils.Int32(int32(v))
	}

	if v, ok := config["write_accelerator_enabled"].(bool); ok {
		osDisk.WriteAcceleratorEnabled = utils.Bool(v)
	}

	return osDisk, nil
}

func findStorageAccountResourceGroup(meta interface{}, storageAccountName string) (string, error) {
	client := meta.(*ArmClient).resourcesClient
	ctx := meta.(*ArmClient).StopContext
	filter := fmt.Sprintf("name eq '%s' and resourceType eq 'Microsoft.Storage/storageAccounts'", storageAccountName)
	expand := ""
	var pager *int32

	rf, err := client.List(ctx, filter, expand, pager)
	if err != nil {
		return "", fmt.Errorf("Error making resource request for query %s: %+v", filter, err)
	}

	results := rf.Values()
	if len(results) != 1 {
		return "", fmt.Errorf("Wrong number of results making resource request for query %s: %d", filter, len(results))
	}

	id, err := parseAzureResourceID(*results[0].ID)
	if err != nil {
		return "", err
	}

	return id.ResourceGroup, nil
}

func resourceArmVirtualMachineStorageOsProfileHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["admin_username"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["computer_name"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineStorageOsProfileWindowsConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["provision_vm_agent"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
		if v, ok := m["enable_automatic_upgrades"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineStorageOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineStorageImageReferenceHash(v interface{}) int {
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
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineGetManagedDiskInfo(disk *compute.ManagedDiskParameters, meta interface{}) (*compute.Disk, error) {
	client := meta.(*ArmClient).diskClient
	ctx := meta.(*ArmClient).StopContext

	if disk == nil || disk.ID == nil {
		return nil, nil
	}

	diskId := *disk.ID
	id, err := parseAzureResourceID(diskId)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Disk ID %q: %+v", diskId, err)
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["disks"]
	diskResp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return &diskResp, nil
}
