package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementSkuZoneDetails struct {
	Capabilities *[]ApiManagementSkuCapabilities `json:"capabilities,omitempty"`
	Name         *[]string                       `json:"name,omitempty"`
}
