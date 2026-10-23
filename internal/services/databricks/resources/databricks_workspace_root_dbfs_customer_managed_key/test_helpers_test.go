// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace_root_dbfs_customer_managed_key_test

import (
	"os"
)

type DatabricksWorkspaceAlternateSubscription struct {
	tenantID       string
	subscriptionID string
}

func altSubscriptionCheck() *DatabricksWorkspaceAlternateSubscription {
	altSubscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID_ALT")
	altTenantID := os.Getenv("ARM_TENANT_ID")

	if altSubscriptionID == "" || altTenantID == "" {
		return nil
	}

	return &DatabricksWorkspaceAlternateSubscription{
		tenantID:       altTenantID,
		subscriptionID: altSubscriptionID,
	}
}
