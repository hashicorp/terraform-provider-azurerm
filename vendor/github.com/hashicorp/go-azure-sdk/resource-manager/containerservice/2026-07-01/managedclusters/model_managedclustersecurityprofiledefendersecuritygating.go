package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSecurityProfileDefenderSecurityGating struct {
	AllowSecretAccess *bool                                                          `json:"allowSecretAccess,omitempty"`
	Enabled           *bool                                                          `json:"enabled,omitempty"`
	Identities        *[]ManagedClusterSecurityProfileDefenderSecurityGatingIdentity `json:"identities,omitempty"`
}
