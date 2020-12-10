package compute

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	intStor "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
	"golang.org/x/net/context"
)

var virtualMachineResourceName = "azurerm_virtual_machine"

// TODO move into internal/tf/suppress/base64.go
func userDataDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return userDataStateFunc(old) == new
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		s = utils.Base64EncodeIfNot(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

// NOTE: the `azurerm_virtual_machine` resource has been superseded by the `azurerm_linux_virtual_machine` and
// 		 `azurerm_windows_virtual_machine` resources - as such this resource is feature-frozen and new
//		 functionality will be added to these new resources instead.
func resourceArmVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineCreateUpdate,
		Read:   resourceArmVirtualMachineRead,
		Update: resourceArmVirtualMachineCreateUpdate,
		Delete: resourceArmVirtualMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zones": azure.SchemaSingleZone(),

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
				ConflictsWith: []string{"zones"},
			},

			"proximity_placement_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,

				// We have to ignore case due to incorrect capitalisation of resource group name in
				// proximity placement group ID in the response we get from the API request
				//
				// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
				DiffSuppressFunc: suppress.CaseDifference,
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
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.ResourceIdentityTypeSystemAssigned),
								string(compute.ResourceIdentityTypeUserAssigned),
								string(compute.ResourceIdentityTypeSystemAssignedUserAssigned),
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity_ids": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
						},
					},
				},
			},

			"license_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"Windows_Client",
					"Windows_Server",
				}, true),
			},

			"vm_size": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			// lintignore:S018
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
							DiffSuppressFunc: suppress.CaseDifference,
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
								string(compute.StorageAccountTypesPremiumLRS),
								string(compute.StorageAccountTypesStandardLRS),
								string(compute.StorageAccountTypesStandardSSDLRS),
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
							DiffSuppressFunc: suppress.CaseDifference,
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
				Computed: true,
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
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"managed_disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StorageAccountTypesPremiumLRS),
								string(compute.StorageAccountTypesStandardLRS),
								string(compute.StorageAccountTypesStandardSSDLRS),
								string(compute.StorageAccountTypesUltraSSDLRS),
							}, true),
						},

						"create_option": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
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

			"additional_capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ultra_ssd_enabled": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			// lintignore:S018
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

			// lintignore:S018
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
						"timezone": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.VirtualMachineTimeZoneCaseInsensitive(),
						},
						"winrm": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"HTTP",
											"HTTPS",
										}, true),
										DiffSuppressFunc: suppress.CaseDifference,
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
									// TODO: should we make `pass` and `component` Optional + Defaulted?
									"pass": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"oobeSystem",
										}, false),
									},
									"component": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Microsoft-Windows-Shell-Setup",
										}, false),
									},
									"setting_name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"AutoLogon",
											"FirstLogonCommands",
										}, false),
									},
									"content": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
				Set:           resourceArmVirtualMachineStorageOsProfileWindowsConfigHash,
				ConflictsWith: []string{"os_profile_linux_config"},
			},

			// lintignore:S018
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
										Required: true,
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
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"primary_network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVirtualMachineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Machine %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_machine", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	osDisk, err := expandAzureRmVirtualMachineOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile := compute.StorageProfile{
		OsDisk: osDisk,
	}

	if _, ok := d.GetOk("storage_image_reference"); ok {
		imageRef, err2 := expandAzureRmVirtualMachineImageReference(d)
		if err2 != nil {
			return err2
		}
		storageProfile.ImageReference = imageRef
	}

	if _, ok := d.GetOk("storage_data_disk"); ok {
		dataDisks, err2 := expandAzureRmVirtualMachineDataDisk(d)
		if err2 != nil {
			return err2
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
	if _, ok := d.GetOk("additional_capabilities"); ok {
		properties.AdditionalCapabilities = expandAzureRmVirtualMachineAdditionalCapabilities(d)
	}

	if _, ok := d.GetOk("os_profile"); ok {
		osProfile, err2 := expandAzureRmVirtualMachineOsProfile(d)
		if err2 != nil {
			return err2
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

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		properties.ProximityPlacementGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	vm := compute.VirtualMachine{
		Name:                     &name,
		Location:                 &location,
		VirtualMachineProperties: &properties,
		Tags:                     expandedTags,
		Zones:                    zones,
	}

	if _, ok := d.GetOk("identity"); ok {
		vm.Identity = expandAzureRmVirtualMachineIdentity(d)
	}

	if _, ok := d.GetOk("plan"); ok {
		vm.Plan = expandAzureRmVirtualMachinePlan(d)
	}

	locks.ByName(name, virtualMachineResourceName)
	defer locks.UnlockByName(name, virtualMachineResourceName)

	future, err := client.CreateOrUpdate(ctx, resGroup, name, vm)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

	ipAddress, err := determineVirtualMachineIPAddress(ctx, meta, read.VirtualMachineProperties)
	if err != nil {
		return fmt.Errorf("Error determining IP Address for Virtual Machine %q (Resource Group %q): %+v", name, resGroup, err)
	}

	provisionerType := "ssh"
	if props := read.VirtualMachineProperties; props != nil {
		if profile := props.OsProfile; profile != nil {
			if profile.WindowsConfiguration != nil {
				provisionerType = "winrm"
			}
		}
	}
	d.SetConnInfo(map[string]string{
		"type": provisionerType,
		"host": ipAddress,
	})

	return resourceArmVirtualMachineRead(d, meta)
}

func resourceArmVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	vmclient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachines"]

	resp, err := vmclient.Get(ctx, resGroup, name, "")
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("plan", flattenAzureRmVirtualMachinePlan(resp.Plan)); err != nil {
		return fmt.Errorf("Error setting `plan`: %#v", err)
	}

	if err := d.Set("identity", flattenAzureRmVirtualMachineIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if props := resp.VirtualMachineProperties; props != nil {
		if availabilitySet := props.AvailabilitySet; availabilitySet != nil {
			// Lowercase due to incorrect capitalisation of resource group name in
			// availability set ID in response from get VM API request
			// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
			d.Set("availability_set_id", strings.ToLower(*availabilitySet.ID))
		}

		if proximityPlacementGroup := props.ProximityPlacementGroup; proximityPlacementGroup != nil {
			d.Set("proximity_placement_group_id", proximityPlacementGroup.ID)
		}

		if profile := props.HardwareProfile; profile != nil {
			d.Set("vm_size", profile.VMSize)
		}

		if profile := props.StorageProfile; profile != nil {
			if err := d.Set("storage_image_reference", schema.NewSet(resourceArmVirtualMachineStorageImageReferenceHash, flattenAzureRmVirtualMachineImageReference(profile.ImageReference))); err != nil {
				return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Image Reference error: %#v", err)
			}

			if osDisk := profile.OsDisk; osDisk != nil {
				diskInfo, err := resourceArmVirtualMachineGetManagedDiskInfo(d, osDisk.ManagedDisk, meta)
				if err != nil {
					return fmt.Errorf("Error flattening `storage_os_disk`: %#v", err)
				}
				if err := d.Set("storage_os_disk", flattenAzureRmVirtualMachineOsDisk(osDisk, diskInfo)); err != nil {
					return fmt.Errorf("Error setting `storage_os_disk`: %#v", err)
				}
			}

			if dataDisks := profile.DataDisks; dataDisks != nil {
				disksInfo := make([]*compute.Disk, len(*dataDisks))
				for i, dataDisk := range *dataDisks {
					diskInfo, err := resourceArmVirtualMachineGetManagedDiskInfo(d, dataDisk.ManagedDisk, meta)
					if err != nil {
						return fmt.Errorf("[DEBUG] Error getting managed data disk detailed information: %#v", err)
					}
					disksInfo[i] = diskInfo
				}
				if err := d.Set("storage_data_disk", flattenAzureRmVirtualMachineDataDisk(dataDisks, disksInfo)); err != nil {
					return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Data Disks error: %#v", err)
				}
			}
		}

		if profile := props.OsProfile; profile != nil {
			if err := d.Set("os_profile", schema.NewSet(resourceArmVirtualMachineStorageOsProfileHash, flattenAzureRmVirtualMachineOsProfile(profile))); err != nil {
				return fmt.Errorf("Error setting `os_profile`: %#v", err)
			}

			if err := d.Set("os_profile_linux_config", schema.NewSet(resourceArmVirtualMachineStorageOsProfileLinuxConfigHash, flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(profile.LinuxConfiguration))); err != nil {
				return fmt.Errorf("Error setting `os_profile_linux_config`: %+v", err)
			}

			if err := d.Set("os_profile_windows_config", schema.NewSet(resourceArmVirtualMachineStorageOsProfileWindowsConfigHash, flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(profile.WindowsConfiguration))); err != nil {
				return fmt.Errorf("Error setting `os_profile_windows_config`: %+v", err)
			}

			if err := d.Set("os_profile_secrets", flattenAzureRmVirtualMachineOsProfileSecrets(profile.Secrets)); err != nil {
				return fmt.Errorf("Error setting `os_profile_secrets`: %+v", err)
			}
		}

		if profile := props.DiagnosticsProfile; profile != nil {
			if err := d.Set("boot_diagnostics", flattenAzureRmVirtualMachineDiagnosticsProfile(profile.BootDiagnostics)); err != nil {
				return fmt.Errorf("Error setting `boot_diagnostics`: %#v", err)
			}
		}
		if err := d.Set("additional_capabilities", flattenAzureRmVirtualMachineAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
			return fmt.Errorf("Error setting `additional_capabilities`: %#v", err)
		}

		if profile := props.NetworkProfile; profile != nil {
			if err := d.Set("network_interface_ids", flattenAzureRmVirtualMachineNetworkInterfaces(profile)); err != nil {
				return fmt.Errorf("Error flattening `network_interface_ids`: %#v", err)
			}

			if profile.NetworkInterfaces != nil {
				for _, nic := range *profile.NetworkInterfaces {
					if props := nic.NetworkInterfaceReferenceProperties; props != nil {
						if props.Primary != nil && *props.Primary {
							d.Set("primary_network_interface_id", nic.ID)
							break
						}
					}
				}
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachines"]

	locks.ByName(name, virtualMachineResourceName)
	defer locks.UnlockByName(name, virtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Machine %q (Resource Group %q): %s", name, resGroup, err)
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Machine %q (Resource Group %q): %s", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Machine %q (Resource Group %q): %s", name, resGroup, err)
	}

	// delete OS Disk if opted in
	deleteOsDisk := d.Get("delete_os_disk_on_termination").(bool)
	deleteDataDisks := d.Get("delete_data_disks_on_termination").(bool)

	if deleteOsDisk || deleteDataDisks {
		storageClient := meta.(*clients.Client).Storage

		props := virtualMachine.VirtualMachineProperties
		if props == nil {
			return fmt.Errorf("Error deleting Disks for Virtual Machine %q - `props` was nil", name)
		}
		storageProfile := props.StorageProfile
		if storageProfile == nil {
			return fmt.Errorf("Error deleting Disks for Virtual Machine %q - `storageProfile` was nil", name)
		}

		if deleteOsDisk {
			log.Printf("[INFO] delete_os_disk_on_termination is enabled, deleting disk from %s", name)
			osDisk := storageProfile.OsDisk
			if osDisk == nil {
				return fmt.Errorf("Error deleting OS Disk for Virtual Machine %q - `osDisk` was nil", name)
			}
			if osDisk.Vhd == nil && osDisk.ManagedDisk == nil {
				return fmt.Errorf("Unable to determine OS Disk Type to Delete it for Virtual Machine %q", name)
			}

			if osDisk.Vhd != nil {
				if err = resourceArmVirtualMachineDeleteVhd(ctx, storageClient, osDisk.Vhd); err != nil {
					return fmt.Errorf("Error deleting OS Disk VHD: %+v", err)
				}
			} else if osDisk.ManagedDisk != nil {
				if err = resourceArmVirtualMachineDeleteManagedDisk(d, osDisk.ManagedDisk, meta); err != nil {
					return fmt.Errorf("Error deleting OS Managed Disk: %+v", err)
				}
			}
		}

		// delete Data disks if opted in
		if deleteDataDisks {
			log.Printf("[INFO] delete_data_disks_on_termination is enabled, deleting each data disk from %q", name)

			dataDisks := storageProfile.DataDisks
			if dataDisks == nil {
				return fmt.Errorf("Error deleting Data Disks for Virtual Machine %q: `dataDisks` was nil", name)
			}

			for _, disk := range *dataDisks {
				if disk.Vhd == nil && disk.ManagedDisk == nil {
					return fmt.Errorf("Unable to determine Data Disk Type to Delete it for Virtual Machine %q / Disk %q", name, *disk.Name)
				}

				if disk.Vhd != nil {
					if err = resourceArmVirtualMachineDeleteVhd(ctx, storageClient, disk.Vhd); err != nil {
						return fmt.Errorf("Error deleting Data Disk VHD: %+v", err)
					}
				} else if disk.ManagedDisk != nil {
					if err = resourceArmVirtualMachineDeleteManagedDisk(d, disk.ManagedDisk, meta); err != nil {
						return fmt.Errorf("Error deleting Data Managed Disk: %+v", err)
					}
				}
			}
		}
	}

	return nil
}

func resourceArmVirtualMachineDeleteVhd(ctx context.Context, storageClient *intStor.Client, vhd *compute.VirtualHardDisk) error {
	if vhd == nil {
		return fmt.Errorf("`vhd` was nil`")
	}
	if vhd.URI == nil {
		return fmt.Errorf("`vhd.URI` was nil`")
	}

	uri := *vhd.URI
	id, err := blobs.ParseResourceID(uri)
	if err != nil {
		return fmt.Errorf("Error parsing %q: %s", uri, err)
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Blob %q (Container %q): %s", id.AccountName, id.BlobName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q (Disk %q)!", id.AccountName, uri)
	}

	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	input := blobs.DeleteInput{
		DeleteSnapshots: false,
	}
	if _, err := blobsClient.Delete(ctx, id.AccountName, id.ContainerName, id.BlobName, input); err != nil {
		return fmt.Errorf("Error deleting Blob %q (Container %q / Account %q / Resource Group %q): %s", id.BlobName, id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func resourceArmVirtualMachineDeleteManagedDisk(d *schema.ResourceData, disk *compute.ManagedDiskParameters, meta interface{}) error {
	if disk == nil {
		return fmt.Errorf("`disk` was nil`")
	}
	if disk.ID == nil {
		return fmt.Errorf("`disk.ID` was nil`")
	}
	managedDiskID := *disk.ID

	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(managedDiskID)
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["disks"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Disk %q (Resource Group %q) %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Managed Disk %q (Resource Group %q) %+v", name, resGroup, err)
	}

	return nil
}

func flattenAzureRmVirtualMachinePlan(plan *compute.Plan) []interface{} {
	if plan == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if plan.Name != nil {
		result["name"] = *plan.Name
	}
	if plan.Publisher != nil {
		result["publisher"] = *plan.Publisher
	}
	if plan.Product != nil {
		result["product"] = *plan.Product
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineImageReference(image *compute.ImageReference) []interface{} {
	if image == nil {
		return []interface{}{}
	}

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

	identityIds := make([]string, 0)
	if identity.UserAssignedIdentities != nil {
		/*
			"userAssignedIdentities": {
			  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tomdevidentity/providers/Microsoft.ManagedIdentity/userAssignedIdentities/tom123": {
				"principalId": "00000000-0000-0000-0000-000000000000",
				"clientId": "00000000-0000-0000-0000-000000000000"
			  }
			}
		*/
		for key := range identity.UserAssignedIdentities {
			identityIds = append(identityIds, key)
		}
	}
	result["identity_ids"] = identityIds

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineDiagnosticsProfile(profile *compute.BootDiagnostics) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if profile.Enabled != nil {
		result["enabled"] = *profile.Enabled
	}

	if profile.StorageURI != nil {
		result["storage_uri"] = *profile.StorageURI
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineAdditionalCapabilities(profile *compute.AdditionalCapabilities) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	if v := profile.UltraSSDEnabled; v != nil {
		result["ultra_ssd_enabled"] = *v
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
	if secrets == nil {
		return []interface{}{}
	}

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
			if disk.ManagedDisk.ID != nil {
				l["managed_disk_id"] = *disk.ManagedDisk.ID
			}
		}
		l["create_option"] = disk.CreateOption
		l["caching"] = string(disk.Caching)
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}
		if v := disk.Lun; v != nil {
			l["lun"] = *v
		}

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
	if config == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if config.ProvisionVMAgent != nil {
		result["provision_vm_agent"] = *config.ProvisionVMAgent
	}

	if config.EnableAutomaticUpdates != nil {
		result["enable_automatic_upgrades"] = *config.EnableAutomaticUpdates
	}

	if config.TimeZone != nil {
		result["timezone"] = *config.TimeZone
	}

	listeners := make([]map[string]interface{}, 0)
	if config.WinRM != nil && config.WinRM.Listeners != nil {
		for _, i := range *config.WinRM.Listeners {
			listener := make(map[string]interface{})
			listener["protocol"] = string(i.Protocol)

			if i.CertificateURL != nil {
				listener["certificate_url"] = *i.CertificateURL
			}

			listeners = append(listeners, listener)
		}
	}

	result["winrm"] = listeners

	content := make([]map[string]interface{}, 0)
	if config.AdditionalUnattendContent != nil {
		for _, i := range *config.AdditionalUnattendContent {
			c := make(map[string]interface{})
			c["pass"] = string(i.PassName)
			c["component"] = string(i.ComponentName)
			c["setting_name"] = string(i.SettingName)

			if i.Content != nil {
				c["content"] = *i.Content
			}

			content = append(content, c)
		}
	}
	result["additional_unattend_config"] = content

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(config *compute.LinuxConfiguration) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if config.DisablePasswordAuthentication != nil {
		result["disable_password_authentication"] = *config.DisablePasswordAuthentication
	}

	if config.SSH != nil && config.SSH.PublicKeys != nil && len(*config.SSH.PublicKeys) > 0 {
		ssh_keys := make([]map[string]interface{}, 0)
		for _, i := range *config.SSH.PublicKeys {
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

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineOsDisk(disk *compute.OSDisk, diskInfo *compute.Disk) []interface{} {
	result := make(map[string]interface{})
	if disk.Name != nil {
		result["name"] = *disk.Name
	}
	if disk.Vhd != nil && disk.Vhd.URI != nil {
		result["vhd_uri"] = *disk.Vhd.URI
	}
	if disk.Image != nil && disk.Image.URI != nil {
		result["image_uri"] = *disk.Image.URI
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

func expandAzureRmVirtualMachinePlan(d *schema.ResourceData) *compute.Plan {
	planConfigs := d.Get("plan").([]interface{})
	if len(planConfigs) == 0 {
		return nil
	}

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

func expandAzureRmVirtualMachineIdentity(d *schema.ResourceData) *compute.VirtualMachineIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := compute.ResourceIdentityType(identity["type"].(string))

	identityIds := make(map[string]*compute.VirtualMachineIdentityUserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &compute.VirtualMachineIdentityUserAssignedIdentitiesValue{}
	}

	vmIdentity := compute.VirtualMachineIdentity{
		Type: identityType,
	}

	if vmIdentity.Type == compute.ResourceIdentityTypeUserAssigned || vmIdentity.Type == compute.ResourceIdentityTypeSystemAssignedUserAssigned {
		vmIdentity.UserAssignedIdentities = identityIds
	}

	return &vmIdentity
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
		winConfig := expandAzureRmVirtualMachineOsProfileWindowsConfig(d)
		if winConfig != nil {
			profile.WindowsConfiguration = winConfig
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		linuxConfig := expandAzureRmVirtualMachineOsProfileLinuxConfig(d)
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
		v = utils.Base64EncodeIfNot(v)
		profile.CustomData = &v
	}

	return profile, nil
}

func expandAzureRmVirtualMachineOsProfileSecrets(d *schema.ResourceData) *[]compute.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").([]interface{})
	secrets := make([]compute.VaultSecretGroup, 0, len(secretsConfig))

	for _, secretConfig := range secretsConfig {
		if secretConfig == nil {
			continue
		}

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
				if certConfig == nil {
					continue
				}
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

func expandAzureRmVirtualMachineOsProfileLinuxConfig(d *schema.ResourceData) *compute.LinuxConfiguration {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*schema.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]interface{})
	disablePasswordAuth := linuxConfig["disable_password_authentication"].(bool)

	config := &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &disablePasswordAuth,
	}

	linuxKeys := linuxConfig["ssh_keys"].([]interface{})
	sshPublicKeys := make([]compute.SSHPublicKey, 0)
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

	return config
}

func expandAzureRmVirtualMachineOsProfileWindowsConfig(d *schema.ResourceData) *compute.WindowsConfiguration {
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

	if v := osProfileConfig["timezone"]; v != nil && v.(string) != "" {
		config.TimeZone = utils.String(v.(string))
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

func expandAzureRmVirtualMachineAdditionalCapabilities(d *schema.ResourceData) *compute.AdditionalCapabilities {
	additionalCapabilities := d.Get("additional_capabilities").([]interface{})
	if len(additionalCapabilities) == 0 {
		return nil
	}

	additionalCapability := additionalCapabilities[0].(map[string]interface{})
	capability := &compute.AdditionalCapabilities{
		UltraSSDEnabled: utils.Bool(additionalCapability["ultra_ssd_enabled"].(bool)),
	}

	return capability
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
		if nic != nil {
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

	// BEGIN: code to be removed after GH-13016 is merged
	if vhdURI != "" && managedDiskID != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
	}
	if vhdURI != "" && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
	}
	// END: code to be removed after GH-13016 is merged
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
		if v, ok := m["timezone"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
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

func resourceArmVirtualMachineGetManagedDiskInfo(d *schema.ResourceData, disk *compute.ManagedDiskParameters, meta interface{}) (*compute.Disk, error) {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if disk == nil || disk.ID == nil {
		return nil, nil
	}

	diskId := *disk.ID
	id, err := azure.ParseAzureResourceID(diskId)
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

func determineVirtualMachineIPAddress(ctx context.Context, meta interface{}, props *compute.VirtualMachineProperties) (string, error) {
	nicClient := meta.(*clients.Client).Network.InterfacesClient
	pipClient := meta.(*clients.Client).Network.PublicIPsClient

	if props == nil {
		return "", nil
	}

	var networkInterface *network.Interface

	if profile := props.NetworkProfile; profile != nil {
		if nicReferences := profile.NetworkInterfaces; nicReferences != nil {
			for _, nicReference := range *nicReferences {
				// pick out the primary if multiple NIC's are assigned
				if len(*nicReferences) > 1 {
					if nicReference.Primary == nil || !*nicReference.Primary {
						continue
					}
				}

				id, err := azure.ParseAzureResourceID(*nicReference.ID)
				if err != nil {
					return "", err
				}

				resourceGroup := id.ResourceGroup
				name := id.Path["networkInterfaces"]

				nic, err := nicClient.Get(ctx, resourceGroup, name, "")
				if err != nil {
					return "", fmt.Errorf("Error obtaining NIC %q (Resource Group %q): %+v", name, resourceGroup, err)
				}

				networkInterface = &nic
				break
			}
		}
	}

	if networkInterface == nil {
		return "", fmt.Errorf("A Network Interface wasn't found on the Virtual Machine")
	}

	if props := networkInterface.InterfacePropertiesFormat; props != nil {
		if configs := props.IPConfigurations; configs != nil {
			for _, config := range *configs {
				if config.PublicIPAddress != nil {
					id, err := azure.ParseAzureResourceID(*config.PublicIPAddress.ID)
					if err != nil {
						return "", err
					}

					resourceGroup := id.ResourceGroup
					name := id.Path["publicIPAddresses"]

					pip, err := pipClient.Get(ctx, resourceGroup, name, "")
					if err != nil {
						return "", fmt.Errorf("Error obtaining Public IP %q (Resource Group %q): %+v", name, resourceGroup, err)
					}

					if pipProps := pip.PublicIPAddressPropertiesFormat; pipProps != nil {
						if ip := pipProps.IPAddress; ip != nil {
							return *ip, nil
						}
					}
				}

				if ip := config.PrivateIPAddress; ip != nil {
					return *ip, nil
				}
			}
		}
	}

	return "", fmt.Errorf("No Public or Private IP Address found on the Primary Network Interface")
}
