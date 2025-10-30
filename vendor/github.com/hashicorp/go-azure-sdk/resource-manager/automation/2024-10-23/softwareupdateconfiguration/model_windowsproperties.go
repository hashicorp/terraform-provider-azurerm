package softwareupdateconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WindowsProperties struct {
	ExcludedKbNumbers             *[]string             `json:"excludedKbNumbers,omitempty"`
	IncludedKbNumbers             *[]string             `json:"includedKbNumbers,omitempty"`
	IncludedUpdateClassifications *WindowsUpdateClasses `json:"includedUpdateClassifications,omitempty"`
	RebootSetting                 *string               `json:"rebootSetting,omitempty"`
}
