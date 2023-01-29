package webtestsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestPropertiesValidationRulesContentValidation struct {
	ContentMatch    *string `json:"ContentMatch,omitempty"`
	IgnoreCase      *bool   `json:"IgnoreCase,omitempty"`
	PassIfTextFound *bool   `json:"PassIfTextFound,omitempty"`
}
