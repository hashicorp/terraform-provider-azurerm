package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterUpdate struct {
	NodeImageSelection *NodeImageSelection       `json:"nodeImageSelection,omitempty"`
	Upgrade            ManagedClusterUpgradeSpec `json:"upgrade"`
}
