package managedidentities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResource struct {
	Id                      *string `json:"id,omitempty"`
	Name                    *string `json:"name,omitempty"`
	ResourceGroup           *string `json:"resourceGroup,omitempty"`
	SubscriptionDisplayName *string `json:"subscriptionDisplayName,omitempty"`
	SubscriptionId          *string `json:"subscriptionId,omitempty"`
	Type                    *string `json:"type,omitempty"`
}
