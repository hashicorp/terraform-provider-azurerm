package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptActionProfile struct {
	Name             string   `json:"name"`
	Parameters       *string  `json:"parameters,omitempty"`
	Services         []string `json:"services"`
	ShouldPersist    *bool    `json:"shouldPersist,omitempty"`
	TimeoutInMinutes *int64   `json:"timeoutInMinutes,omitempty"`
	Type             string   `json:"type"`
	Url              string   `json:"url"`
}
