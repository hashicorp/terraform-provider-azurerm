package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SasDefinitionAttributes struct {
	Created         *int64                 `json:"created,omitempty"`
	Enabled         *bool                  `json:"enabled,omitempty"`
	RecoverableDays *int64                 `json:"recoverableDays,omitempty"`
	RecoveryLevel   *DeletionRecoveryLevel `json:"recoveryLevel,omitempty"`
	Updated         *int64                 `json:"updated,omitempty"`
}
