package rolemanagementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApprovalStage struct {
	ApprovalStageTimeOutInDays      *int64     `json:"approvalStageTimeOutInDays,omitempty"`
	EscalationApprovers             *[]UserSet `json:"escalationApprovers,omitempty"`
	EscalationTimeInMinutes         *int64     `json:"escalationTimeInMinutes,omitempty"`
	IsApproverJustificationRequired *bool      `json:"isApproverJustificationRequired,omitempty"`
	IsEscalationEnabled             *bool      `json:"isEscalationEnabled,omitempty"`
	PrimaryApprovers                *[]UserSet `json:"primaryApprovers,omitempty"`
}
