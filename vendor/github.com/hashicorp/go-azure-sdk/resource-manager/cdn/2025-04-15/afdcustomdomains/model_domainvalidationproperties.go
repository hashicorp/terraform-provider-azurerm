package afdcustomdomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainValidationProperties struct {
	ExpirationDate  *string `json:"expirationDate,omitempty"`
	ValidationToken *string `json:"validationToken,omitempty"`
}
