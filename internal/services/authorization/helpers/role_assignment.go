// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
)

func GetTenantIdBySubscriptionId(ctx context.Context, client *subscriptions.SubscriptionsClient, subscriptionId string) (string, error) {
	id := commonids.NewSubscriptionID(subscriptionId)
	resp, err := client.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("retrieving %s: %+v", id, err)
	}
	tenantId := ""
	if model := resp.Model; model != nil && model.TenantId != nil {
		tenantId = *model.TenantId
	}

	if tenantId == "" {
		return "", fmt.Errorf("retrieving %s: `tenantId` was nil", id)
	}
	return tenantId, nil
}
