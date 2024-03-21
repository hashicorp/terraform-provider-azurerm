package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Actions struct {
	ActionGroups     *[]string          `json:"actionGroups,omitempty"`
	ActionProperties *map[string]string `json:"actionProperties,omitempty"`
	CustomProperties *map[string]string `json:"customProperties,omitempty"`
}
