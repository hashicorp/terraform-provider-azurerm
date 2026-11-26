package externalauths

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenClaimMappingsProfileUpdate struct {
	Groups   *GroupClaimProfileUpdate    `json:"groups,omitempty"`
	Username *UsernameClaimProfileUpdate `json:"username,omitempty"`
}
