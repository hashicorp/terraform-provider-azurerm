package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyAttributes struct {
	Attestation     *KeyAttestation        `json:"attestation,omitempty"`
	Created         *int64                 `json:"created,omitempty"`
	Enabled         *bool                  `json:"enabled,omitempty"`
	Exp             *int64                 `json:"exp,omitempty"`
	Exportable      *bool                  `json:"exportable,omitempty"`
	HsmPlatform     *string                `json:"hsmPlatform,omitempty"`
	Nbf             *int64                 `json:"nbf,omitempty"`
	RecoverableDays *int64                 `json:"recoverableDays,omitempty"`
	RecoveryLevel   *DeletionRecoveryLevel `json:"recoveryLevel,omitempty"`
	Updated         *int64                 `json:"updated,omitempty"`
}
