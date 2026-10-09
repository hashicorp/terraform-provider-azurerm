package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultReferenceWithStatus struct {
	ReferenceStatus *string `json:"referenceStatus,omitempty"`
	SecretUri       *string `json:"secretUri,omitempty"`
}
