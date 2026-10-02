package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSProfileProvisioningData struct {
	AdminPassword *string `json:"adminPassword,omitempty"`
	CustomData    *string `json:"customData,omitempty"`
}
