package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlatformWorkloadIdentityProfile struct {
	PlatformWorkloadIdentities *map[string]PlatformWorkloadIdentity `json:"platformWorkloadIdentities,omitempty"`
	UpgradeableTo              *string                              `json:"upgradeableTo,omitempty"`
}
