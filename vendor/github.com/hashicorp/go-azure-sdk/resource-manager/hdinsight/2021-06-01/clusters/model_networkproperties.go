package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProperties struct {
	PrivateLink                *PrivateLink                `json:"privateLink,omitempty"`
	ResourceProviderConnection *ResourceProviderConnection `json:"resourceProviderConnection,omitempty"`
}
