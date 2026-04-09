package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpgradeGraphProperties struct {
	ApplianceVersion  *string             `json:"applianceVersion,omitempty"`
	SupportedVersions *[]SupportedVersion `json:"supportedVersions,omitempty"`
}
