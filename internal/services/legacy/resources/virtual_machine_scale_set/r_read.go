// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	opts := virtualmachinescalesets.DefaultGetOperationOptions()
	opts.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	resp, err := client.Get(ctx, *id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found. Removing from State", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.VirtualMachineScaleSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", model.Zones)

		if err := d.Set("sku", flattenAzureRmVirtualMachineScaleSetSku(model.Sku)); err != nil {
			return fmt.Errorf("[DEBUG] setting `sku`: %#v", err)
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("[DEBUG] setting `identity`: %+v", err)
		}

		if plan := model.Plan; plan != nil {
			if err := d.Set("plan", flattenAzureRmVirtualMachineScaleSetPlan(plan)); err != nil {
				return fmt.Errorf("[DEBUG] setting `plan`: %#v", err)
			}
		}

		if props := model.Properties; props != nil {
			if upgradePolicy := props.UpgradePolicy; upgradePolicy != nil {
				d.Set("upgrade_policy_mode", pointer.From(upgradePolicy.Mode))
				if policy := upgradePolicy.AutomaticOSUpgradePolicy; policy != nil {
					d.Set("automatic_os_upgrade", policy.EnableAutomaticOSUpgrade)
				}

				if rollingUpgradePolicy := upgradePolicy.RollingUpgradePolicy; rollingUpgradePolicy != nil {
					if err := d.Set("rolling_upgrade_policy", flattenAzureRmVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicy)); err != nil {
						return fmt.Errorf("[DEBUG] setting Virtual Machine Scale Set Rolling Upgrade Policy error: %#v", err)
					}
				}

				if proximityPlacementGroup := props.ProximityPlacementGroup; proximityPlacementGroup != nil {
					d.Set("proximity_placement_group_id", proximityPlacementGroup.Id)
				}
			}
			d.Set("overprovision", props.Overprovision)
			d.Set("single_placement_group", props.SinglePlacementGroup)

			if profile := props.VirtualMachineProfile; profile != nil {
				d.Set("license_type", profile.LicenseType)
				d.Set("priority", pointer.FromEnum(profile.Priority))
				d.Set("eviction_policy", pointer.FromEnum(profile.EvictionPolicy))

				if err := d.Set("os_profile", flattenAzureRMVirtualMachineScaleSetOsProfile(d, profile.OsProfile)); err != nil {
					return fmt.Errorf("[DEBUG] setting `os_profile`: %#v", err)
				}

				if osProfile := profile.OsProfile; osProfile != nil {
					if linuxConfiguration := osProfile.LinuxConfiguration; linuxConfiguration != nil {
						if err := d.Set("os_profile_linux_config", flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(linuxConfiguration)); err != nil {
							return fmt.Errorf("[DEBUG] setting `os_profile_linux_config`: %#v", err)
						}
					}

					if secrets := osProfile.Secrets; secrets != nil {
						if err := d.Set("os_profile_secrets", flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(secrets)); err != nil {
							return fmt.Errorf("[DEBUG] setting `os_profile_secrets`: %#v", err)
						}
					}

					if windowsConfiguration := osProfile.WindowsConfiguration; windowsConfiguration != nil {
						if err := d.Set("os_profile_windows_config", flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(windowsConfiguration)); err != nil {
							return fmt.Errorf("[DEBUG] setting `os_profile_windows_config`: %#v", err)
						}
					}
				}

				if diagnosticsProfile := profile.DiagnosticsProfile; diagnosticsProfile != nil {
					if bootDiagnostics := diagnosticsProfile.BootDiagnostics; bootDiagnostics != nil {
						if err := d.Set("boot_diagnostics", flattenAzureRmVirtualMachineScaleSetBootDiagnostics(bootDiagnostics)); err != nil {
							return fmt.Errorf("[DEBUG] setting `boot_diagnostics`: %#v", err)
						}
					}
				}

				if networkProfile := profile.NetworkProfile; networkProfile != nil {
					if hp := networkProfile.HealthProbe; hp != nil {
						if id := hp.Id; id != nil {
							d.Set("health_probe_id", id)
						}
					}

					if err := d.Set("network_profile", flattenAzureRmVirtualMachineScaleSetNetworkProfile(networkProfile)); err != nil {
						return fmt.Errorf("[DEBUG] setting `network_profile`: %#v", err)
					}
				}

				if storageProfile := profile.StorageProfile; storageProfile != nil {
					if dataDisks := props.VirtualMachineProfile.StorageProfile.DataDisks; dataDisks != nil {
						if err := d.Set("storage_profile_data_disk", flattenAzureRmVirtualMachineScaleSetStorageProfileDataDisk(dataDisks)); err != nil {
							return fmt.Errorf("[DEBUG] setting `storage_profile_data_disk`: %#v", err)
						}
					}

					if imageRef := storageProfile.ImageReference; imageRef != nil {
						if err := d.Set("storage_profile_image_reference", flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(imageRef)); err != nil {
							return fmt.Errorf("[DEBUG] setting `storage_profile_image_reference`: %#v", err)
						}
					}

					if osDisk := storageProfile.OsDisk; osDisk != nil {
						if err := d.Set("storage_profile_os_disk", flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(osDisk)); err != nil {
							return fmt.Errorf("[DEBUG] setting `storage_profile_os_disk`: %#v", err)
						}
					}
				}

				if extensionProfile := props.VirtualMachineProfile.ExtensionProfile; extensionProfile != nil {
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
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func resourceVirtualMachineScaleSetStorageProfileImageReferenceHash(v any) int {
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
		if v, ok := m["version"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
		if v, ok := m["id"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetStorageProfileOsDiskHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", m["name"].(string))

		if v, ok := m["vhd_containers"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(*pluginsdk.Set).List())
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetNetworkConfigurationHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", m["name"].(string))
		fmt.Fprintf(&buf, "%t-", m["primary"].(bool))

		if v, ok := m["accelerated_networking"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
		if v, ok := m["ip_forwarding"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
		if v, ok := m["network_security_group_id"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
		if v, ok := m["dns_settings"].(map[string]any); ok {
			if k, ok := v["dns_servers"]; ok {
				fmt.Fprintf(&buf, "%s-", k)
			}
		}
		if ipConfig, ok := m["ip_configuration"].([]any); ok {
			for _, it := range ipConfig {
				config := it.(map[string]any)
				if name, ok := config["name"]; ok {
					fmt.Fprintf(&buf, "%s-", name.(string))
				}
				if subnetid, ok := config["subnet_id"]; ok {
					fmt.Fprintf(&buf, "%s-", subnetid.(string))
				}
				if appPoolId, ok := config["application_gateway_backend_address_pool_ids"]; ok {
					fmt.Fprintf(&buf, "%s-", appPoolId.(*pluginsdk.Set).List())
				}
				if appSecGroup, ok := config["application_security_group_ids"]; ok {
					fmt.Fprintf(&buf, "%s-", appSecGroup.(*pluginsdk.Set).List())
				}
				if lbPoolIds, ok := config["load_balancer_backend_address_pool_ids"]; ok {
					fmt.Fprintf(&buf, "%s-", lbPoolIds.(*pluginsdk.Set).List())
				}
				if lbInNatRules, ok := config["load_balancer_inbound_nat_rules_ids"]; ok {
					fmt.Fprintf(&buf, "%s-", lbInNatRules.(*pluginsdk.Set).List())
				}
				if primary, ok := config["primary"]; ok {
					fmt.Fprintf(&buf, "%t-", primary.(bool))
				}
				if publicIPConfig, ok := config["public_ip_address_configuration"].([]any); ok {
					for _, publicIPIt := range publicIPConfig {
						publicip := publicIPIt.(map[string]any)
						if publicIPConfigName, ok := publicip["name"]; ok {
							fmt.Fprintf(&buf, "%s-", publicIPConfigName.(string))
						}
						if idle_timeout, ok := publicip["idle_timeout"]; ok {
							fmt.Fprintf(&buf, "%d-", idle_timeout.(int))
						}
						if dnsLabel, ok := publicip["domain_name_label"]; ok {
							fmt.Fprintf(&buf, "%s-", dnsLabel.(string))
						}
					}
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetOsProfileLinuxConfigHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%t-", m["disable_password_authentication"].(bool))

		if sshKeys, ok := m["ssh_keys"].([]any); ok {
			for _, item := range sshKeys {
				k := item.(map[string]any)
				if path, ok := k["path"]; ok {
					fmt.Fprintf(&buf, "%s-", path.(string))
				}
				if data, ok := k["key_data"]; ok {
					fmt.Fprintf(&buf, "%s-", data.(string))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceVirtualMachineScaleSetOsProfileWindowsConfigHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["provision_vm_agent"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
		if v, ok := m["enable_automatic_upgrades"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}
	}

	return pluginsdk.HashString(buf.String())
}

// When upgrade_policy_mode is not Rolling, we will just ignore rolling_upgrade_policy (returns true).
func azureRmVirtualMachineScaleSetSuppressRollingUpgradePolicyDiff(k, _, new string, d *pluginsdk.ResourceData) bool {
	if k == "rolling_upgrade_policy.#" && new == "0" {
		return strings.ToLower(d.Get("upgrade_policy_mode").(string)) != "rolling"
	}
	return false
}

// Make sure rolling_upgrade_policy is default value when upgrade_policy_mode is not Rolling.
func azureRmVirtualMachineScaleSetCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	mode := d.Get("upgrade_policy_mode").(string)
	if strings.ToLower(mode) != "rolling" {
		if policyRaw, ok := d.GetOk("rolling_upgrade_policy.0"); ok {
			policy := policyRaw.(map[string]any)
			isDefault := (policy["max_batch_instance_percent"].(int) == 20) &&
				(policy["max_unhealthy_instance_percent"].(int) == 20) &&
				(policy["max_unhealthy_upgraded_instance_percent"].(int) == 20) &&
				(policy["pause_time_between_batches"] == "PT0S")
			if !isDefault {
				return fmt.Errorf("ff `upgrade_policy_mode` is `%s`, `rolling_upgrade_policy` must be removed or set to default values", mode)
			}
		}
	}
	return nil
}
