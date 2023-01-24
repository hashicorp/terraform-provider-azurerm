package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12ValidationSettings struct {
	AllowLeadingAndTrailingSpacesAndZeroes    bool                    `json:"allowLeadingAndTrailingSpacesAndZeroes"`
	CheckDuplicateGroupControlNumber          bool                    `json:"checkDuplicateGroupControlNumber"`
	CheckDuplicateInterchangeControlNumber    bool                    `json:"checkDuplicateInterchangeControlNumber"`
	CheckDuplicateTransactionSetControlNumber bool                    `json:"checkDuplicateTransactionSetControlNumber"`
	InterchangeControlNumberValidityDays      int64                   `json:"interchangeControlNumberValidityDays"`
	TrailingSeparatorPolicy                   TrailingSeparatorPolicy `json:"trailingSeparatorPolicy"`
	TrimLeadingAndTrailingSpacesAndZeroes     bool                    `json:"trimLeadingAndTrailingSpacesAndZeroes"`
	ValidateCharacterSet                      bool                    `json:"validateCharacterSet"`
	ValidateEDITypes                          bool                    `json:"validateEDITypes"`
	ValidateXSDTypes                          bool                    `json:"validateXSDTypes"`
}
