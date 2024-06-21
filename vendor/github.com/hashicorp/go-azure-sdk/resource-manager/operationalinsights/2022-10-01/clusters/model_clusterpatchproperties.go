package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPatchProperties struct {
	BillingType        *BillingType        `json:"billingType,omitempty"`
	KeyVaultProperties *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
}
