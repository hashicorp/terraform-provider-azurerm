package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureSku struct {
	Capacity *int64       `json:"capacity,omitempty"`
	Name     AzureSkuName `json:"name"`
	Tier     AzureSkuTier `json:"tier"`
}
