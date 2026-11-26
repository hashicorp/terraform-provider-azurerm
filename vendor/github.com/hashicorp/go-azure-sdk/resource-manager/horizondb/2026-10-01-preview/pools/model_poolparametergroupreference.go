package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolParameterGroupReference struct {
	ApplyMethod              ParameterApplyMethod      `json:"applyMethod"`
	ParameterGroupResourceId string                    `json:"parameterGroupResourceId"`
	SyncStatus               *ParameterGroupSyncStatus `json:"syncStatus,omitempty"`
}
