package certificate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteCertificateError struct {
	Code    string                    `json:"code"`
	Details *[]DeleteCertificateError `json:"details,omitempty"`
	Message string                    `json:"message"`
	Target  *string                   `json:"target,omitempty"`
}
