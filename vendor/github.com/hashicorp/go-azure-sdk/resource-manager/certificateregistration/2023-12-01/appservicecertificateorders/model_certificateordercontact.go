package appservicecertificateorders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateOrderContact struct {
	Email     *string `json:"email,omitempty"`
	NameFirst *string `json:"nameFirst,omitempty"`
	NameLast  *string `json:"nameLast,omitempty"`
	Phone     *string `json:"phone,omitempty"`
}
