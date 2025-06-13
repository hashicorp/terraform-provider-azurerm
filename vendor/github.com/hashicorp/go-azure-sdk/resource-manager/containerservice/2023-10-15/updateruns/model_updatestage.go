package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateStage struct {
	AfterStageWaitInSeconds *int64         `json:"afterStageWaitInSeconds,omitempty"`
	Groups                  *[]UpdateGroup `json:"groups,omitempty"`
	Name                    string         `json:"name"`
}
