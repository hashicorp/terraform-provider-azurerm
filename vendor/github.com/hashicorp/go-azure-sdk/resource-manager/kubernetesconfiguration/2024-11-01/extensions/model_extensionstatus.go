package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionStatus struct {
	Code          *string    `json:"code,omitempty"`
	DisplayStatus *string    `json:"displayStatus,omitempty"`
	Level         *LevelType `json:"level,omitempty"`
	Message       *string    `json:"message,omitempty"`
	Time          *string    `json:"time,omitempty"`
}
