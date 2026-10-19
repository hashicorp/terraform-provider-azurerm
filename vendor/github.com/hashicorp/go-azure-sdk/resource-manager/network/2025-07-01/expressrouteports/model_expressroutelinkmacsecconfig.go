package expressrouteports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkMacSecConfig struct {
	CakSecretIdentifier *string                         `json:"cakSecretIdentifier,omitempty"`
	Cipher              *ExpressRouteLinkMacSecCipher   `json:"cipher,omitempty"`
	CknSecretIdentifier *string                         `json:"cknSecretIdentifier,omitempty"`
	SciState            *ExpressRouteLinkMacSecSciState `json:"sciState,omitempty"`
}
