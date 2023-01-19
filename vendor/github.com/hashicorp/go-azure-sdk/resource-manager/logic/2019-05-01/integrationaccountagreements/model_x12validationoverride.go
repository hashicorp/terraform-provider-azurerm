package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12ValidationOverride struct {
	AllowLeadingAndTrailingSpacesAndZeroes bool                    `json:"allowLeadingAndTrailingSpacesAndZeroes"`
	MessageId                              string                  `json:"messageId"`
	TrailingSeparatorPolicy                TrailingSeparatorPolicy `json:"trailingSeparatorPolicy"`
	TrimLeadingAndTrailingSpacesAndZeroes  bool                    `json:"trimLeadingAndTrailingSpacesAndZeroes"`
	ValidateCharacterSet                   bool                    `json:"validateCharacterSet"`
	ValidateEDITypes                       bool                    `json:"validateEDITypes"`
	ValidateXSDTypes                       bool                    `json:"validateXSDTypes"`
}
