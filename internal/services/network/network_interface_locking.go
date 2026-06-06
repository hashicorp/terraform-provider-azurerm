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
	subnetIDsToLock         []string
	virtualNetworkIDsToLock []string
}

func (details networkInterfaceIPConfigurationLockingDetails) lock() {
	locks.MultipleByID(&details.virtualNetworkIDsToLock)
	locks.MultipleByID(&details.subnetIDsToLock)
}

func (details networkInterfaceIPConfigurationLockingDetails) unlock() {
	locks.UnlockMultipleByID(&details.subnetIDsToLock)
	locks.UnlockMultipleByID(&details.virtualNetworkIDsToLock)
}

func determineResourcesToLockFromIPConfiguration(input *[]networkinterfaces.NetworkInterfaceIPConfiguration) (*networkInterfaceIPConfigurationLockingDetails, error) {
	if input == nil {
		return &networkInterfaceIPConfigurationLockingDetails{
			subnetIDsToLock:         []string{},
			virtualNetworkIDsToLock: []string{},
		}, nil
	}

	subnetIDsToLock := make([]string, 0)
	virtualNetworkIDsToLock := make([]string, 0)

	for _, config := range *input {
		if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
			continue
		}

		id, err := commonids.ParseSubnetID(*config.Properties.Subnet.Id)
		if err != nil {
			return nil, err
		}

		virtualNetworkID := commonids.NewVirtualNetworkID(id.SubscriptionId, id.SubscriptionId, id.VirtualNetworkName)
		if !utils.SliceContainsValue(virtualNetworkIDsToLock, virtualNetworkID.ID()) {
			virtualNetworkIDsToLock = append(virtualNetworkIDsToLock, virtualNetworkID.ID())
		}

		if !utils.SliceContainsValue(subnetIDsToLock, id.ID()) {
			subnetIDsToLock = append(subnetIDsToLock, id.ID())
		}
	}

	return &networkInterfaceIPConfigurationLockingDetails{
		subnetIDsToLock:         subnetIDsToLock,
		virtualNetworkIDsToLock: virtualNetworkIDsToLock,
	}, nil
}
