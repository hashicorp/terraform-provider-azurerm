package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JitApproverDefinition struct {
	DisplayName *string          `json:"displayName,omitempty"`
	Id          string           `json:"id"`
	Type        *JitApproverType `json:"type,omitempty"`
}
