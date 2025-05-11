package automationrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddIncidentTaskActionProperties struct {
	Description *string `json:"description,omitempty"`
	Title       string  `json:"title"`
}
