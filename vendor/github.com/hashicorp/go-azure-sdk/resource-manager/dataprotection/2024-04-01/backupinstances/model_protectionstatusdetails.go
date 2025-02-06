package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionStatusDetails struct {
	ErrorDetails *UserFacingError `json:"errorDetails,omitempty"`
	Status       *Status          `json:"status,omitempty"`
}
