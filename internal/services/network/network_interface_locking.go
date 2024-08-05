// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type networkInterfaceIPConfigurationLockingDetails struct {
	subnetNamesToLock         []string
	virtualNetworkNamesToLock []string
}

func (details networkInterfaceIPConfigurationLockingDetails) lock() {
	locks.MultipleByName(&details.virtualNetworkNamesToLock, VirtualNetworkResourceName)
	locks.MultipleByName(&details.subnetNamesToLock, SubnetResourceName)
}

func (details networkInterfaceIPConfigurationLockingDetails) unlock() {
	locks.UnlockMultipleByName(&details.subnetNamesToLock, SubnetResourceName)
	locks.UnlockMultipleByName(&details.virtualNetworkNamesToLock, VirtualNetworkResourceName)
}

func determineResourcesToLockFromIPConfiguration(input *[]networkinterfaces.NetworkInterfaceIPConfiguration) (*networkInterfaceIPConfigurationLockingDetails, error) {
	if input == nil {
		return &networkInterfaceIPConfigurationLockingDetails{
			subnetNamesToLock:         []string{},
			virtualNetworkNamesToLock: []string{},
		}, nil
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, config := range *input {
		if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
			continue
		}

		id, err := commonids.ParseSubnetID(*config.Properties.Subnet.Id)
		if err != nil {
			return nil, err
		}

		virtualNetworkName := id.VirtualNetworkName
		subnetName := id.SubnetName

		if !utils.SliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
			virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
		}

		if !utils.SliceContainsValue(subnetNamesToLock, subnetName) {
			subnetNamesToLock = append(subnetNamesToLock, subnetName)
		}
	}

	return &networkInterfaceIPConfigurationLockingDetails{
		subnetNamesToLock:         subnetNamesToLock,
		virtualNetworkNamesToLock: virtualNetworkNamesToLock,
	}, nil
}
