package compute_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OrchestratedVirtualMachineScaleSetResource struct {
}

func (t OrchestratedVirtualMachineScaleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachineScaleSets"]

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := clients.Compute.VMScaleSetClient.Get(ctx, resGroup, name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		return nil, fmt.Errorf("retrieving Virtual Machine Scale Set %q", id)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (OrchestratedVirtualMachineScaleSetResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VirtualMachineScaleSetID(state.ID)
	if err != nil {
		return nil, err
	}

	// this is a preview feature we don't want to use right now
	var forceDelete *bool = nil
	future, err := client.Compute.VMScaleSetClient.Delete(ctx, id.ResourceGroup, id.Name, forceDelete)
	if err != nil {
		return nil, fmt.Errorf("Bad: deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Compute.VMScaleSetClient.Client); err != nil {
		return nil, fmt.Errorf("Bad: waiting for deletion of %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (OrchestratedVirtualMachineScaleSetResource) hasApplicationGateway(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.VirtualMachineScaleSetID(state.ID)
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	read, err := client.Compute.VMScaleSetClient.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		return err
	}

	if props := read.VirtualMachineScaleSetProperties; props != nil {
		if vmProfile := props.VirtualMachineProfile; vmProfile != nil {
			if nwProfile := vmProfile.NetworkProfile; nwProfile != nil {
				if nics := nwProfile.NetworkInterfaceConfigurations; nics != nil {
					for _, nic := range *nics {
						if nic.IPConfigurations == nil {
							continue
						}

						for _, config := range *nic.IPConfigurations {
							if config.ApplicationGatewayBackendAddressPools == nil {
								continue
							}

							if len(*config.ApplicationGatewayBackendAddressPools) > 0 {
								return nil
							}
						}
					}
				}
			}
		}
	}

	return fmt.Errorf("application gateway configuration was missing")
}

func (OrchestratedVirtualMachineScaleSetResource) hasLoadBalancer(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.VirtualMachineScaleSetID(state.ID)
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	read, err := client.Compute.VMScaleSetClient.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		return err
	}

	if props := read.VirtualMachineScaleSetProperties; props != nil {
		if vmProfile := props.VirtualMachineProfile; vmProfile != nil {
			if nwProfile := vmProfile.NetworkProfile; nwProfile != nil {
				if nics := nwProfile.NetworkInterfaceConfigurations; nics != nil {
					for _, nic := range *nics {
						if nic.IPConfigurations == nil {
							continue
						}

						for _, config := range *nic.IPConfigurations {
							if config.LoadBalancerBackendAddressPools == nil {
								continue
							}

							if len(*config.LoadBalancerBackendAddressPools) > 0 {
								return nil
							}
						}
					}
				}
			}
		}
	}

	return fmt.Errorf("load balancer configuration was missing")
}
