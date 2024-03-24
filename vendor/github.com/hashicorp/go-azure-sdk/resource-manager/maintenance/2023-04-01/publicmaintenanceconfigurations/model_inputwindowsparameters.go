package publicmaintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputWindowsParameters struct {
	ClassificationsToInclude  *[]string `json:"classificationsToInclude,omitempty"`
	ExcludeKbsRequiringReboot *bool     `json:"excludeKbsRequiringReboot,omitempty"`
	KbNumbersToExclude        *[]string `json:"kbNumbersToExclude,omitempty"`
	KbNumbersToInclude        *[]string `json:"kbNumbersToInclude,omitempty"`
}
