package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeImageSelection struct {
	CustomNodeImageVersions *[]NodeImageVersion    `json:"customNodeImageVersions,omitempty"`
	Type                    NodeImageSelectionType `json:"type"`
}
