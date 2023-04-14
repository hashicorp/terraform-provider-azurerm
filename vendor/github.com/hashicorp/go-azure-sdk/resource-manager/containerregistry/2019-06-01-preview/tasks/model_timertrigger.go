package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimerTrigger struct {
	Name     string         `json:"name"`
	Schedule string         `json:"schedule"`
	Status   *TriggerStatus `json:"status,omitempty"`
}
