package settings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Setting struct {
	Name  string           `json:"name"`
	Type  *SettingTypeEnum `json:"type,omitempty"`
	Value string           `json:"value"`
}
