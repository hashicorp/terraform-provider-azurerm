package batchaccount

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoStorageProperties struct {
	AuthenticationMode    *AutoStorageAuthenticationMode `json:"authenticationMode,omitempty"`
	LastKeySync           string                         `json:"lastKeySync"`
	NodeIdentityReference *ComputeNodeIdentityReference  `json:"nodeIdentityReference,omitempty"`
	StorageAccountId      string                         `json:"storageAccountId"`
}

func (o *AutoStorageProperties) GetLastKeySyncAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.LastKeySync, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoStorageProperties) SetLastKeySyncAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastKeySync = formatted
}
