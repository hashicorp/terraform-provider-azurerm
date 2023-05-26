package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionDetail struct {
	Detail        *string `json:"Detail,omitempty"`
	MechanismType *string `json:"MechanismType,omitempty"`
	Name          *string `json:"Name,omitempty"`
	SendTime      *string `json:"SendTime,omitempty"`
	Status        *string `json:"Status,omitempty"`
	SubState      *string `json:"SubState,omitempty"`
}
