// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
)

func FindNetworkInterfaceIPConfiguration(input *[]networkinterfaces.NetworkInterfaceIPConfiguration, name string) *networkinterfaces.NetworkInterfaceIPConfiguration {
	if input == nil {
		return nil
	}

	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		if *v.Name == name {
			return &v
		}
	}

	return nil
}

func updateNetworkInterfaceIPConfiguration(config networkinterfaces.NetworkInterfaceIPConfiguration, configs *[]networkinterfaces.NetworkInterfaceIPConfiguration) *[]networkinterfaces.NetworkInterfaceIPConfiguration {
	output := make([]networkinterfaces.NetworkInterfaceIPConfiguration, 0)
	if configs == nil {
		return &output
	}

	for _, v := range *configs {
		if v.Name == nil {
			continue
		}

		if *v.Name != *config.Name {
			output = append(output, v)
		} else {
			output = append(output, config)
		}
	}

	return &output
}
