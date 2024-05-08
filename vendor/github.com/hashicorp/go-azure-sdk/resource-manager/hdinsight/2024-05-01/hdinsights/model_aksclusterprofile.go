package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AksClusterProfile struct {
	AksClusterAgentPoolIdentityProfile *IdentityProfile `json:"aksClusterAgentPoolIdentityProfile,omitempty"`
	AksClusterResourceId               *string          `json:"aksClusterResourceId,omitempty"`
	AksVersion                         *string          `json:"aksVersion,omitempty"`
}
