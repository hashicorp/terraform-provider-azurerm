package encryptionscopes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionScopeKeyVaultProperties struct {
	CurrentVersionedKeyIdentifier *string `json:"currentVersionedKeyIdentifier,omitempty"`
	KeyUri                        *string `json:"keyUri,omitempty"`
	LastKeyRotationTimestamp      *string `json:"lastKeyRotationTimestamp,omitempty"`
}

func (o *EncryptionScopeKeyVaultProperties) GetLastKeyRotationTimestampAsTime() (*time.Time, error) {
	if o.LastKeyRotationTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastKeyRotationTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *EncryptionScopeKeyVaultProperties) SetLastKeyRotationTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastKeyRotationTimestamp = &formatted
}
