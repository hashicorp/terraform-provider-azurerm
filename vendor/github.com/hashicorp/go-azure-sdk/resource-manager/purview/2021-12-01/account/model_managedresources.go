package account

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedResources struct {
	EventHubNamespace *string `json:"eventHubNamespace,omitempty"`
	ResourceGroup     *string `json:"resourceGroup,omitempty"`
	StorageAccount    *string `json:"storageAccount,omitempty"`
}
