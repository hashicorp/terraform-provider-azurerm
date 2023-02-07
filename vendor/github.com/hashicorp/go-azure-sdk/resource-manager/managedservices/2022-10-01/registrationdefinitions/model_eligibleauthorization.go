package registrationdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EligibleAuthorization struct {
	JustInTimeAccessPolicy *JustInTimeAccessPolicy `json:"justInTimeAccessPolicy,omitempty"`
	PrincipalId            string                  `json:"principalId"`
	PrincipalIdDisplayName *string                 `json:"principalIdDisplayName,omitempty"`
	RoleDefinitionId       string                  `json:"roleDefinitionId"`
}
