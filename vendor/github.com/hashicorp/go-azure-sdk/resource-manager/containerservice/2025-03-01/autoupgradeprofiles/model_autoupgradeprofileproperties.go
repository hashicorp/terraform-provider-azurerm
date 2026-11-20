package autoupgradeprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoUpgradeProfileProperties struct {
	AutoUpgradeProfileStatus *AutoUpgradeProfileStatus            `json:"autoUpgradeProfileStatus,omitempty"`
	Channel                  UpgradeChannel                       `json:"channel"`
	Disabled                 *bool                                `json:"disabled,omitempty"`
	NodeImageSelection       *AutoUpgradeNodeImageSelection       `json:"nodeImageSelection,omitempty"`
	ProvisioningState        *AutoUpgradeProfileProvisioningState `json:"provisioningState,omitempty"`
	UpdateStrategyId         *string                              `json:"updateStrategyId,omitempty"`
}
