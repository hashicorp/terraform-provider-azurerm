package vminstancehybrididentitymetadatas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMInstanceHybridIdentityMetadataProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	PublicKey         *string            `json:"publicKey,omitempty"`
	ResourceUid       *string            `json:"resourceUid,omitempty"`
}
