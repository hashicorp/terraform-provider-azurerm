package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MapsAccountProperties struct {
	Cors              *CorsRules        `json:"cors,omitempty"`
	DisableLocalAuth  *bool             `json:"disableLocalAuth,omitempty"`
	Encryption        *Encryption       `json:"encryption,omitempty"`
	LinkedResources   *[]LinkedResource `json:"linkedResources,omitempty"`
	ProvisioningState *string           `json:"provisioningState,omitempty"`
	UniqueId          *string           `json:"uniqueId,omitempty"`
}
