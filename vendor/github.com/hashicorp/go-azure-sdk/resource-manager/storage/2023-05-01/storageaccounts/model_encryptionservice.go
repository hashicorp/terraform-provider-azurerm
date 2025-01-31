package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionService struct {
	Enabled         *bool    `json:"enabled,omitempty"`
	KeyType         *KeyType `json:"keyType,omitempty"`
	LastEnabledTime *string  `json:"lastEnabledTime,omitempty"`
}

func (o *EncryptionService) GetLastEnabledTimeAsTime() (*time.Time, error) {
	if o.LastEnabledTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastEnabledTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EncryptionService) SetLastEnabledTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastEnabledTime = &formatted
}
