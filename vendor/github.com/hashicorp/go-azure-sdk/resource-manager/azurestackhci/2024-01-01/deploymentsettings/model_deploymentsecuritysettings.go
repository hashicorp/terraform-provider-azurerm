package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSecuritySettings struct {
	BitlockerBootVolume           *bool `json:"bitlockerBootVolume,omitempty"`
	BitlockerDataVolumes          *bool `json:"bitlockerDataVolumes,omitempty"`
	CredentialGuardEnforced       *bool `json:"credentialGuardEnforced,omitempty"`
	DriftControlEnforced          *bool `json:"driftControlEnforced,omitempty"`
	DrtmProtection                *bool `json:"drtmProtection,omitempty"`
	HvciProtection                *bool `json:"hvciProtection,omitempty"`
	SideChannelMitigationEnforced *bool `json:"sideChannelMitigationEnforced,omitempty"`
	SmbClusterEncryption          *bool `json:"smbClusterEncryption,omitempty"`
	SmbSigningEnforced            *bool `json:"smbSigningEnforced,omitempty"`
	WdacEnforced                  *bool `json:"wdacEnforced,omitempty"`
}
