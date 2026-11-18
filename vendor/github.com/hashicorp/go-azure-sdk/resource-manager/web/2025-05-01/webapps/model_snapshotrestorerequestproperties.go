package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotRestoreRequestProperties struct {
	IgnoreConflictingHostNames *bool                   `json:"ignoreConflictingHostNames,omitempty"`
	Overwrite                  bool                    `json:"overwrite"`
	RecoverConfiguration       *bool                   `json:"recoverConfiguration,omitempty"`
	RecoverySource             *SnapshotRecoverySource `json:"recoverySource,omitempty"`
	SnapshotTime               *string                 `json:"snapshotTime,omitempty"`
	UseDRSecondary             *bool                   `json:"useDRSecondary,omitempty"`
}
