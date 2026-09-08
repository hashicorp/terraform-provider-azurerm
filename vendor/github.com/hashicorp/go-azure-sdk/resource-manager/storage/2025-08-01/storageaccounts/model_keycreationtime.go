package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyCreationTime struct {
	Key1 *string `json:"key1,omitempty"`
	Key2 *string `json:"key2,omitempty"`
}

func (o *KeyCreationTime) GetKey1AsTime() (*time.Time, error) {
	if o.Key1 == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Key1, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyCreationTime) SetKey1AsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Key1 = &formatted
}

func (o *KeyCreationTime) GetKey2AsTime() (*time.Time, error) {
	if o.Key2 == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Key2, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyCreationTime) SetKey2AsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Key2 = &formatted
}
