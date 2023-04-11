package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceTriggerUpdateParameters struct {
	Name                string                  `json:"name"`
	SourceRepository    *SourceUpdateParameters `json:"sourceRepository,omitempty"`
	SourceTriggerEvents *[]SourceTriggerEvent   `json:"sourceTriggerEvents,omitempty"`
	Status              *TriggerStatus          `json:"status,omitempty"`
}
