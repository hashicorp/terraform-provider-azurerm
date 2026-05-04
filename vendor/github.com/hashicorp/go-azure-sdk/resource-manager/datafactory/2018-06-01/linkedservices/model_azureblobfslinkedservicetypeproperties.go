package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobFSLinkedServiceTypeProperties struct {
	AccountKey                     *interface{}         `json:"accountKey,omitempty"`
	AzureCloudType                 *interface{}         `json:"azureCloudType,omitempty"`
	Credential                     *CredentialReference `json:"credential,omitempty"`
	EncryptedCredential            *string              `json:"encryptedCredential,omitempty"`
	SasToken                       SecretBase           `json:"sasToken"`
	SasUri                         *interface{}         `json:"sasUri,omitempty"`
	ServicePrincipalCredential     SecretBase           `json:"servicePrincipalCredential"`
	ServicePrincipalCredentialType *interface{}         `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalId             *interface{}         `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey            SecretBase           `json:"servicePrincipalKey"`
	Tenant                         *interface{}         `json:"tenant,omitempty"`
	Url                            *interface{}         `json:"url,omitempty"`
}

var _ json.Unmarshaler = &AzureBlobFSLinkedServiceTypeProperties{}

func (s *AzureBlobFSLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountKey                     *interface{}         `json:"accountKey,omitempty"`
		AzureCloudType                 *interface{}         `json:"azureCloudType,omitempty"`
		Credential                     *CredentialReference `json:"credential,omitempty"`
		EncryptedCredential            *string              `json:"encryptedCredential,omitempty"`
		SasUri                         *interface{}         `json:"sasUri,omitempty"`
		ServicePrincipalCredentialType *interface{}         `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId             *interface{}         `json:"servicePrincipalId,omitempty"`
		Tenant                         *interface{}         `json:"tenant,omitempty"`
		Url                            *interface{}         `json:"url,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountKey = decoded.AccountKey
	s.AzureCloudType = decoded.AzureCloudType
	s.Credential = decoded.Credential
	s.EncryptedCredential = decoded.EncryptedCredential
	s.SasUri = decoded.SasUri
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBlobFSLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["sasToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SasToken' for 'AzureBlobFSLinkedServiceTypeProperties': %+v", err)
		}
		s.SasToken = impl
	}

	if v, ok := temp["servicePrincipalCredential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalCredential' for 'AzureBlobFSLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalCredential = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureBlobFSLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
