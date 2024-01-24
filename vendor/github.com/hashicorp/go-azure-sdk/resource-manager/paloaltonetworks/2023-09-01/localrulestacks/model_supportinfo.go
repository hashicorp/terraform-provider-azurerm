package localrulestacks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportInfo struct {
	AccountId           *string      `json:"accountId,omitempty"`
	AccountRegistered   *BooleanEnum `json:"accountRegistered,omitempty"`
	FreeTrial           *BooleanEnum `json:"freeTrial,omitempty"`
	FreeTrialCreditLeft *int64       `json:"freeTrialCreditLeft,omitempty"`
	FreeTrialDaysLeft   *int64       `json:"freeTrialDaysLeft,omitempty"`
	HelpURL             *string      `json:"helpURL,omitempty"`
	ProductSerial       *string      `json:"productSerial,omitempty"`
	ProductSku          *string      `json:"productSku,omitempty"`
	RegisterURL         *string      `json:"registerURL,omitempty"`
	SupportURL          *string      `json:"supportURL,omitempty"`
	UserDomainSupported *BooleanEnum `json:"userDomainSupported,omitempty"`
	UserRegistered      *BooleanEnum `json:"userRegistered,omitempty"`
}
