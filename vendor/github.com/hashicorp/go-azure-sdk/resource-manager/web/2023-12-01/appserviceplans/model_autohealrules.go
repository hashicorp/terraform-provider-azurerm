package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoHealRules struct {
	Actions  *AutoHealActions  `json:"actions,omitempty"`
	Triggers *AutoHealTriggers `json:"triggers,omitempty"`
}
