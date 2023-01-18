package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityProfile struct {
	OpenAccess       *EnableState `json:"openAccess,omitempty"`
	RegistrationCode *string      `json:"registrationCode,omitempty"`
}
