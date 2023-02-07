package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeviceSecrets struct {
	BmcDefaultUserPassword                *Secret `json:"bmcDefaultUserPassword,omitempty"`
	HcsDataVolumeBitLockerExternalKey     *Secret `json:"hcsDataVolumeBitLockerExternalKey,omitempty"`
	HcsInternalVolumeBitLockerExternalKey *Secret `json:"hcsInternalVolumeBitLockerExternalKey,omitempty"`
	RotateKeyForDataVolumeBitlocker       *Secret `json:"rotateKeyForDataVolumeBitlocker,omitempty"`
	RotateKeysForSedDrivesSerialized      *Secret `json:"rotateKeysForSedDrivesSerialized,omitempty"`
	SedEncryptionExternalKey              *Secret `json:"sedEncryptionExternalKey,omitempty"`
	SedEncryptionExternalKeyId            *Secret `json:"sedEncryptionExternalKeyId,omitempty"`
	SystemVolumeBitLockerRecoveryKey      *Secret `json:"systemVolumeBitLockerRecoveryKey,omitempty"`
}
