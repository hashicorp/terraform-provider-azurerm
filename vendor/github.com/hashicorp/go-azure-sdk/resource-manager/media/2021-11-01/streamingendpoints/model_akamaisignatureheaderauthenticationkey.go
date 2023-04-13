package streamingendpoints

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AkamaiSignatureHeaderAuthenticationKey struct {
	Base64Key  *string `json:"base64Key,omitempty"`
	Expiration *string `json:"expiration,omitempty"`
	Identifier *string `json:"identifier,omitempty"`
}

func (o *AkamaiSignatureHeaderAuthenticationKey) GetExpirationAsTime() (*time.Time, error) {
	if o.Expiration == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Expiration, "2006-01-02T15:04:05Z07:00")
}

func (o *AkamaiSignatureHeaderAuthenticationKey) SetExpirationAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Expiration = &formatted
}
