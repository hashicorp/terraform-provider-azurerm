package webapps

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiKVReferenceProperties struct {
	ActiveVersion *string                            `json:"activeVersion,omitempty"`
	Details       *string                            `json:"details,omitempty"`
	IdentityType  *identity.SystemAndUserAssignedMap `json:"identityType,omitempty"`
	Reference     *string                            `json:"reference,omitempty"`
	SecretName    *string                            `json:"secretName,omitempty"`
	SecretVersion *string                            `json:"secretVersion,omitempty"`
	Source        *ConfigReferenceSource             `json:"source,omitempty"`
	Status        *ResolveStatus                     `json:"status,omitempty"`
	VaultName     *string                            `json:"vaultName,omitempty"`
}
