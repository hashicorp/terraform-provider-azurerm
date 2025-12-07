// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func ValidateStrataCloudManagerTenantNameExists(ctx context.Context, client *clients.Client, subscriptionId string, tenantName string) error {
	if strings.TrimSpace(tenantName) == "" {
		return fmt.Errorf("strata cloud manager tenant name cannot be empty")
	}

	subscriptionIdParsed := commonids.NewSubscriptionID(subscriptionId)

	resp, err := client.PaloAlto.PaloAltoClient_v2025_05_23.PaloAltoNetworksCloudngfws.PaloAltoNetworksCloudngfwOperationslistCloudManagerTenants(ctx, subscriptionIdParsed)
	if err != nil {
		return fmt.Errorf("retrieving list of available Strata Cloud Manager tenants: %v", err)
	}

	if resp.Model == nil || resp.Model.Value == nil {
		return fmt.Errorf("unable to retrieve list of available Strata Cloud Manager tenants")
	}

	availableTenants := resp.Model.Value
	for _, tenant := range availableTenants {
		if strings.EqualFold(tenant, tenantName) {
			return nil
		}
	}

	return fmt.Errorf("strata cloud manager tenant %q is not available. Available tenants are: %s", tenantName, strings.Join(availableTenants, ", "))
}
