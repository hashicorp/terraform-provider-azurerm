// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	intStor "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func resourceVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
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
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func resourceVirtualMachineDeleteVhd(ctx context.Context, storageClient *intStor.Client, subscriptionId string, vhd *virtualmachines.VirtualHardDisk) error {
	if vhd == nil {
		return errors.New("`vhd` was nil`")
	}
	if vhd.Uri == nil {
		return errors.New("`vhd.Uri` was nil`")
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

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %s", err)
	}

	input := blobs.DeleteInput{}
	if _, err := blobsClient.Delete(ctx, id.ContainerName, id.BlobName, input); err != nil {
		return fmt.Errorf("deleting Blob %q (Container %q in %s): %+v", id.BlobName, id.ContainerName, id.AccountId, err)
	}

	return nil
}

func resourceVirtualMachineDeleteManagedDisk(d *pluginsdk.ResourceData, disk *virtualmachines.ManagedDiskParameters, meta any) error {
	if disk == nil {
		return errors.New("`disk` was nil`")
	}
	if disk.Id == nil {
		return errors.New("`disk.Id` was nil`")
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

func resourceVirtualMachineStorageOsProfileHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", m["admin_username"].(string))
		fmt.Fprintf(&buf, "%s-", m["computer_name"].(string))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineStorageOsProfileWindowsConfigHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["provision_vm_agent"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
		if v, ok := m["enable_automatic_upgrades"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
		if v, ok := m["timezone"]; ok {
			fmt.Fprintf(&buf, "%s-", strings.ToLower(v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineStorageOsProfileLinuxConfigHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["disable_password_authentication"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineStorageImageReferenceHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["publisher"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
		if v, ok := m["offer"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
		if v, ok := m["sku"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
		if v, ok := m["id"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineGetManagedDiskInfo(d *pluginsdk.ResourceData, disk *virtualmachines.ManagedDiskParameters, meta any) (*disks.Disk, error) {
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

func determineVirtualMachineIPAddress(ctx context.Context, meta any, props *virtualmachines.VirtualMachineProperties) (string, error) {
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
		return "", errors.New("a Network Interface wasn't found on the Virtual Machine")
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

	return "", errors.New("no Public or Private IP Address found on the Primary Network Interface")
}
