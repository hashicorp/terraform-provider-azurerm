// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
)

// DataboxEdgeDeviceSku type
type DataboxEdgeDeviceSku struct {
	Name string
	Tier string
}

// DataboxEdgeDeviceSkuName parses the input string into a DataboxEdgeDeviceSku type
func DataboxEdgeDeviceSkuName(input string) (*DataboxEdgeDeviceSku, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return nil, fmt.Errorf("unable to parse Databox Edge Device 'sku_name' %q", input)
	}

	skuParts := strings.Split(input, "-")

	if strings.TrimSpace(skuParts[0]) == "" {
		return nil, fmt.Errorf("invalid Databox Edge Device 'sku_name' %q", input)
	}

	// There is only one possible Tier so always set value to Standard
	databoxEdgeDeviceSku := DataboxEdgeDeviceSku{
		Name: skuParts[0],
		Tier: string(devices.SkuTierStandard),
	}

	return &databoxEdgeDeviceSku, nil
}
