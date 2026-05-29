package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceTrigger struct {
	Name                string               `json:"name"`
	SourceRepository    SourceProperties     `json:"sourceRepository"`
	SourceTriggerEvents []SourceTriggerEvent `json:"sourceTriggerEvents"`
	Status              *TriggerStatus       `json:"status,omitempty"`
}
