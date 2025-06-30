// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
)

func ManagedRedisCacheSkuName(input string) (*redisenterprise.Sku, error) {
	skuParts := strings.Split(input, "-")
	skuName := skuParts[0]
	var capacity *int64

	if len(skuParts) == 2 {
		capacityValue, err := strconv.ParseInt(skuParts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Managed Redis Cluster 'sku_name', 'capacity' segment must be of type int64, got %q", skuParts[1])
		}
		capacity = &capacityValue
	}

	return &redisenterprise.Sku{
		Name:     redisenterprise.SkuName(skuName),
		Capacity: capacity,
	}, nil
}
