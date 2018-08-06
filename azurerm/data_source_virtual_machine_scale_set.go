package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualMachineScaleSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualMachineScaleSetRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"vm_hostnames": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"vm_statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"vm_primary_private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceArmVirtualMachineScaleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	vmsClient := meta.(*ArmClient).vmScaleSetVMsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	vmssResp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error reading VM Scale Set: %+v", err)
	}

	d.SetId(*vmssResp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := vmssResp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	vmsResp, err := vmsClient.ListComplete(ctx, resourceGroupName, name, "", "", "")
	if err != nil {
		if utils.ResponseWasNotFound(vmsResp.Response().Response) {
			log.Printf("[INFO] VMs for AzureRM Virtual Machine Scale Set (%s) were Not Found.", name)
		}
		return fmt.Errorf("Error making Read request for VMs on Azure Virtual Machine Scale Set %s: %+v", name, err)
	}

	vm_hostnames, vm_statuses, vm_primary_private_ip_addresses := flattenAzureRmVirtualMachineScaleSetVMs(vmsResp, meta, resourceGroupName, name)
	if err := d.Set("vm_hostnames", vm_hostnames); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set hostname array error: %#v", err)
	}

	if err := d.Set("vm_statuses", vm_statuses); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set vm hostnames array error: %#v", err)
	}

	if err := d.Set("vm_primary_private_ip_addresses", vm_primary_private_ip_addresses); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set vm primary private ip array error: %#v", err)
	}

	return nil
}

func flattenAzureRmVirtualMachineScaleSetVMs(vmsItt compute.VirtualMachineScaleSetVMListResultIterator, meta interface{}, resGroup string, name string) ([]string, []string, []string) {
	vm_statuses := make([]string, 0, len(*vmsItt.Response().Value))
	vm_hostnames := make([]string, 0, len(*vmsItt.Response().Value))
	vm_primary_private_ip_addresses := make([]string, 0, len(*vmsItt.Response().Value))
	for vmsItt.NotDone() {
		vmssVM := vmsItt.Value()
		s := map[string]interface{}{
			"name":     *vmssVM.Name,
			"type":     *vmssVM.Type,
			"hostname": *vmssVM.OsProfile.ComputerName,
			"status":   *vmssVM.ProvisioningState,
		}
		primary_private_ip_address := flattenAzureRmVirtualMachineScaleSetVMNetworkInterfaces(meta, resGroup, name, *vmssVM.InstanceID)
		vm_statuses = append(vm_statuses, s["status"].(string))
		vm_hostnames = append(vm_hostnames, s["hostname"].(string))
		vm_primary_private_ip_addresses = append(vm_primary_private_ip_addresses, primary_private_ip_address)
		vmsItt.Next()
	}
	return vm_hostnames, vm_statuses, vm_primary_private_ip_addresses
}

func flattenAzureRmVirtualMachineScaleSetVMNetworkInterfaces(meta interface{}, resGroup string, name string, vmid string) string {
	ifaClient := meta.(*ArmClient).ifaceClient
	ctx := meta.(*ArmClient).StopContext
	ifaResp, err := ifaClient.ListVirtualMachineScaleSetVMNetworkInterfaces(ctx, resGroup, name, vmid)
	if err != nil {
		if utils.ResponseWasNotFound(ifaResp.Response().Response) {
			log.Printf("[INFO] Network Interfaces for VM(%s) in AzureRM Virtual Machine Scale Set (%s) were Not Found.", vmid, name)
		}
		log.Printf("[ERROR] making Read request for Network Interfaces for VM(%s) in Azure Virtual Machine Scale Set %s: %+v", vmid, name, err)
	}
	primary_private_ip := flattenAzureRmVirtualMachineScaleSetVMNetworkInterfaceIpConfs(ifaResp)
	return primary_private_ip
}

func flattenAzureRmVirtualMachineScaleSetVMNetworkInterfaceIpConfs(nicsItt network.InterfaceListResultPage) string {
	primary_nic_private_ip_address := ""
	for nicsItt.NotDone() {
		nics := nicsItt.Values()
		for _, nic := range nics {
			s := map[string]interface{}{
				"id":      *nic.ID,
				"primary": *nic.Primary,
			}
			for _, ipconf := range *nic.IPConfigurations {
				if ipconf.Primary != nil && *ipconf.Primary {
					s["private_ip_address"] = *ipconf.PrivateIPAddress
					if nic.Primary != nil && *nic.Primary {
						primary_nic_private_ip_address = *ipconf.PrivateIPAddress
					}
				}
			}
		}
		nicsItt.Next()
	}
	return primary_nic_private_ip_address
}
