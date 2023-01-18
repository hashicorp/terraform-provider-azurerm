package publicmaintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputLinuxParameters struct {
	ClassificationsToInclude  *[]string `json:"classificationsToInclude,omitempty"`
	PackageNameMasksToExclude *[]string `json:"packageNameMasksToExclude,omitempty"`
	PackageNameMasksToInclude *[]string `json:"packageNameMasksToInclude,omitempty"`
}
