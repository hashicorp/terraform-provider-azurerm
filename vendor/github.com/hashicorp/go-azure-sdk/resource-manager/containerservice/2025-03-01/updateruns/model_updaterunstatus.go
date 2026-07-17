package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateRunStatus struct {
	NodeImageSelection *NodeImageSelectionStatus `json:"nodeImageSelection,omitempty"`
	Stages             *[]UpdateStageStatus      `json:"stages,omitempty"`
	Status             *UpdateStatus             `json:"status,omitempty"`
}
