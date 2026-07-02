package nginxconfigurationresponses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticItem struct {
	Category    *string `json:"category,omitempty"`
	Description string  `json:"description"`
	Directive   string  `json:"directive"`
	File        string  `json:"file"`
	Id          *string `json:"id,omitempty"`
	Level       Level   `json:"level"`
	Line        float64 `json:"line"`
	Message     string  `json:"message"`
	Rule        string  `json:"rule"`
}
