package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentSettings struct {
	Enable         *bool     `json:"enable,omitempty"`
	RunImmediately *bool     `json:"runImmediately,omitempty"`
	Schedule       *Schedule `json:"schedule,omitempty"`
}
