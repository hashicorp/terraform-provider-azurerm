package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileStorageLinkedServiceTypeProperties struct {
	AccountKey          *AzureKeyVaultSecretReference `json:"accountKey,omitempty"`
	ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
	Credential          *CredentialReference          `json:"credential,omitempty"`
	EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
	FileShare           *interface{}                  `json:"fileShare,omitempty"`
	Host                *interface{}                  `json:"host,omitempty"`
	Password            SecretBase                    `json:"password"`
	SasToken            *AzureKeyVaultSecretReference `json:"sasToken,omitempty"`
	SasUri              *interface{}                  `json:"sasUri,omitempty"`
	ServiceEndpoint     *interface{}                  `json:"serviceEndpoint,omitempty"`
	Snapshot            *interface{}                  `json:"snapshot,omitempty"`
	UserId              *interface{}                  `json:"userId,omitempty"`
}

var _ json.Unmarshaler = &AzureFileStorageLinkedServiceTypeProperties{}

func (s *AzureFileStorageLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountKey          *AzureKeyVaultSecretReference `json:"accountKey,omitempty"`
		ConnectionString    *interface{}                  `json:"connectionString,omitempty"`
		Credential          *CredentialReference          `json:"credential,omitempty"`
		EncryptedCredential *string                       `json:"encryptedCredential,omitempty"`
		FileShare           *interface{}                  `json:"fileShare,omitempty"`
		Host                *interface{}                  `json:"host,omitempty"`
		SasToken            *AzureKeyVaultSecretReference `json:"sasToken,omitempty"`
		SasUri              *interface{}                  `json:"sasUri,omitempty"`
		ServiceEndpoint     *interface{}                  `json:"serviceEndpoint,omitempty"`
		Snapshot            *interface{}                  `json:"snapshot,omitempty"`
		UserId              *interface{}                  `json:"userId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountKey = decoded.AccountKey
	s.ConnectionString = decoded.ConnectionString
	s.Credential = decoded.Credential
	s.EncryptedCredential = decoded.EncryptedCredential
	s.FileShare = decoded.FileShare
	s.Host = decoded.Host
	s.SasToken = decoded.SasToken
	s.SasUri = decoded.SasUri
	s.ServiceEndpoint = decoded.ServiceEndpoint
	s.Snapshot = decoded.Snapshot
	s.UserId = decoded.UserId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureFileStorageLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'AzureFileStorageLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
