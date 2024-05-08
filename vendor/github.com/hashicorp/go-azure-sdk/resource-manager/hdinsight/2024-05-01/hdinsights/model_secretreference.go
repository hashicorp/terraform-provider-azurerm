package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretReference struct {
	KeyVaultObjectName string             `json:"keyVaultObjectName"`
	ReferenceName      string             `json:"referenceName"`
	Type               KeyVaultObjectType `json:"type"`
	Version            *string            `json:"version,omitempty"`
}
