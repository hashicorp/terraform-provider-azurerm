package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlotDifferenceProperties struct {
	Description        *string `json:"description,omitempty"`
	DiffRule           *string `json:"diffRule,omitempty"`
	Level              *string `json:"level,omitempty"`
	SettingName        *string `json:"settingName,omitempty"`
	SettingType        *string `json:"settingType,omitempty"`
	ValueInCurrentSlot *string `json:"valueInCurrentSlot,omitempty"`
	ValueInTargetSlot  *string `json:"valueInTargetSlot,omitempty"`
}
