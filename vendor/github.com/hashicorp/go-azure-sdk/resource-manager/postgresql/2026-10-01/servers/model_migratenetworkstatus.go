package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateNetworkStatus struct {
	ResourceGroupName *string                `json:"resourceGroupName,omitempty"`
	ServerName        *string                `json:"serverName,omitempty"`
	State             *NetworkMigrationState `json:"state,omitempty"`
	SubscriptionId    *string                `json:"subscriptionId,omitempty"`
}
