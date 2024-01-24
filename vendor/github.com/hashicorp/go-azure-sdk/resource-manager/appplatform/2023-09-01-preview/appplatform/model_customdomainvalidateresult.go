package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainValidateResult struct {
	IsValid *bool   `json:"isValid,omitempty"`
	Message *string `json:"message,omitempty"`
}
