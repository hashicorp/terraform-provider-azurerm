// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_cluster

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
)

func expandEventHubClusterSkuName(skuName string) eventhubsclusters.ClusterSku {
	// "sku_name" is validated to be in this format above, and is required
	skuParts := strings.Split(skuName, "_")
	name := skuParts[0]
	capacity, _ := strconv.Atoi(skuParts[1])
	return eventhubsclusters.ClusterSku{
		Name:     eventhubsclusters.ClusterSkuName(name),
		Capacity: pointer.To(int64(capacity)),
	}
}
