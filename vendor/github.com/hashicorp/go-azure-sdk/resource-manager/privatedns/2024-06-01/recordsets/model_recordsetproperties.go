package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSetProperties struct {
	ARecords         *[]ARecord         `json:"aRecords,omitempty"`
	AaaaRecords      *[]AaaaRecord      `json:"aaaaRecords,omitempty"`
	CnameRecord      *CnameRecord       `json:"cnameRecord,omitempty"`
	Fqdn             *string            `json:"fqdn,omitempty"`
	IsAutoRegistered *bool              `json:"isAutoRegistered,omitempty"`
	Metadata         *map[string]string `json:"metadata,omitempty"`
	MxRecords        *[]MxRecord        `json:"mxRecords,omitempty"`
	PtrRecords       *[]PtrRecord       `json:"ptrRecords,omitempty"`
	SoaRecord        *SoaRecord         `json:"soaRecord,omitempty"`
	SrvRecords       *[]SrvRecord       `json:"srvRecords,omitempty"`
	Ttl              *int64             `json:"ttl,omitempty"`
	TxtRecords       *[]TxtRecord       `json:"txtRecords,omitempty"`
}
