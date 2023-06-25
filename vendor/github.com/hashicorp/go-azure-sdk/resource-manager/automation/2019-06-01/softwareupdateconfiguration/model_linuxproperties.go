package softwareupdateconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxProperties struct {
	ExcludedPackageNameMasks       *[]string           `json:"excludedPackageNameMasks,omitempty"`
	IncludedPackageClassifications *LinuxUpdateClasses `json:"includedPackageClassifications,omitempty"`
	IncludedPackageNameMasks       *[]string           `json:"includedPackageNameMasks,omitempty"`
	RebootSetting                  *string             `json:"rebootSetting,omitempty"`
}
