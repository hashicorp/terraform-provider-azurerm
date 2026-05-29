package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerUpdateParameters struct {
	BaseImageTrigger *BaseImageTriggerUpdateParameters `json:"baseImageTrigger,omitempty"`
	SourceTriggers   *[]SourceTriggerUpdateParameters  `json:"sourceTriggers,omitempty"`
	TimerTriggers    *[]TimerTriggerUpdateParameters   `json:"timerTriggers,omitempty"`
}
