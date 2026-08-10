package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultSigningKeyParameters struct {
	ResourceGroupName string                           `json:"resourceGroupName"`
	SecretName        string                           `json:"secretName"`
	SecretVersion     string                           `json:"secretVersion"`
	SubscriptionId    string                           `json:"subscriptionId"`
	TypeName          KeyVaultSigningKeyParametersType `json:"typeName"`
	VaultName         string                           `json:"vaultName"`
}
