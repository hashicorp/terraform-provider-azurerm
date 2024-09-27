package hybrididentitymetadata

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridIdentityMetadataProperties struct {
	Identity          *identity.SystemAssigned `json:"identity,omitempty"`
	ProvisioningState *string                  `json:"provisioningState,omitempty"`
	PublicKey         *string                  `json:"publicKey,omitempty"`
	ResourceUid       *string                  `json:"resourceUid,omitempty"`
}
