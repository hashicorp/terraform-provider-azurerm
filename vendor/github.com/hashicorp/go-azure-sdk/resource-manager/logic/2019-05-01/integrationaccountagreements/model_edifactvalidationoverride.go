package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactValidationOverride struct {
	AllowLeadingAndTrailingSpacesAndZeroes bool                    `json:"allowLeadingAndTrailingSpacesAndZeroes"`
	EnforceCharacterSet                    bool                    `json:"enforceCharacterSet"`
	MessageId                              string                  `json:"messageId"`
	TrailingSeparatorPolicy                TrailingSeparatorPolicy `json:"trailingSeparatorPolicy"`
	TrimLeadingAndTrailingSpacesAndZeroes  bool                    `json:"trimLeadingAndTrailingSpacesAndZeroes"`
	ValidateEDITypes                       bool                    `json:"validateEDITypes"`
	ValidateXSDTypes                       bool                    `json:"validateXSDTypes"`
}
