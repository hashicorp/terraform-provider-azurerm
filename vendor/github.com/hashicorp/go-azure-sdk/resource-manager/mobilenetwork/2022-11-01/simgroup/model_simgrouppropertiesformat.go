package simgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimGroupPropertiesFormat struct {
	EncryptionKey     *KeyVaultKey             `json:"encryptionKey,omitempty"`
	MobileNetwork     *MobileNetworkResourceId `json:"mobileNetwork,omitempty"`
	ProvisioningState *ProvisioningState       `json:"provisioningState,omitempty"`
}
