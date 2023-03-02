package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateCustomDomainOutput struct {
	CustomDomainValidated *bool   `json:"customDomainValidated,omitempty"`
	Message               *string `json:"message,omitempty"`
	Reason                *string `json:"reason,omitempty"`
}
