package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentSettingsSpec struct {
	Logs *[]AgentSetting `json:"logs,omitempty"`
}
