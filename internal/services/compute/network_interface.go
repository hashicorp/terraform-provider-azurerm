// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// nolint: deadcode unused
type connectionInfo struct {
	// primaryPrivateAddress is the Primary Private IP Address for this VM
	primaryPrivateAddress string

	// privateAddresses is a slice of the Private IP Addresses supported by this VM
	privateAddresses []string

	// primaryPublicAddress is the Primary Public IP Address for this VM
	primaryPublicAddress string

	// publicAddresses is a slice of the Public IP Addresses supported by this VM
	publicAddresses []string
}

// retrieveConnectionInformation retrieves all of the Public and Private IP Addresses assigned to a Virtual Machine
// nolint: deadcode unused
func retrieveConnectionInformation(ctx context.Context, nicsClient *networkinterfaces.NetworkInterfacesClient, pipsClient *publicipaddresses.PublicIPAddressesClient, input *virtualmachines.VirtualMachineProperties) connectionInfo {
	if input == nil || input.NetworkProfile == nil || input.NetworkProfile.NetworkInterfaces == nil {
		return connectionInfo{}
	}

	privateIPAddresses := make([]string, 0)
	publicIPAddresses := make([]string, 0)

	for _, v := range *input.NetworkProfile.NetworkInterfaces {
		if v.Id == nil {
			continue
		}

		nic := retrieveIPAddressesForNIC(ctx, nicsClient, pipsClient, *v.Id)
		if nic == nil {
			continue
		}

		privateIPAddresses = append(privateIPAddresses, nic.privateIPAddresses...)
		publicIPAddresses = append(publicIPAddresses, nic.publicIPAddresses...)
	}

	primaryPrivateAddress := ""
	if len(privateIPAddresses) > 0 {
		primaryPrivateAddress = privateIPAddresses[0]
	}
	primaryPublicAddress := ""
	if len(publicIPAddresses) > 0 {
		primaryPublicAddress = publicIPAddresses[0]
	}

	return connectionInfo{
		primaryPrivateAddress: primaryPrivateAddress,
		privateAddresses:      privateIPAddresses,
		primaryPublicAddress:  primaryPublicAddress,
		publicAddresses:       publicIPAddresses,
	}
}

// nolint: deadcode unused
type interfaceDetails struct {
	// privateIPAddresses is a slice of the Private IP Addresses supported by this VM
	privateIPAddresses []string

	// publicIPAddresses is a slice of the Public IP Addresses supported by this VM
	publicIPAddresses []string
}

// retrieveIPAddressesForNIC returns the Public and Private IP Addresses associated
// with the specified Network Interface
// nolint: deadcode unused
func retrieveIPAddressesForNIC(ctx context.Context, nicClient *networkinterfaces.NetworkInterfacesClient, pipClient *publicipaddresses.PublicIPAddressesClient, nicID string) *interfaceDetails {
	id, err := commonids.ParseNetworkInterfaceID(nicID)
	if err != nil {
		return nil
	}

	nic, err := nicClient.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return nil
	}

	if nic.Model == nil || nic.Model.Properties == nil || nic.Model.Properties.IPConfigurations == nil {
		return nil
	}

	privateIPAddresses := make([]string, 0)
	publicIPAddresses := make([]string, 0)
	for _, config := range *nic.Model.Properties.IPConfigurations {
		if props := config.Properties; props != nil {
			if props.PrivateIPAddress != nil {
				privateIPAddresses = append(privateIPAddresses, *props.PrivateIPAddress)
			}

			if pip := props.PublicIPAddress; pip != nil {
				if pip.Id != nil {
					publicIPAddress, err := retrievePublicIPAddress(ctx, pipClient, *pip.Id)
					if err != nil {
						continue
					}

					if publicIPAddress != nil {
						publicIPAddresses = append(publicIPAddresses, *publicIPAddress)
					}
				}
			}
		}
	}

	return &interfaceDetails{
		privateIPAddresses: privateIPAddresses,
		publicIPAddresses:  publicIPAddresses,
	}
}

// retrievePublicIPAddress returns the Public IP Address associated with an Azure Public IP
// nolint: deadcode unused
func retrievePublicIPAddress(ctx context.Context, client *publicipaddresses.PublicIPAddressesClient, publicIPAddressID string) (*string, error) {
	id, err := commonids.ParsePublicIPAddressID(publicIPAddressID)
	if err != nil {
		return nil, err
	}

	pip, err := client.Get(ctx, *id, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		return nil, err
	}

	if model := pip.Model; model != nil {
		if props := model.Properties; props != nil {
			// if it's Static it'll always have an IP Address assigned
			// however there's a bug here where Dynamic IP's can take some time until it's assigned after attachment
			// TODO: fix the bug with Dynamic IP's here
			return props.IPAddress, nil
		}
	}

	return nil, nil
}

// setConnectionInformation sets the connection information required for Provisioners
// to connect to the Virtual Machine. A Public IP Address is used if one is available
// but this falls back to a Private IP Address (which should always exist)
// nolint: deadcode unused
func setConnectionInformation(d *pluginsdk.ResourceData, input connectionInfo, isWindows bool) {
	provisionerType := "ssh"
	if isWindows {
		provisionerType = "winrm"
	}

	ipAddress := input.primaryPublicAddress
	if ipAddress == "" {
		ipAddress = input.primaryPrivateAddress
	}

	d.SetConnInfo(map[string]string{
		"type": provisionerType,
		"host": ipAddress,
	})
}
