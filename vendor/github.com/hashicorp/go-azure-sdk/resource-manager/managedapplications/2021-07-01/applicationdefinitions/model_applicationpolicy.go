package applicationdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationPolicy struct {
	Name               *string `json:"name,omitempty"`
	Parameters         *string `json:"parameters,omitempty"`
	PolicyDefinitionId *string `json:"policyDefinitionId,omitempty"`
}
