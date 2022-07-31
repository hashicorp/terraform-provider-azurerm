package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSetProperties struct {
	AAAARecords       *[]AaaaRecord      `json:"AAAARecords,omitempty"`
	ARecords          *[]ARecord         `json:"ARecords,omitempty"`
	CNAMERecord       *CnameRecord       `json:"CNAMERecord,omitempty"`
	CaaRecords        *[]CaaRecord       `json:"caaRecords,omitempty"`
	Fqdn              *string            `json:"fqdn,omitempty"`
	MXRecords         *[]MxRecord        `json:"MXRecords,omitempty"`
	Metadata          *map[string]string `json:"metadata,omitempty"`
	NSRecords         *[]NsRecord        `json:"NSRecords,omitempty"`
	PTRRecords        *[]PtrRecord       `json:"PTRRecords,omitempty"`
	ProvisioningState *string            `json:"provisioningState,omitempty"`
	SOARecord         *SoaRecord         `json:"SOARecord,omitempty"`
	SRVRecords        *[]SrvRecord       `json:"SRVRecords,omitempty"`
	TTL               *int64             `json:"TTL,omitempty"`
	TXTRecords        *[]TxtRecord       `json:"TXTRecords,omitempty"`
	TargetResource    *SubResource       `json:"targetResource,omitempty"`
}
