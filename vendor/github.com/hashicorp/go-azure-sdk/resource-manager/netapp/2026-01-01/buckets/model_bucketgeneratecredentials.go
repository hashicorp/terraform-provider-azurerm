package buckets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketGenerateCredentials struct {
	AccessKey     *string `json:"accessKey,omitempty"`
	KeyPairExpiry *string `json:"keyPairExpiry,omitempty"`
	SecretKey     *string `json:"secretKey,omitempty"`
}

func (o *BucketGenerateCredentials) GetKeyPairExpiryAsTime() (*time.Time, error) {
	if o.KeyPairExpiry == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.KeyPairExpiry, "2006-01-02T15:04:05Z07:00")
}

func (o *BucketGenerateCredentials) SetKeyPairExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.KeyPairExpiry = &formatted
}
