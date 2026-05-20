package datadogmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResubscribeProperties struct {
	AzureSubscriptionId *string      `json:"azureSubscriptionId,omitempty"`
	ResourceGroup       *string      `json:"resourceGroup,omitempty"`
	Sku                 *ResourceSku `json:"sku,omitempty"`
}
