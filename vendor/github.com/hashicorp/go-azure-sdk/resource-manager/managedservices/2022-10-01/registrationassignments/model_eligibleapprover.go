package registrationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EligibleApprover struct {
	PrincipalId            string  `json:"principalId"`
	PrincipalIdDisplayName *string `json:"principalIdDisplayName,omitempty"`
}
