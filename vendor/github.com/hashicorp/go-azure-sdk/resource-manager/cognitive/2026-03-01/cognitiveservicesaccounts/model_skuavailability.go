package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuAvailability struct {
	Kind         *string `json:"kind,omitempty"`
	Message      *string `json:"message,omitempty"`
	Reason       *string `json:"reason,omitempty"`
	SkuAvailable *bool   `json:"skuAvailable,omitempty"`
	SkuName      *string `json:"skuName,omitempty"`
	Type         *string `json:"type,omitempty"`
}
