// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package legacy

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	compute2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	intStor "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func userDataDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
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
//
//	`azurerm_windows_virtual_machine` resources - as such this resource is feature-frozen and new
//	functionality will be added to these new resources instead.
func resourceVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineCreateUpdate,
		Read:   resourceVirtualMachineRead,
		Update: resourceVirtualMachineCreateUpdate,
		Delete: resourceVirtualMachineDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachines.ParseVirtualMachineID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"zones": {
				// @tombuildsstuff: since this is the legacy VM resource this is intentionally not using commonschema for consistency
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"plan": {
				Type:     pluginsdk.TypeList,
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

			"availability_set_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				StateFunc: func(id interface{}) string {
					return strings.ToLower(id.(string))
				},
				ConflictsWith: []string{"zones"},
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

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Windows_Client",
					"Windows_Server",
				}, false),
			},

			"vm_size": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			// lintignore:S018
			"storage_image_reference": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"publisher": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"offer": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"sku": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceVirtualMachineStorageImageReferenceHash,
			},

			"storage_os_disk": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"os_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualmachines.OperatingSystemTypesLinux),
								string(virtualmachines.OperatingSystemTypesWindows),
							}, false),
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"vhd_uri": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							ConflictsWith: []string{
								"storage_os_disk.0.managed_disk_id",
								"storage_os_disk.0.managed_disk_type",
							},
						},

						"managed_disk_id": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_os_disk.0.vhd_uri"},
						},

						"managed_disk_type": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_os_disk.0.vhd_uri"},
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualmachines.StorageAccountTypesPremiumLRS),
								string(virtualmachines.StorageAccountTypesStandardLRS),
								string(virtualmachines.StorageAccountTypesStandardSSDLRS),
							}, false),
						},

						"image_uri": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"caching": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"create_option": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.DiskSizeGB,
						},

						"write_accelerator_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"delete_os_disk_on_termination": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"storage_data_disk": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"vhd_uri": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"managed_disk_id": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"managed_disk_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualmachines.StorageAccountTypesPremiumLRS),
								string(virtualmachines.StorageAccountTypesStandardLRS),
								string(virtualmachines.StorageAccountTypesStandardSSDLRS),
								string(virtualmachines.StorageAccountTypesUltraSSDLRS),
							}, false),
						},

						"create_option": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
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
							ValidateFunc: validate.DiskSizeGB,
						},

						"lun": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"write_accelerator_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"delete_data_disks_on_termination": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"boot_diagnostics": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},

						"storage_uri": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"additional_capabilities": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ultra_ssd_enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			// lintignore:S018
			"os_profile": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"computer_name": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Required: true,
						},

						"admin_username": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"admin_password": {
							Type:      pluginsdk.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"custom_data": {
							Type:      pluginsdk.TypeString,
							ForceNew:  true,
							Optional:  true,
							Computed:  true,
							StateFunc: userDataStateFunc,
						},
					},
				},
				Set: resourceVirtualMachineStorageOsProfileHash,
			},

			// lintignore:S018
			"os_profile_windows_config": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"provision_vm_agent": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"enable_automatic_upgrades": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"timezone": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.VirtualMachineTimeZoneCaseInsensitive(),
						},
						"winrm": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"protocol": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"HTTP",
											"HTTPS",
										}, false),
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
										ValidateFunc: validation.StringInSlice([]string{
											"oobeSystem",
										}, false),
									},
									"component": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Microsoft-Windows-Shell-Setup",
										}, false),
									},
									"setting_name": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"AutoLogon",
											"FirstLogonCommands",
										}, false),
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
				Set:           resourceVirtualMachineStorageOsProfileWindowsConfigHash,
				ConflictsWith: []string{"os_profile_linux_config"},
			},

			// lintignore:S018
			"os_profile_linux_config": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disable_password_authentication": {
							Type:     pluginsdk.TypeBool,
							Required: true,
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
										Required: true,
									},
								},
							},
						},
					},
				},
				Set:           resourceVirtualMachineStorageOsProfileLinuxConfigHash,
				ConflictsWith: []string{"os_profile_windows_config"},
			},

			"os_profile_secrets": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"source_vault_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
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

			"network_interface_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"primary_network_interface_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine creation.")
	id := virtualmachines.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_machine", id.ID())
		}
	}

	osDisk, err := expandAzureRmVirtualMachineOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile := virtualmachines.StorageProfile{
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
	properties := virtualmachines.VirtualMachineProperties{
		NetworkProfile: &networkProfile,
		HardwareProfile: &virtualmachines.HardwareProfile{
			VMSize: pointer.To(virtualmachines.VirtualMachineSizeTypes(vmSize)),
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
		availSet := virtualmachines.SubResource{
			Id: &availabilitySet,
		}

		properties.AvailabilitySet = &availSet
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		properties.ProximityPlacementGroup = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	vm := virtualmachines.VirtualMachine{
		Name:       &id.VirtualMachineName,
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &properties,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
		Zones:      expandZones(d.Get("zones").([]interface{})),
	}

	if _, ok := d.GetOk("identity"); ok {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		vm.Identity = identityExpanded
	}

	if _, ok := d.GetOk("plan"); ok {
		vm.Plan = expandAzureRmVirtualMachinePlan(d)
	}

	locks.ByName(id.VirtualMachineName, compute2.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, compute2.VirtualMachineResourceName)

	if err := client.CreateOrUpdateThenPoll(ctx, id, vm, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	read, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return err
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if read.Model.Id == nil {
		return fmt.Errorf("retrieving %s: `id` was nil", id)
	}

	d.SetId(id.ID())

	ipAddress, err := determineVirtualMachineIPAddress(ctx, meta, read.Model.Properties)
	if err != nil {
		return fmt.Errorf("determining IP Address for %s: %+v", id, err)
	}

	provisionerType := "ssh"
	if props := read.Model.Properties; props != nil {
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

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.VirtualMachineName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("zones", model.Zones)
		d.Set("location", location.Normalize(model.Location))

		if err := d.Set("plan", flattenAzureRmVirtualMachinePlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %#v", err)
		}

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if availabilitySet := props.AvailabilitySet; availabilitySet != nil {
				// Lowercase due to incorrect capitalisation of resource group name in
				// availability set ID in response from get VM API request
				// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
				d.Set("availability_set_id", strings.ToLower(*availabilitySet.Id))
			}

			if proximityPlacementGroup := props.ProximityPlacementGroup; proximityPlacementGroup != nil {
				d.Set("proximity_placement_group_id", proximityPlacementGroup.Id)
			}

			if profile := props.HardwareProfile; profile != nil {
				d.Set("vm_size", pointer.From(profile.VMSize))
			}

			if profile := props.StorageProfile; profile != nil {
				if err := d.Set("storage_image_reference", pluginsdk.NewSet(resourceVirtualMachineStorageImageReferenceHash, flattenAzureRmVirtualMachineImageReference(profile.ImageReference))); err != nil {
					return fmt.Errorf("error setting Virtual Machine Storage Image Reference error: %#v", err)
				}

				if osDisk := profile.OsDisk; osDisk != nil {
					diskInfo, err := resourceVirtualMachineGetManagedDiskInfo(d, osDisk.ManagedDisk, meta)
					if err != nil {
						return fmt.Errorf("flattening `storage_os_disk`: %#v", err)
					}
					if err := d.Set("storage_os_disk", flattenAzureRmVirtualMachineOsDisk(osDisk, diskInfo)); err != nil {
						return fmt.Errorf("setting `storage_os_disk`: %#v", err)
					}
				}

				if dataDisks := profile.DataDisks; dataDisks != nil {
					disksInfo := make([]*disks.Disk, len(*dataDisks))
					for i, dataDisk := range *dataDisks {
						diskInfo, err := resourceVirtualMachineGetManagedDiskInfo(d, dataDisk.ManagedDisk, meta)
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
				if err := d.Set("os_profile", pluginsdk.NewSet(resourceVirtualMachineStorageOsProfileHash, flattenAzureRmVirtualMachineOsProfile(profile))); err != nil {
					return fmt.Errorf("setting `os_profile`: %#v", err)
				}

				if err := d.Set("os_profile_linux_config", pluginsdk.NewSet(resourceVirtualMachineStorageOsProfileLinuxConfigHash, flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(profile.LinuxConfiguration))); err != nil {
					return fmt.Errorf("setting `os_profile_linux_config`: %+v", err)
				}

				if err := d.Set("os_profile_windows_config", pluginsdk.NewSet(resourceVirtualMachineStorageOsProfileWindowsConfigHash, flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(profile.WindowsConfiguration))); err != nil {
					return fmt.Errorf("setting `os_profile_windows_config`: %+v", err)
				}

				if err := d.Set("os_profile_secrets", flattenAzureRmVirtualMachineOsProfileSecrets(profile.Secrets)); err != nil {
					return fmt.Errorf("setting `os_profile_secrets`: %+v", err)
				}
			}

			if profile := props.DiagnosticsProfile; profile != nil {
				if err := d.Set("boot_diagnostics", flattenAzureRmVirtualMachineDiagnosticsProfile(profile.BootDiagnostics)); err != nil {
					return fmt.Errorf("setting `boot_diagnostics`: %#v", err)
				}
			}
			if err := d.Set("additional_capabilities", flattenAzureRmVirtualMachineAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
				return fmt.Errorf("setting `additional_capabilities`: %#v", err)
			}

			if profile := props.NetworkProfile; profile != nil {
				if err := d.Set("network_interface_ids", flattenAzureRmVirtualMachineNetworkInterfaces(profile)); err != nil {
					return fmt.Errorf("flattening `network_interface_ids`: %#v", err)
				}

				if profile.NetworkInterfaces != nil {
					for _, nic := range *profile.NetworkInterfaces {
						if props := nic.Properties; props != nil {
							if props.Primary != nil && *props.Primary {
								d.Set("primary_network_interface_id", nic.Id)
								break
							}
						}
					}
				}
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, compute2.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, compute2.VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(virtualMachine.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(virtualMachine.HttpResponse) {
		// @tombuildsstuff: sending `nil` here omits this value from being sent - which matches
		// the previous behaviour - we're only splitting this out so it's clear why
		opts := virtualmachines.DefaultDeleteOperationOptions()
		opts.ForceDeletion = nil
		if err := client.DeleteThenPoll(ctx, *id, opts); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	// delete OS Disk if opted in
	deleteOsDisk := d.Get("delete_os_disk_on_termination").(bool)
	deleteDataDisks := d.Get("delete_data_disks_on_termination").(bool)

	if deleteOsDisk || deleteDataDisks {
		storageClient := meta.(*clients.Client).Storage

		model := virtualMachine.Model
		if model == nil {
			return fmt.Errorf("deleting Disks for %s - `model` was nil", id)
		}
		props := model.Properties
		if props == nil {
			return fmt.Errorf("deleting Disks for %s - `props` was nil", id)
		}
		storageProfile := props.StorageProfile
		if storageProfile == nil {
			return fmt.Errorf("deleting Disks for %s - `storageProfile` was nil", id)
		}

		if deleteOsDisk {
			log.Printf("[INFO] delete_os_disk_on_termination is enabled, deleting disk from %s", id)
			osDisk := storageProfile.OsDisk
			if osDisk == nil {
				return fmt.Errorf("deleting OS Disk for %s - `osDisk` was nil", id)
			}
			if osDisk.Vhd == nil && osDisk.ManagedDisk == nil {
				return fmt.Errorf("unable to determine OS Disk Type to Delete it for %s", id)
			}

			if osDisk.Vhd != nil {
				if err = resourceVirtualMachineDeleteVhd(ctx, storageClient, id.SubscriptionId, osDisk.Vhd); err != nil {
					return fmt.Errorf("deleting OS Disk VHD: %+v", err)
				}
			} else if osDisk.ManagedDisk != nil {
				if err = resourceVirtualMachineDeleteManagedDisk(d, osDisk.ManagedDisk, meta); err != nil {
					return fmt.Errorf("deleting OS Managed Disk: %+v", err)
				}
			}
		}

		// delete Data disks if opted in
		if deleteDataDisks {
			log.Printf("[INFO] delete_data_disks_on_termination is enabled, deleting each data disk from %s", id)

			dataDisks := storageProfile.DataDisks
			if dataDisks == nil {
				return fmt.Errorf("deleting Data Disks for %s: `dataDisks` was nil", id)
			}

			for _, disk := range *dataDisks {
				if disk.Vhd == nil && disk.ManagedDisk == nil {
					return fmt.Errorf("unable to determine Data Disk Type to Delete it for %s / Disk %q", id, *disk.Name)
				}

				if disk.Vhd != nil {
					if err = resourceVirtualMachineDeleteVhd(ctx, storageClient, id.SubscriptionId, disk.Vhd); err != nil {
						return fmt.Errorf("deleting Data Disk VHD: %+v", err)
					}
				} else if disk.ManagedDisk != nil {
					if err = resourceVirtualMachineDeleteManagedDisk(d, disk.ManagedDisk, meta); err != nil {
						return fmt.Errorf("deleting Data Managed Disk: %+v", err)
					}
				}
			}
		}
	}

	return nil
}

func resourceVirtualMachineDeleteVhd(ctx context.Context, storageClient *intStor.Client, subscriptionId string, vhd *virtualmachines.VirtualHardDisk) error {
	if vhd == nil {
		return fmt.Errorf("`vhd` was nil`")
	}
	if vhd.Uri == nil {
		return fmt.Errorf("`vhd.Uri` was nil`")
	}

	uri := *vhd.Uri
	id, err := blobs.ParseBlobID(uri, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing %q: %s", uri, err)
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %s", id.AccountId.AccountName, id.BlobName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q (Disk %q)", id.AccountId.AccountName, uri)
	}

	if err != nil {
		return fmt.Errorf("building Blobs Client: %s", err)
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %s", err)
	}

	input := blobs.DeleteInput{
		DeleteSnapshots: false,
	}
	if _, err := blobsClient.Delete(ctx, id.ContainerName, id.BlobName, input); err != nil {
		return fmt.Errorf("deleting Blob %q (Container %q in %s): %+v", id.BlobName, id.ContainerName, id.AccountId, err)
	}

	return nil
}

func resourceVirtualMachineDeleteManagedDisk(d *pluginsdk.ResourceData, disk *virtualmachines.ManagedDiskParameters, meta interface{}) error {
	if disk == nil {
		return fmt.Errorf("`disk` was nil`")
	}
	if disk.Id == nil {
		return fmt.Errorf("`disk.Id` was nil`")
	}
	managedDiskID := *disk.Id

	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseManagedDiskIDInsensitively(managedDiskID)
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting Managed Disk %s: %+v", *id, err)
	}

	return nil
}

func flattenAzureRmVirtualMachinePlan(plan *virtualmachines.Plan) []interface{} {
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

func flattenAzureRmVirtualMachineImageReference(image *virtualmachines.ImageReference) []interface{} {
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
	if image.Id != nil {
		result["id"] = *image.Id
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineDiagnosticsProfile(profile *virtualmachines.BootDiagnostics) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if profile.Enabled != nil {
		result["enabled"] = *profile.Enabled
	}

	if profile.StorageUri != nil {
		result["storage_uri"] = *profile.StorageUri
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineAdditionalCapabilities(profile *virtualmachines.AdditionalCapabilities) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	if v := profile.UltraSSDEnabled; v != nil {
		result["ultra_ssd_enabled"] = *v
	}
	return []interface{}{result}
}

func flattenAzureRmVirtualMachineNetworkInterfaces(profile *virtualmachines.NetworkProfile) []interface{} {
	result := make([]interface{}, 0)
	for _, nic := range *profile.NetworkInterfaces {
		result = append(result, *nic.Id)
	}
	return result
}

func flattenAzureRmVirtualMachineOsProfileSecrets(secrets *[]virtualmachines.VaultSecretGroup) []interface{} {
	if secrets == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, secret := range *secrets {
		s := map[string]interface{}{
			"source_vault_id": *secret.SourceVault.Id,
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

func flattenAzureRmVirtualMachineDataDisk(disks *[]virtualmachines.DataDisk, disksInfo []*disks.Disk) interface{} {
	result := make([]interface{}, len(*disks))
	for i, disk := range *disks {
		l := make(map[string]interface{})
		l["name"] = *disk.Name
		if disk.Vhd != nil {
			l["vhd_uri"] = *disk.Vhd.Uri
		}
		if disk.ManagedDisk != nil {
			l["managed_disk_type"] = string(pointer.From(disk.ManagedDisk.StorageAccountType))
			if disk.ManagedDisk.Id != nil {
				l["managed_disk_id"] = *disk.ManagedDisk.Id
			}
		}
		l["create_option"] = disk.CreateOption
		l["caching"] = string(pointer.From(disk.Caching))
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}

		l["lun"] = disk.Lun

		if v := disk.WriteAcceleratorEnabled; v != nil {
			l["write_accelerator_enabled"] = *disk.WriteAcceleratorEnabled
		}

		flattenAzureRmVirtualMachineReviseDiskInfo(l, disksInfo[i])

		result[i] = l
	}
	return result
}

func flattenAzureRmVirtualMachineOsProfile(input *virtualmachines.OSProfile) []interface{} {
	result := make(map[string]interface{})
	result["computer_name"] = pointer.From(input.ComputerName)
	result["admin_username"] = pointer.From(input.AdminUsername)
	if input.CustomData != nil {
		result["custom_data"] = *input.CustomData
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(config *virtualmachines.WindowsConfiguration) []interface{} {
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
			listener["protocol"] = string(pointer.From(i.Protocol))

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
			c["pass"] = string(pointer.From(i.PassName))
			c["component"] = string(pointer.From(i.ComponentName))
			c["setting_name"] = string(pointer.From(i.SettingName))

			if i.Content != nil {
				c["content"] = *i.Content
			}

			content = append(content, c)
		}
	}
	result["additional_unattend_config"] = content

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(config *virtualmachines.LinuxConfiguration) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if config.DisablePasswordAuthentication != nil {
		result["disable_password_authentication"] = *config.DisablePasswordAuthentication
	}

	if config.Ssh != nil && config.Ssh.PublicKeys != nil && len(*config.Ssh.PublicKeys) > 0 {
		ssh_keys := make([]map[string]interface{}, 0)
		for _, i := range *config.Ssh.PublicKeys {
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

func flattenAzureRmVirtualMachineOsDisk(disk *virtualmachines.OSDisk, diskInfo *disks.Disk) []interface{} {
	result := make(map[string]interface{})
	if disk.Name != nil {
		result["name"] = *disk.Name
	}
	if disk.Vhd != nil && disk.Vhd.Uri != nil {
		result["vhd_uri"] = *disk.Vhd.Uri
	}
	if disk.Image != nil && disk.Image.Uri != nil {
		result["image_uri"] = *disk.Image.Uri
	}
	if disk.ManagedDisk != nil {
		result["managed_disk_type"] = string(pointer.From(disk.ManagedDisk.StorageAccountType))
		if disk.ManagedDisk.Id != nil {
			result["managed_disk_id"] = *disk.ManagedDisk.Id
		}
	}
	result["create_option"] = disk.CreateOption
	result["caching"] = disk.Caching
	if disk.DiskSizeGB != nil {
		result["disk_size_gb"] = *disk.DiskSizeGB
	}
	result["os_type"] = string(pointer.From(disk.OsType))

	if v := disk.WriteAcceleratorEnabled; v != nil {
		result["write_accelerator_enabled"] = *disk.WriteAcceleratorEnabled
	}

	flattenAzureRmVirtualMachineReviseDiskInfo(result, diskInfo)

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineReviseDiskInfo(result map[string]interface{}, diskInfo *disks.Disk) {
	if diskInfo != nil {
		if diskInfo.Sku != nil {
			result["managed_disk_type"] = string(pointer.From(diskInfo.Sku.Name))
		}
		if diskInfo.Properties != nil && diskInfo.Properties.DiskSizeGB != nil {
			result["disk_size_gb"] = *diskInfo.Properties.DiskSizeGB
		}
	}
}

func expandAzureRmVirtualMachinePlan(d *pluginsdk.ResourceData) *virtualmachines.Plan {
	planConfigs := d.Get("plan").([]interface{})
	if len(planConfigs) == 0 {
		return nil
	}

	planConfig := planConfigs[0].(map[string]interface{})

	publisher := planConfig["publisher"].(string)
	name := planConfig["name"].(string)
	product := planConfig["product"].(string)

	return &virtualmachines.Plan{
		Publisher: &publisher,
		Name:      &name,
		Product:   &product,
	}
}

func expandAzureRmVirtualMachineOsProfile(d *pluginsdk.ResourceData) (*virtualmachines.OSProfile, error) {
	osProfiles := d.Get("os_profile").(*pluginsdk.Set).List()

	osProfile := osProfiles[0].(map[string]interface{})

	adminUsername := osProfile["admin_username"].(string)
	adminPassword := osProfile["admin_password"].(string)
	computerName := osProfile["computer_name"].(string)

	profile := &virtualmachines.OSProfile{
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
		return nil, fmt.Errorf("either a `os_profile_linux_config` or a `os_profile_windows_config` must be specified")
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

func expandAzureRmVirtualMachineOsProfileSecrets(d *pluginsdk.ResourceData) *[]virtualmachines.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").([]interface{})
	secrets := make([]virtualmachines.VaultSecretGroup, 0, len(secretsConfig))

	for _, secretConfig := range secretsConfig {
		if secretConfig == nil {
			continue
		}

		config := secretConfig.(map[string]interface{})
		sourceVaultId := config["source_vault_id"].(string)

		vaultSecretGroup := virtualmachines.VaultSecretGroup{
			SourceVault: &virtualmachines.SubResource{
				Id: &sourceVaultId,
			},
		}

		if v := config["vault_certificates"]; v != nil {
			certsConfig := v.([]interface{})
			certs := make([]virtualmachines.VaultCertificate, 0, len(certsConfig))
			for _, certConfig := range certsConfig {
				if certConfig == nil {
					continue
				}
				config := certConfig.(map[string]interface{})

				certUrl := config["certificate_url"].(string)
				cert := virtualmachines.VaultCertificate{
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

func expandAzureRmVirtualMachineOsProfileLinuxConfig(d *pluginsdk.ResourceData) *virtualmachines.LinuxConfiguration {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*pluginsdk.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]interface{})
	disablePasswordAuth := linuxConfig["disable_password_authentication"].(bool)

	config := &virtualmachines.LinuxConfiguration{
		DisablePasswordAuthentication: &disablePasswordAuth,
	}

	linuxKeys := linuxConfig["ssh_keys"].([]interface{})
	sshPublicKeys := make([]virtualmachines.SshPublicKey, 0)
	for _, key := range linuxKeys {
		sshKey, ok := key.(map[string]interface{})
		if !ok {
			continue
		}
		path := sshKey["path"].(string)
		keyData := sshKey["key_data"].(string)

		sshPublicKey := virtualmachines.SshPublicKey{
			Path:    &path,
			KeyData: &keyData,
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	if len(sshPublicKeys) > 0 {
		config.Ssh = &virtualmachines.SshConfiguration{
			PublicKeys: &sshPublicKeys,
		}
	}

	return config
}

func expandAzureRmVirtualMachineOsProfileWindowsConfig(d *pluginsdk.ResourceData) *virtualmachines.WindowsConfiguration {
	osProfilesWindowsConfig := d.Get("os_profile_windows_config").(*pluginsdk.Set).List()

	osProfileConfig := osProfilesWindowsConfig[0].(map[string]interface{})
	config := &virtualmachines.WindowsConfiguration{}

	if v := osProfileConfig["provision_vm_agent"]; v != nil {
		provision := v.(bool)
		config.ProvisionVMAgent = &provision
	}

	if v := osProfileConfig["enable_automatic_upgrades"]; v != nil {
		update := v.(bool)
		config.EnableAutomaticUpdates = &update
	}

	if v := osProfileConfig["timezone"]; v != nil && v.(string) != "" {
		config.TimeZone = pointer.To(v.(string))
	}

	if v := osProfileConfig["winrm"]; v != nil {
		winRm := v.([]interface{})
		if len(winRm) > 0 {
			winRmListeners := make([]virtualmachines.WinRMListener, 0, len(winRm))
			for _, winRmConfig := range winRm {
				config := winRmConfig.(map[string]interface{})

				protocol := config["protocol"].(string)
				winRmListener := virtualmachines.WinRMListener{
					Protocol: pointer.To(virtualmachines.ProtocolTypes(protocol)),
				}
				if v := config["certificate_url"].(string); v != "" {
					winRmListener.CertificateURL = &v
				}

				winRmListeners = append(winRmListeners, winRmListener)
			}
			config.WinRM = &virtualmachines.WinRMConfiguration{
				Listeners: &winRmListeners,
			}
		}
	}
	if v := osProfileConfig["additional_unattend_config"]; v != nil {
		additionalConfig := v.([]interface{})
		if len(additionalConfig) > 0 {
			additionalConfigContent := make([]virtualmachines.AdditionalUnattendContent, 0, len(additionalConfig))
			for _, addConfig := range additionalConfig {
				config := addConfig.(map[string]interface{})
				pass := config["pass"].(string)
				component := config["component"].(string)
				settingName := config["setting_name"].(string)
				content := config["content"].(string)

				addContent := virtualmachines.AdditionalUnattendContent{
					PassName:      pointer.To(virtualmachines.PassNames(pass)),
					ComponentName: pointer.To(virtualmachines.ComponentNames(component)),
					SettingName:   pointer.To(virtualmachines.SettingNames(settingName)),
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

func expandAzureRmVirtualMachineDataDisk(d *pluginsdk.ResourceData) ([]virtualmachines.DataDisk, error) {
	disks := d.Get("storage_data_disk").([]interface{})
	data_disks := make([]virtualmachines.DataDisk, 0, len(disks))
	for _, disk_config := range disks {
		config := disk_config.(map[string]interface{})

		name := config["name"].(string)
		createOption := config["create_option"].(string)
		vhdURI := config["vhd_uri"].(string)
		managedDiskType := config["managed_disk_type"].(string)
		managedDiskID := config["managed_disk_id"].(string)
		lun := int64(config["lun"].(int))

		data_disk := virtualmachines.DataDisk{
			Name:         &name,
			Lun:          lun,
			CreateOption: virtualmachines.DiskCreateOptionTypes(createOption),
		}

		if vhdURI != "" {
			data_disk.Vhd = &virtualmachines.VirtualHardDisk{
				Uri: &vhdURI,
			}
		}

		managedDisk := &virtualmachines.ManagedDiskParameters{}

		if managedDiskType != "" {
			managedDisk.StorageAccountType = pointer.To(virtualmachines.StorageAccountTypes(managedDiskType))
			data_disk.ManagedDisk = managedDisk
		}

		if managedDiskID != "" {
			managedDisk.Id = &managedDiskID
			data_disk.ManagedDisk = managedDisk
		}

		if vhdURI != "" && managedDiskID != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
		}
		if vhdURI != "" && managedDiskType != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
		}
		if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(data_disk.CreateOption), string(virtualmachines.DiskCreateOptionTypesAttach)) {
			return nil, fmt.Errorf("[ERROR] Must specify `vhd_uri` or `managed_disk_id` to attach")
		}

		if v := config["caching"].(string); v != "" {
			data_disk.Caching = pointer.To(virtualmachines.CachingTypes(v))
		}

		if v, ok := config["disk_size_gb"].(int); ok {
			data_disk.DiskSizeGB = pointer.To(int64(v))
		}

		if v, ok := config["write_accelerator_enabled"].(bool); ok {
			data_disk.WriteAcceleratorEnabled = utils.Bool(v)
		}

		data_disks = append(data_disks, data_disk)
	}

	return data_disks, nil
}

func expandAzureRmVirtualMachineDiagnosticsProfile(d *pluginsdk.ResourceData) *virtualmachines.DiagnosticsProfile {
	bootDiagnostics := d.Get("boot_diagnostics").([]interface{})

	diagnosticsProfile := &virtualmachines.DiagnosticsProfile{}
	if len(bootDiagnostics) > 0 {
		bootDiagnostic := bootDiagnostics[0].(map[string]interface{})

		diagnostic := &virtualmachines.BootDiagnostics{
			Enabled:    utils.Bool(bootDiagnostic["enabled"].(bool)),
			StorageUri: pointer.To(bootDiagnostic["storage_uri"].(string)),
		}

		diagnosticsProfile.BootDiagnostics = diagnostic

		return diagnosticsProfile
	}

	return nil
}

func expandAzureRmVirtualMachineAdditionalCapabilities(d *pluginsdk.ResourceData) *virtualmachines.AdditionalCapabilities {
	additionalCapabilities := d.Get("additional_capabilities").([]interface{})
	if len(additionalCapabilities) == 0 || additionalCapabilities[0] == nil {
		return nil
	}

	additionalCapability := additionalCapabilities[0].(map[string]interface{})
	capability := &virtualmachines.AdditionalCapabilities{
		UltraSSDEnabled: utils.Bool(additionalCapability["ultra_ssd_enabled"].(bool)),
	}

	return capability
}

func expandAzureRmVirtualMachineImageReference(d *pluginsdk.ResourceData) (*virtualmachines.ImageReference, error) {
	storageImageRefs := d.Get("storage_image_reference").(*pluginsdk.Set).List()

	storageImageRef := storageImageRefs[0].(map[string]interface{})
	imageID := storageImageRef["id"].(string)
	publisher := storageImageRef["publisher"].(string)

	imageReference := virtualmachines.ImageReference{}

	if imageID != "" && publisher != "" {
		return nil, fmt.Errorf("conflict between `id` and `publisher` (only one or the other can be used)")
	}

	if imageID != "" {
		imageReference.Id = pointer.To(storageImageRef["id"].(string))
	} else {
		offer := storageImageRef["offer"].(string)
		sku := storageImageRef["sku"].(string)
		version := storageImageRef["version"].(string)

		imageReference = virtualmachines.ImageReference{
			Publisher: &publisher,
			Offer:     &offer,
			Sku:       &sku,
			Version:   &version,
		}
	}

	return &imageReference, nil
}

func expandAzureRmVirtualMachineNetworkProfile(d *pluginsdk.ResourceData) virtualmachines.NetworkProfile {
	nicIds := d.Get("network_interface_ids").([]interface{})
	primaryNicId := d.Get("primary_network_interface_id").(string)
	network_interfaces := make([]virtualmachines.NetworkInterfaceReference, 0, len(nicIds))

	network_profile := virtualmachines.NetworkProfile{}

	for _, nic := range nicIds {
		if nic != nil {
			id := nic.(string)
			primary := id == primaryNicId

			network_interface := virtualmachines.NetworkInterfaceReference{
				Id: &id,
				Properties: &virtualmachines.NetworkInterfaceReferenceProperties{
					Primary: &primary,
				},
			}
			network_interfaces = append(network_interfaces, network_interface)
		}
	}

	network_profile.NetworkInterfaces = &network_interfaces

	return network_profile
}

func expandAzureRmVirtualMachineOsDisk(d *pluginsdk.ResourceData) (*virtualmachines.OSDisk, error) {
	disks := d.Get("storage_os_disk").([]interface{})

	config := disks[0].(map[string]interface{})

	name := config["name"].(string)
	imageURI := config["image_uri"].(string)
	createOption := config["create_option"].(string)
	vhdURI := config["vhd_uri"].(string)
	managedDiskType := config["managed_disk_type"].(string)
	managedDiskID := config["managed_disk_id"].(string)

	osDisk := &virtualmachines.OSDisk{
		Name:         &name,
		CreateOption: virtualmachines.DiskCreateOptionTypes(createOption),
	}

	if vhdURI != "" {
		osDisk.Vhd = &virtualmachines.VirtualHardDisk{
			Uri: &vhdURI,
		}
	}

	managedDisk := &virtualmachines.ManagedDiskParameters{}

	if managedDiskType != "" {
		managedDisk.StorageAccountType = pointer.To(virtualmachines.StorageAccountTypes(managedDiskType))
		osDisk.ManagedDisk = managedDisk
	}

	if managedDiskID != "" {
		managedDisk.Id = &managedDiskID
		osDisk.ManagedDisk = managedDisk
	}

	// BEGIN: code to be removed after GH-13016 is merged
	if vhdURI != "" && managedDiskID != "" {
		return nil, fmt.Errorf("conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
	}
	if vhdURI != "" && managedDiskType != "" {
		return nil, fmt.Errorf("conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
	}
	// END: code to be removed after GH-13016 is merged
	if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(osDisk.CreateOption), string(virtualmachines.DiskCreateOptionTypesAttach)) {
		return nil, fmt.Errorf("must specify `vhd_uri` or `managed_disk_id` to attach")
	}

	if v := config["image_uri"].(string); v != "" {
		osDisk.Image = &virtualmachines.VirtualHardDisk{
			Uri: &imageURI,
		}
	}

	if v := config["os_type"].(string); v != "" {
		osDisk.OsType = pointer.To(virtualmachines.OperatingSystemTypes(v))
	}

	if v := config["caching"].(string); v != "" {
		osDisk.Caching = pointer.To(virtualmachines.CachingTypes(v))
	}

	if v := config["disk_size_gb"].(int); v != 0 {
		osDisk.DiskSizeGB = pointer.To(int64(v))
	}

	if v, ok := config["write_accelerator_enabled"].(bool); ok {
		osDisk.WriteAcceleratorEnabled = utils.Bool(v)
	}

	return osDisk, nil
}

func resourceVirtualMachineStorageOsProfileHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["admin_username"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["computer_name"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineStorageOsProfileWindowsConfigHash(v interface{}) int {
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

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineStorageOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineStorageImageReferenceHash(v interface{}) int {
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

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineGetManagedDiskInfo(d *pluginsdk.ResourceData, disk *virtualmachines.ManagedDiskParameters, meta interface{}) (*disks.Disk, error) {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if disk == nil || disk.Id == nil {
		return nil, nil
	}

	diskId := *disk.Id
	id, err := commonids.ParseManagedDiskIDInsensitively(diskId)
	if err != nil {
		return nil, fmt.Errorf("parsing Disk ID %q: %+v", diskId, err)
	}

	diskResp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return diskResp.Model, nil
}

func determineVirtualMachineIPAddress(ctx context.Context, meta interface{}, props *virtualmachines.VirtualMachineProperties) (string, error) {
	nicClient := meta.(*clients.Client).Network.NetworkInterfaces
	pipClient := meta.(*clients.Client).Network.PublicIPAddresses

	if props == nil {
		return "", nil
	}

	var networkInterface *networkinterfaces.NetworkInterface

	if profile := props.NetworkProfile; profile != nil {
		if nicReferences := profile.NetworkInterfaces; nicReferences != nil {
			for _, nicReference := range *nicReferences {
				// pick out the primary if multiple NIC's are assigned
				if len(*nicReferences) > 1 {
					if nicReference.Properties == nil || nicReference.Properties.Primary == nil || !*nicReference.Properties.Primary {
						continue
					}
				}

				id, err := commonids.ParseNetworkInterfaceID(*nicReference.Id)
				if err != nil {
					return "", err
				}

				nic, err := nicClient.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
				if err != nil {
					return "", fmt.Errorf("retrieving %s: %+v", id, err)
				}

				networkInterface = nic.Model
				break
			}
		}
	}

	if networkInterface == nil {
		return "", fmt.Errorf("A Network Interface wasn't found on the Virtual Machine")
	}

	if props := networkInterface.Properties; props != nil {
		if configs := props.IPConfigurations; configs != nil {
			for _, config := range *configs {
				if configProps := config.Properties; configProps != nil {
					if configProps.PublicIPAddress != nil {
						id, err := commonids.ParsePublicIPAddressID(*configProps.PublicIPAddress.Id)
						if err != nil {
							return "", err
						}

						pip, err := pipClient.Get(ctx, *id, publicipaddresses.DefaultGetOperationOptions())
						if err != nil {
							return "", fmt.Errorf("retrieving %s: %+v", id, err)
						}

						if model := pip.Model; model != nil {
							if pipProps := model.Properties; pipProps != nil {
								if ip := pipProps.IPAddress; ip != nil {
									return *ip, nil
								}
							}
						}
					}

					if ip := configProps.PrivateIPAddress; ip != nil {
						return *ip, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("no Public or Private IP Address found on the Primary Network Interface")
}
