package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationJitAccessPolicy struct {
	JitAccessEnabled         bool                     `json:"jitAccessEnabled"`
	JitApprovalMode          *JitApprovalMode         `json:"jitApprovalMode,omitempty"`
	JitApprovers             *[]JitApproverDefinition `json:"jitApprovers,omitempty"`
	MaximumJitAccessDuration *string                  `json:"maximumJitAccessDuration,omitempty"`
}
