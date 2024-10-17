package autonomousdatabasebackups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseBackupUpdateProperties struct {
	RetentionPeriodInDays *int64 `json:"retentionPeriodInDays,omitempty"`
}
