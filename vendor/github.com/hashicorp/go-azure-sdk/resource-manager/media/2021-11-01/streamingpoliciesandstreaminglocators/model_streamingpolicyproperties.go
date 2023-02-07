package streamingpoliciesandstreaminglocators

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPolicyProperties struct {
	CommonEncryptionCbcs        *CommonEncryptionCbcs `json:"commonEncryptionCbcs,omitempty"`
	CommonEncryptionCenc        *CommonEncryptionCenc `json:"commonEncryptionCenc,omitempty"`
	Created                     *string               `json:"created,omitempty"`
	DefaultContentKeyPolicyName *string               `json:"defaultContentKeyPolicyName,omitempty"`
	EnvelopeEncryption          *EnvelopeEncryption   `json:"envelopeEncryption,omitempty"`
	NoEncryption                *NoEncryption         `json:"noEncryption,omitempty"`
}

func (o *StreamingPolicyProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingPolicyProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}
