package network

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type networkInterfaceIPConfigurationLockingDetails struct {
	subnetNamesToLock         []string
	virtualNetworkNamesToLock []string
}

func (details networkInterfaceIPConfigurationLockingDetails) lock() {
	locks.MultipleByName(&details.subnetNamesToLock, SubnetResourceName)
	locks.MultipleByName(&details.virtualNetworkNamesToLock, VirtualNetworkResourceName)
}

func (details networkInterfaceIPConfigurationLockingDetails) unlock() {
	locks.UnlockMultipleByName(&details.subnetNamesToLock, SubnetResourceName)
	locks.UnlockMultipleByName(&details.virtualNetworkNamesToLock, VirtualNetworkResourceName)
}

func determineResourcesToLockFromIPConfiguration(input *[]network.InterfaceIPConfiguration) (*networkInterfaceIPConfigurationLockingDetails, error) {
	if input == nil {
		return &networkInterfaceIPConfigurationLockingDetails{
			subnetNamesToLock:         []string{},
			virtualNetworkNamesToLock: []string{},
		}, nil
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, config := range *input {
		if config.Subnet == nil || config.Subnet.ID == nil {
			continue
		}

		id, err := parse.SubnetID(*config.Subnet.ID)
		if err != nil {
			return nil, err
		}

		virtualNetworkName := id.VirtualNetworkName
		subnetName := id.Name

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
