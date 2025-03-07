package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBatchLinkedServiceTypeProperties struct {
	AccessKey           SecretBase             `json:"accessKey"`
	AccountName         interface{}            `json:"accountName"`
	BatchUri            interface{}            `json:"batchUri"`
	Credential          *CredentialReference   `json:"credential,omitempty"`
	EncryptedCredential *string                `json:"encryptedCredential,omitempty"`
	LinkedServiceName   LinkedServiceReference `json:"linkedServiceName"`
	PoolName            interface{}            `json:"poolName"`
}

var _ json.Unmarshaler = &AzureBatchLinkedServiceTypeProperties{}

func (s *AzureBatchLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountName         interface{}            `json:"accountName"`
		BatchUri            interface{}            `json:"batchUri"`
		Credential          *CredentialReference   `json:"credential,omitempty"`
		EncryptedCredential *string                `json:"encryptedCredential,omitempty"`
		LinkedServiceName   LinkedServiceReference `json:"linkedServiceName"`
		PoolName            interface{}            `json:"poolName"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountName = decoded.AccountName
	s.BatchUri = decoded.BatchUri
	s.Credential = decoded.Credential
	s.EncryptedCredential = decoded.EncryptedCredential
	s.LinkedServiceName = decoded.LinkedServiceName
	s.PoolName = decoded.PoolName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBatchLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accessKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessKey' for 'AzureBatchLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessKey = impl
	}

	return nil
}
