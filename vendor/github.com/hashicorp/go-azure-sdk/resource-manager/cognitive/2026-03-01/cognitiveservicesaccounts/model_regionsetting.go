package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionSetting struct {
	Customsubdomain *string  `json:"customsubdomain,omitempty"`
	Name            *string  `json:"name,omitempty"`
	Value           *float64 `json:"value,omitempty"`
}
