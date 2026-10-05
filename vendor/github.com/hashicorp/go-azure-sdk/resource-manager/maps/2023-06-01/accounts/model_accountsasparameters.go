package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountSasParameters struct {
	Expiry           string     `json:"expiry"`
	MaxRatePerSecond int64      `json:"maxRatePerSecond"`
	PrincipalId      string     `json:"principalId"`
	Regions          *[]string  `json:"regions,omitempty"`
	SigningKey       SigningKey `json:"signingKey"`
	Start            string     `json:"start"`
}
