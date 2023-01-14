package legacy

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	compute2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	intStor "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
	"github.com/tombuildsstuff/kermit/sdk/compute/2022-08-01/compute"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-05-01/network"
	"golang.org/x/net/context"
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
			_, err := parse.VirtualMachineID(id)
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
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"identity_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: commonids.ValidateUserAssignedIdentityID,
							},
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

			"vm_size": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			//lintignore:S018
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
								string(compute.OperatingSystemTypesLinux),
								string(compute.OperatingSystemTypesWindows),
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
								string(compute.StorageAccountTypesPremiumLRS),
								string(compute.StorageAccountTypesStandardLRS),
								string(compute.StorageAccountTypesStandardSSDLRS),
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
								string(compute.StorageAccountTypesPremiumLRS),
								string(compute.StorageAccountTypesStandardLRS),
								string(compute.StorageAccountTypesStandardSSDLRS),
								string(compute.StorageAccountTypesUltraSSDLRS),
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

			//lintignore:S018
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

			//lintignore:S018
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

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Legacy.VMClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine creation.")
	id := parse.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_machine", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)
	zones := expandZones(d.Get("zones").([]interface{}))

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
		Name:                     &id.Name,
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

	locks.ByName(id.Name, compute2.VirtualMachineResourceName)
	defer locks.UnlockByName(id.Name, compute2.VirtualMachineResourceName)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, vm)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read %s", id)
	}

	d.SetId(id.ID())

	ipAddress, err := determineVirtualMachineIPAddress(ctx, meta, read.VirtualMachineProperties)
	if err != nil {
		return fmt.Errorf("determining IP Address for %s: %+v", id, err)
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

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	vmclient := meta.(*clients.Client).Legacy.VMClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := vmclient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Virtual Machine %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zones", resp.Zones)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("plan", flattenAzureRmVirtualMachinePlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %#v", err)
	}

	identity, err := flattenAzureRmVirtualMachineIdentity(resp.Identity)
	if err != nil {
		return err
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
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
			if err := d.Set("storage_image_reference", pluginsdk.NewSet(resourceVirtualMachineStorageImageReferenceHash, flattenAzureRmVirtualMachineImageReference(profile.ImageReference))); err != nil {
				return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Image Reference error: %#v", err)
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
				disksInfo := make([]*compute.Disk, len(*dataDisks))
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

			osProfile := flattenAzureRmVirtualMachineOsProfile(profile, id.Name)
			if err := d.Set("os_profile", pluginsdk.NewSet(resourceVirtualMachineStorageOsProfileHash, osProfile)); err != nil {
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

func resourceVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Legacy.VMClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, compute2.VirtualMachineResourceName)
	defer locks.UnlockByName(id.Name, compute2.VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving Virtual Machine %q : %s", id.String(), err)
	}

	// @tombuildsstuff: sending `nil` here omits this value from being sent - which matches
	// the previous behaviour - we're only splitting this out so it's clear why
	var forceDeletion *bool = nil
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, forceDeletion)
	if err != nil {
		return fmt.Errorf("deleting Virtual Machine %q : %s", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Virtual Machine %q : %s", id.String(), err)
	}

	// delete OS Disk if opted in
	deleteOsDisk := d.Get("delete_os_disk_on_termination").(bool)
	deleteDataDisks := d.Get("delete_data_disks_on_termination").(bool)

	if deleteOsDisk || deleteDataDisks {
		storageClient := meta.(*clients.Client).Storage

		props := virtualMachine.VirtualMachineProperties
		if props == nil {
			return fmt.Errorf("deleting Disks for Virtual Machine %q - `props` was nil", id.Name)
		}
		storageProfile := props.StorageProfile
		if storageProfile == nil {
			return fmt.Errorf("deleting Disks for Virtual Machine %q - `storageProfile` was nil", id.Name)
		}

		if deleteOsDisk {
			log.Printf("[INFO] delete_os_disk_on_termination is enabled, deleting disk from %s", id.Name)
			osDisk := storageProfile.OsDisk
			if osDisk == nil {
				return fmt.Errorf("deleting OS Disk for Virtual Machine %q - `osDisk` was nil", id.Name)
			}
			if osDisk.Vhd == nil && osDisk.ManagedDisk == nil {
				return fmt.Errorf("Unable to determine OS Disk Type to Delete it for Virtual Machine %q", id.Name)
			}

			if osDisk.Vhd != nil {
				if err = resourceVirtualMachineDeleteVhd(ctx, storageClient, osDisk.Vhd); err != nil {
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
			log.Printf("[INFO] delete_data_disks_on_termination is enabled, deleting each data disk from %q", id.Name)

			dataDisks := storageProfile.DataDisks
			if dataDisks == nil {
				return fmt.Errorf("deleting Data Disks for Virtual Machine %q: `dataDisks` was nil", id.Name)
			}

			for _, disk := range *dataDisks {
				if disk.Vhd == nil && disk.ManagedDisk == nil {
					return fmt.Errorf("Unable to determine Data Disk Type to Delete it for Virtual Machine %q / Disk %q", id.Name, *disk.Name)
				}

				if disk.Vhd != nil {
					if err = resourceVirtualMachineDeleteVhd(ctx, storageClient, disk.Vhd); err != nil {
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

func resourceVirtualMachineDeleteVhd(ctx context.Context, storageClient *intStor.Client, vhd *compute.VirtualHardDisk) error {
	if vhd == nil {
		return fmt.Errorf("`vhd` was nil`")
	}
	if vhd.URI == nil {
		return fmt.Errorf("`vhd.URI` was nil`")
	}

	uri := *vhd.URI
	id, err := blobs.ParseResourceID(uri)
	if err != nil {
		return fmt.Errorf("parsing %q: %s", uri, err)
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %s", id.AccountName, id.BlobName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q (Disk %q)!", id.AccountName, uri)
	}

	if err != nil {
		return fmt.Errorf("building Blobs Client: %s", err)
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Blobs Client: %s", err)
	}

	input := blobs.DeleteInput{
		DeleteSnapshots: false,
	}
	if _, err := blobsClient.Delete(ctx, id.AccountName, id.ContainerName, id.BlobName, input); err != nil {
		return fmt.Errorf("deleting Blob %q (Container %q / Account %q / Resource Group %q): %s", id.BlobName, id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func resourceVirtualMachineDeleteManagedDisk(d *pluginsdk.ResourceData, disk *compute.ManagedDiskParameters, meta interface{}) error {
	if disk == nil {
		return fmt.Errorf("`disk` was nil`")
	}
	if disk.ID == nil {
		return fmt.Errorf("`disk.ID` was nil`")
	}
	managedDiskID := *disk.ID

	client := meta.(*clients.Client).Legacy.DisksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disks.ParseDiskID(managedDiskID)
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.DiskName)
	if err != nil {
		return fmt.Errorf("deleting Managed Disk %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
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

func flattenAzureRmVirtualMachineIdentity(identity *compute.VirtualMachineIdentity) ([]interface{}, error) {
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
		/*
			"userAssignedIdentities": {
			  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tomdevidentity/providers/Microsoft.ManagedIdentity/userAssignedIdentities/tom123": {
				"principalId": "00000000-0000-0000-0000-000000000000",
				"clientId": "00000000-0000-0000-0000-000000000000"
			  }
			}
		*/
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

func flattenAzureRmVirtualMachineOsProfile(input *compute.OSProfile, defaultName string) []interface{} {
	result := make(map[string]interface{})

	if input.ComputerName != nil {
		result["computer_name"] = *input.ComputerName
	} else {
		result["computer_name"] = defaultName
	}

	if input.AdminUsername != nil {
		result["admin_username"] = *input.AdminUsername
	} else {
		result["admin_username"] = ""
	}

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

func expandAzureRmVirtualMachinePlan(d *pluginsdk.ResourceData) *compute.Plan {
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

func expandAzureRmVirtualMachineIdentity(d *pluginsdk.ResourceData) *compute.VirtualMachineIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := compute.ResourceIdentityType(identity["type"].(string))

	identityIds := make(map[string]*compute.UserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &compute.UserAssignedIdentitiesValue{}
	}

	vmIdentity := compute.VirtualMachineIdentity{
		Type: identityType,
	}

	if vmIdentity.Type == compute.ResourceIdentityTypeUserAssigned || vmIdentity.Type == compute.ResourceIdentityTypeSystemAssignedUserAssigned {
		vmIdentity.UserAssignedIdentities = identityIds
	}

	return &vmIdentity
}

func expandAzureRmVirtualMachineOsProfile(d *pluginsdk.ResourceData) (*compute.OSProfile, error) {
	osProfiles := d.Get("os_profile").(*pluginsdk.Set).List()

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

func expandAzureRmVirtualMachineOsProfileSecrets(d *pluginsdk.ResourceData) *[]compute.VaultSecretGroup {
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

func expandAzureRmVirtualMachineOsProfileLinuxConfig(d *pluginsdk.ResourceData) *compute.LinuxConfiguration {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*pluginsdk.Set).List()

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

func expandAzureRmVirtualMachineOsProfileWindowsConfig(d *pluginsdk.ResourceData) *compute.WindowsConfiguration {
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

func expandAzureRmVirtualMachineDataDisk(d *pluginsdk.ResourceData) ([]compute.DataDisk, error) {
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
		if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(data_disk.CreateOption), string(compute.DiskCreateOptionAttach)) {
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

func expandAzureRmVirtualMachineDiagnosticsProfile(d *pluginsdk.ResourceData) *compute.DiagnosticsProfile {
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

func expandAzureRmVirtualMachineAdditionalCapabilities(d *pluginsdk.ResourceData) *compute.AdditionalCapabilities {
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

func expandAzureRmVirtualMachineImageReference(d *pluginsdk.ResourceData) (*compute.ImageReference, error) {
	storageImageRefs := d.Get("storage_image_reference").(*pluginsdk.Set).List()

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

func expandAzureRmVirtualMachineNetworkProfile(d *pluginsdk.ResourceData) compute.NetworkProfile {
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

func expandAzureRmVirtualMachineOsDisk(d *pluginsdk.ResourceData) (*compute.OSDisk, error) {
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
	if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(osDisk.CreateOption), string(compute.DiskCreateOptionAttach)) {
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

func resourceVirtualMachineGetManagedDiskInfo(d *pluginsdk.ResourceData, disk *compute.ManagedDiskParameters, meta interface{}) (*compute.Disk, error) {
	client := meta.(*clients.Client).Legacy.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if disk == nil || disk.ID == nil {
		return nil, nil
	}

	diskId := *disk.ID
	id, err := disks.ParseDiskID(diskId)
	if err != nil {
		return nil, fmt.Errorf("parsing Disk ID %q: %+v", diskId, err)
	}

	diskResp, err := client.Get(ctx, id.ResourceGroupName, id.DiskName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
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

				id, err := networkParse.NetworkInterfaceID(*nicReference.ID)
				if err != nil {
					return "", err
				}

				nic, err := nicClient.Get(ctx, id.ResourceGroup, id.Name, "")
				if err != nil {
					return "", fmt.Errorf("obtaining NIC %q : %+v", id.String(), err)
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
					id, err := networkParse.PublicIpAddressID(*config.PublicIPAddress.ID)
					if err != nil {
						return "", err
					}

					pip, err := pipClient.Get(ctx, id.ResourceGroup, id.Name, "")
					if err != nil {
						return "", fmt.Errorf("obtaining Public IP %q : %+v", id.String(), err)
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
