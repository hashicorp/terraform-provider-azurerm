package volumegroups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultProperties struct {
	CurrentVersionedKeyExpirationTimestamp *string `json:"currentVersionedKeyExpirationTimestamp,omitempty"`
	CurrentVersionedKeyIdentifier          *string `json:"currentVersionedKeyIdentifier,omitempty"`
	KeyName                                *string `json:"keyName,omitempty"`
	KeyVaultUri                            *string `json:"keyVaultUri,omitempty"`
	KeyVersion                             *string `json:"keyVersion,omitempty"`
	LastKeyRotationTimestamp               *string `json:"lastKeyRotationTimestamp,omitempty"`
}

func (o *KeyVaultProperties) GetCurrentVersionedKeyExpirationTimestampAsTime() (*time.Time, error) {
	if o.CurrentVersionedKeyExpirationTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CurrentVersionedKeyExpirationTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyVaultProperties) SetCurrentVersionedKeyExpirationTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CurrentVersionedKeyExpirationTimestamp = &formatted
}

func (o *KeyVaultProperties) GetLastKeyRotationTimestampAsTime() (*time.Time, error) {
	if o.LastKeyRotationTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastKeyRotationTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyVaultProperties) SetLastKeyRotationTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastKeyRotationTimestamp = &formatted
}
