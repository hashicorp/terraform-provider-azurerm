// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type networkInterfaceIPConfigurationLockingDetails struct {
	subnetIdsToLock         []string
	virtualNetworkIdsToLock []string
}

func (details networkInterfaceIPConfigurationLockingDetails) lock() {
	locks.MultipleByID(&details.virtualNetworkIdsToLock)
	locks.MultipleByID(&details.subnetIdsToLock)
}

func (details networkInterfaceIPConfigurationLockingDetails) unlock() {
	locks.UnlockMultipleByID(&details.subnetIdsToLock)
	locks.UnlockMultipleByID(&details.virtualNetworkIdsToLock)
}

func determineResourcesToLockFromIPConfiguration(input *[]networkinterfaces.NetworkInterfaceIPConfiguration) (*networkInterfaceIPConfigurationLockingDetails, error) {
	if input == nil {
		return &networkInterfaceIPConfigurationLockingDetails{
			subnetIdsToLock:         []string{},
			virtualNetworkIdsToLock: []string{},
		}, nil
	}

	subnetIdsToLock := make([]string, 0)
	virtualNetworkIdsToLock := make([]string, 0)

	for _, config := range *input {
		if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
			continue
		}

		id, err := commonids.ParseSubnetID(*config.Properties.Subnet.Id)
		if err != nil {
			return nil, err
		}

		vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)

		if !utils.SliceContainsValue(virtualNetworkIdsToLock, vnetId.ID()) {
			virtualNetworkIdsToLock = append(virtualNetworkIdsToLock, vnetId.ID())
		}

		if !utils.SliceContainsValue(subnetIdsToLock, id.ID()) {
			subnetIdsToLock = append(subnetIdsToLock, id.ID())
		}
	}

	return &networkInterfaceIPConfigurationLockingDetails{
		subnetIdsToLock:         subnetIdsToLock,
		virtualNetworkIdsToLock: virtualNetworkIdsToLock,
	}, nil
}
