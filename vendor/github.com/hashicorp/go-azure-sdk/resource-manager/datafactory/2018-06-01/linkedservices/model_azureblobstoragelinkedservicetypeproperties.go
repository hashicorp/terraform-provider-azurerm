package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobStorageLinkedServiceTypeProperties struct {
	AccountKey          *AzureKeyVaultSecretReference   `json:"accountKey,omitempty"`
	AccountKind         *interface{}                    `json:"accountKind,omitempty"`
	AuthenticationType  *AzureStorageAuthenticationType `json:"authenticationType,omitempty"`
	AzureCloudType      *interface{}                    `json:"azureCloudType,omitempty"`
	ConnectionString    *interface{}                    `json:"connectionString,omitempty"`
	ContainerUri        *interface{}                    `json:"containerUri,omitempty"`
	Credential          *CredentialReference            `json:"credential,omitempty"`
	EncryptedCredential *string                         `json:"encryptedCredential,omitempty"`
	SasToken            *AzureKeyVaultSecretReference   `json:"sasToken,omitempty"`
	SasUri              *interface{}                    `json:"sasUri,omitempty"`
	ServiceEndpoint     *interface{}                    `json:"serviceEndpoint,omitempty"`
	ServicePrincipalId  *interface{}                    `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey SecretBase                      `json:"servicePrincipalKey"`
	Tenant              *interface{}                    `json:"tenant,omitempty"`
}

var _ json.Unmarshaler = &AzureBlobStorageLinkedServiceTypeProperties{}

func (s *AzureBlobStorageLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountKey          *AzureKeyVaultSecretReference   `json:"accountKey,omitempty"`
		AccountKind         *interface{}                    `json:"accountKind,omitempty"`
		AuthenticationType  *AzureStorageAuthenticationType `json:"authenticationType,omitempty"`
		AzureCloudType      *interface{}                    `json:"azureCloudType,omitempty"`
		ConnectionString    *interface{}                    `json:"connectionString,omitempty"`
		ContainerUri        *interface{}                    `json:"containerUri,omitempty"`
		Credential          *CredentialReference            `json:"credential,omitempty"`
		EncryptedCredential *string                         `json:"encryptedCredential,omitempty"`
		SasToken            *AzureKeyVaultSecretReference   `json:"sasToken,omitempty"`
		SasUri              *interface{}                    `json:"sasUri,omitempty"`
		ServiceEndpoint     *interface{}                    `json:"serviceEndpoint,omitempty"`
		ServicePrincipalId  *interface{}                    `json:"servicePrincipalId,omitempty"`
		Tenant              *interface{}                    `json:"tenant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountKey = decoded.AccountKey
	s.AccountKind = decoded.AccountKind
	s.AuthenticationType = decoded.AuthenticationType
	s.AzureCloudType = decoded.AzureCloudType
	s.ConnectionString = decoded.ConnectionString
	s.ContainerUri = decoded.ContainerUri
	s.Credential = decoded.Credential
	s.EncryptedCredential = decoded.EncryptedCredential
	s.SasToken = decoded.SasToken
	s.SasUri = decoded.SasUri
	s.ServiceEndpoint = decoded.ServiceEndpoint
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBlobStorageLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureBlobStorageLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
