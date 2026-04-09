package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceCreatedBy struct {
	UserId    *string `json:"userId,omitempty"`
	UserName  *string `json:"userName,omitempty"`
	UserOrgId *string `json:"userOrgId,omitempty"`
}
