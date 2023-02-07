package diskencryptionsets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionSetProperties struct {
	ActiveKey                         *KeyForDiskEncryptionSet   `json:"activeKey,omitempty"`
	AutoKeyRotationError              *ApiError                  `json:"autoKeyRotationError,omitempty"`
	EncryptionType                    *DiskEncryptionSetType     `json:"encryptionType,omitempty"`
	FederatedClientId                 *string                    `json:"federatedClientId,omitempty"`
	LastKeyRotationTimestamp          *string                    `json:"lastKeyRotationTimestamp,omitempty"`
	PreviousKeys                      *[]KeyForDiskEncryptionSet `json:"previousKeys,omitempty"`
	ProvisioningState                 *string                    `json:"provisioningState,omitempty"`
	RotationToLatestKeyVersionEnabled *bool                      `json:"rotationToLatestKeyVersionEnabled,omitempty"`
}

func (o *EncryptionSetProperties) GetLastKeyRotationTimestampAsTime() (*time.Time, error) {
	if o.LastKeyRotationTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastKeyRotationTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *EncryptionSetProperties) SetLastKeyRotationTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastKeyRotationTimestamp = &formatted
}
