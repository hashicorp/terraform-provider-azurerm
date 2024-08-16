package site

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SitePropertiesFormat struct {
	NetworkFunctions  *[]SubResource     `json:"networkFunctions,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
