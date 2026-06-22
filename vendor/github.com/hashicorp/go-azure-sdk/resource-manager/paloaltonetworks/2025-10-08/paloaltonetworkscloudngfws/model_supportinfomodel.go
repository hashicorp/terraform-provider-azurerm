package paloaltonetworkscloudngfws

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportInfoModel struct {
	AccountId                 *string             `json:"accountId,omitempty"`
	AccountRegistrationStatus *RegistrationStatus `json:"accountRegistrationStatus,omitempty"`
	Credits                   *int64              `json:"credits,omitempty"`
	EndDateForCredits         *string             `json:"endDateForCredits,omitempty"`
	FreeTrial                 *EnableStatus       `json:"freeTrial,omitempty"`
	FreeTrialCreditLeft       *int64              `json:"freeTrialCreditLeft,omitempty"`
	FreeTrialDaysLeft         *int64              `json:"freeTrialDaysLeft,omitempty"`
	HelpURL                   *string             `json:"helpURL,omitempty"`
	HubURL                    *string             `json:"hubUrl,omitempty"`
	MonthlyCreditLeft         *int64              `json:"monthlyCreditLeft,omitempty"`
	ProductSerial             *string             `json:"productSerial,omitempty"`
	ProductSku                *string             `json:"productSku,omitempty"`
	RegisterURL               *string             `json:"registerURL,omitempty"`
	StartDateForCredits       *string             `json:"startDateForCredits,omitempty"`
	SupportURL                *string             `json:"supportURL,omitempty"`
}
