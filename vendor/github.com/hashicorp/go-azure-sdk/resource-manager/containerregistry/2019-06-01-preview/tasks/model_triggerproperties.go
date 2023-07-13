package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerProperties struct {
	BaseImageTrigger *BaseImageTrigger `json:"baseImageTrigger,omitempty"`
	SourceTriggers   *[]SourceTrigger  `json:"sourceTriggers,omitempty"`
	TimerTriggers    *[]TimerTrigger   `json:"timerTriggers,omitempty"`
}
