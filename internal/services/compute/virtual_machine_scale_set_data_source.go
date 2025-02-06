// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/vmsspublicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualMachineScaleSetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"network_interface": VirtualMachineScaleSetNetworkInterfaceSchemaForDataSource(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"instances": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"computer_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"instance_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"latest_model_applied": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"public_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"virtual_machine_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"zone": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"power_state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	instancesClient := meta.(*clients.Client).Compute.VirtualMachineScaleSetVMsClient
	virtualMachinesClient := meta.(*clients.Client).Compute.VirtualMachinesClient
	networkInterfacesClient := meta.(*clients.Client).Network.NetworkInterfacesClient
	publicIPAddressesClient := meta.(*clients.Client).Network.PublicIPAddresses
	vmssPublicIpAddressesClient := meta.(*clients.Client).Network.VMSSPublicIPAddressesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if profile := props.VirtualMachineProfile; profile != nil {
				if nwProfile := profile.NetworkProfile; nwProfile != nil {
					flattenedNics := FlattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
					if err := d.Set("network_interface", flattenedNics); err != nil {
						return fmt.Errorf("setting `network_interface`: %+v", err)
					}
				}
			}
		}
	}

	instances := make([]interface{}, 0)
	virtualMachineScaleSetId := virtualmachinescalesetvms.NewVirtualMachineScaleSetID(subscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName)

	// If the VMSS is in Uniform Orchestration Mode, we can use instanceView for the VMSS instances
	// Flexible VMSS instances cannot use instanceView from the VMSS API
	// Instead we need to use the VM API for instanceView
	optionsVMSS := virtualmachinescalesetvms.DefaultListOperationOptions()
	optionsVM := virtualmachines.DefaultGetOperationOptions()
	var orchestrationMode string
	if props := resp.Model.Properties; props != nil {
		if *props.OrchestrationMode == virtualmachinescalesets.OrchestrationModeUniform {
			expandStr := "instanceView"
			optionsVMSS.Expand = &expandStr
			orchestrationMode = "Uniform"
		}
		if *props.OrchestrationMode == virtualmachinescalesets.OrchestrationModeFlexible {
			optionsVM.Expand = pointer.To(virtualmachines.InstanceViewTypesInstanceView)
			orchestrationMode = "Flexible"
		}
	}

	result, err := instancesClient.ListComplete(ctx, virtualMachineScaleSetId, optionsVMSS)
	if err != nil {
		return fmt.Errorf("listing VM Instances for %q: %+v", id, err)
	}

	var connInfo *connectionInfo
	var vmModel *virtualmachines.VirtualMachine
	for _, item := range result.Items {
		if item.InstanceId != nil {
			vmId := networkinterfaces.NewVirtualMachineID(subscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName, *item.InstanceId)
			nics, err := networkInterfacesClient.ListVirtualMachineScaleSetVMNetworkInterfacesComplete(ctx, vmId)
			if err != nil {
				if !response.WasNotFound(nics.LatestHttpResponse) {
					return fmt.Errorf("listing Network Interfaces for VM Instance %q for %q: %+v", *item.InstanceId, id, err)
				}

				// Network Interfaces of VM in Flexible VMSS are accessed from single VM
				virtualMachineId := virtualmachines.NewVirtualMachineID(subscriptionId, id.ResourceGroupName, *item.InstanceId)
				vm, err := virtualMachinesClient.Get(ctx, virtualMachineId, optionsVM)
				if err != nil {
					return fmt.Errorf("retrieving VM Instance %q for %q: %+v", *item.InstanceId, id, err)
				}
				connInfoRaw := retrieveConnectionInformation(ctx, networkInterfacesClient, publicIPAddressesClient, vm.Model.Properties)
				connInfo = &connInfoRaw
				vmModel = vm.Model
			} else {
				connInfo, err = getVirtualMachineScaleSetVMConnectionInfo(ctx, nics.Items, id.ResourceGroupName, id.VirtualMachineScaleSetName, *item.InstanceId, vmssPublicIpAddressesClient)
				if err != nil {
					return err
				}
			}

			flattenedInstances := flattenVirtualMachineScaleSetVM(item, connInfo, vmModel, orchestrationMode)
			instances = append(instances, flattenedInstances)
		}
	}
	if err := d.Set("instances", instances); err != nil {
		return fmt.Errorf("setting `instances`: %+v", err)
	}

	return nil
}

func getVirtualMachineScaleSetVMConnectionInfo(ctx context.Context, networkInterfaces []networkinterfaces.NetworkInterface, resourceGroupName string, virtualMachineScaleSetName string, virtualmachineIndex string, publicIPAddressesClient *vmsspublicipaddresses.VMSSPublicIPAddressesClient) (*connectionInfo, error) {
	if len(networkInterfaces) == 0 {
		return nil, nil
	}

	primaryPublicAddress := ""
	primaryPrivateAddress := ""
	publicIPAddresses := make([]string, 0)
	privateIPAddresses := make([]string, 0)

	for _, nic := range networkInterfaces {
		if props := nic.Properties; props != nil {
			if ipConfigs := props.IPConfigurations; ipConfigs != nil {
				for _, config := range *ipConfigs {
					if configProps := config.Properties; configProps != nil {
						if pip := configProps.PublicIPAddress; pip != nil {
							pipID, err := commonids.ParseVirtualMachineScaleSetPublicIPAddressIDInsensitively(*pip.Id)
							if err != nil {
								return nil, err
							}

							publicIPAddress, err := publicIPAddressesClient.PublicIPAddressesGetVirtualMachineScaleSetPublicIPAddress(ctx, *pipID, vmsspublicipaddresses.DefaultPublicIPAddressesGetVirtualMachineScaleSetPublicIPAddressOperationOptions())
							if err != nil {
								return nil, fmt.Errorf("reading Public IP Address for VM Instance %q for Virtual Machine Scale Set %q (Resource Group %q): %+v", virtualmachineIndex, virtualMachineScaleSetName, resourceGroupName, err)
							}

							if publicIPAddress.Model == nil || publicIPAddress.Model.Properties == nil {
								return nil, fmt.Errorf("retrieving %s: `model` or `properties` was nil", pipID)
							}

							if pointer.From(props.Primary) && pointer.From(configProps.Primary) {
								primaryPublicAddress = pointer.From(publicIPAddress.Model.Properties.IPAddress)
							}
							publicIPAddresses = append(publicIPAddresses, pointer.From(publicIPAddress.Model.Properties.IPAddress))
						}

						if configProps.PrivateIPAddress != nil {
							if pointer.From(props.Primary) && pointer.From(configProps.Primary) {
								primaryPrivateAddress = pointer.From(configProps.PrivateIPAddress)
							}
							privateIPAddresses = append(privateIPAddresses, pointer.From(configProps.PrivateIPAddress))
						}
					}
				}
			}
		}
	}

	if primaryPublicAddress == "" && len(publicIPAddresses) > 0 {
		primaryPublicAddress = publicIPAddresses[0]
	}

	if primaryPrivateAddress == "" && len(privateIPAddresses) > 0 {
		primaryPrivateAddress = privateIPAddresses[0]
	}

	return &connectionInfo{
		primaryPublicAddress:  primaryPublicAddress,
		publicAddresses:       publicIPAddresses,
		primaryPrivateAddress: primaryPrivateAddress,
		privateAddresses:      privateIPAddresses,
	}, nil
}

func flattenVirtualMachineScaleSetVM(input virtualmachinescalesetvms.VirtualMachineScaleSetVM, connectionInfo *connectionInfo, vm *virtualmachines.VirtualMachine, mode string) map[string]interface{} {
	output := make(map[string]interface{})
	output["name"] = *input.Name
	output["instance_id"] = *input.InstanceId

	if mode == "Flexible" && vm != nil {
		if props := vm.Properties; props != nil {
			if props.VMId != nil {
				output["virtual_machine_id"] = *props.VMId
			}

			if profile := props.OsProfile; profile != nil && profile.ComputerName != nil {
				output["computer_name"] = *profile.ComputerName
			}

			if instance := props.InstanceView; instance != nil {
				if statuses := instance.Statuses; statuses != nil {
					for _, status := range *statuses {
						if status.Code != nil && strings.HasPrefix(strings.ToLower(*status.Code), "powerstate/") {
							output["power_state"] = strings.SplitN(*status.Code, "/", 2)[1]
						}
					}
				}
			}
		}
	}

	if mode == "Uniform" {
		if props := input.Properties; props != nil {
			if props.LatestModelApplied != nil {
				output["latest_model_applied"] = *props.LatestModelApplied
			}

			if props.VMId != nil {
				output["virtual_machine_id"] = *props.VMId
			}

			if profile := props.OsProfile; profile != nil && profile.ComputerName != nil {
				output["computer_name"] = *profile.ComputerName
			}

			if instance := props.InstanceView; instance != nil {
				if statuses := instance.Statuses; statuses != nil {
					for _, status := range *statuses {
						if status.Code != nil && strings.HasPrefix(strings.ToLower(*status.Code), "powerstate/") {
							output["power_state"] = strings.SplitN(*status.Code, "/", 2)[1]
						}
					}
				}
			}
		}
	}

	zone := ""
	if input.Zones != nil {
		if zones := *input.Zones; len(zones) > 0 {
			zone = zones[0]
		}
	}
	output["zone"] = zone

	if connectionInfo != nil {
		output["private_ip_address"] = connectionInfo.primaryPrivateAddress
		output["private_ip_addresses"] = connectionInfo.privateAddresses
		output["public_ip_address"] = connectionInfo.primaryPublicAddress
		output["public_ip_addresses"] = connectionInfo.publicAddresses
	}

	return output
}
