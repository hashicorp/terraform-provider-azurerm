package registries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultProperties struct {
	Identity                 *string `json:"identity,omitempty"`
	KeyIdentifier            *string `json:"keyIdentifier,omitempty"`
	KeyRotationEnabled       *bool   `json:"keyRotationEnabled,omitempty"`
	LastKeyRotationTimestamp *string `json:"lastKeyRotationTimestamp,omitempty"`
	VersionedKeyIdentifier   *string `json:"versionedKeyIdentifier,omitempty"`
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
