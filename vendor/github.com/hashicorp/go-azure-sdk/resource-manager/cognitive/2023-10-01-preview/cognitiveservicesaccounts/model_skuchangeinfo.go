package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuChangeInfo struct {
	CountOfDowngrades              *float64 `json:"countOfDowngrades,omitempty"`
	CountOfUpgradesAfterDowngrades *float64 `json:"countOfUpgradesAfterDowngrades,omitempty"`
	LastChangeDate                 *string  `json:"lastChangeDate,omitempty"`
}
