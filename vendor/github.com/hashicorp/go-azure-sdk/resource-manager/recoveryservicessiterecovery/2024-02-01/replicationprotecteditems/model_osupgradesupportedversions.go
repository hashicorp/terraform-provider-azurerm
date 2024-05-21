package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSUpgradeSupportedVersions struct {
	SupportedSourceOsVersion  *string   `json:"supportedSourceOsVersion,omitempty"`
	SupportedTargetOsVersions *[]string `json:"supportedTargetOsVersions,omitempty"`
}
