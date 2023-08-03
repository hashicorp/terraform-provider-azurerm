package recoverableservers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoverableServerProperties struct {
	Edition                     *string `json:"edition,omitempty"`
	HardwareGeneration          *string `json:"hardwareGeneration,omitempty"`
	LastAvailableBackupDateTime *string `json:"lastAvailableBackupDateTime,omitempty"`
	ServiceLevelObjective       *string `json:"serviceLevelObjective,omitempty"`
	VCore                       *int64  `json:"vCore,omitempty"`
	Version                     *string `json:"version,omitempty"`
}
