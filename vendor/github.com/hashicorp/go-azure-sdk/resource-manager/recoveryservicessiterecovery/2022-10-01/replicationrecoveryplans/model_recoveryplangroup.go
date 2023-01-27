package replicationrecoveryplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanGroup struct {
	EndGroupActions           *[]RecoveryPlanAction        `json:"endGroupActions,omitempty"`
	GroupType                 RecoveryPlanGroupType        `json:"groupType"`
	ReplicationProtectedItems *[]RecoveryPlanProtectedItem `json:"replicationProtectedItems,omitempty"`
	StartGroupActions         *[]RecoveryPlanAction        `json:"startGroupActions,omitempty"`
}
