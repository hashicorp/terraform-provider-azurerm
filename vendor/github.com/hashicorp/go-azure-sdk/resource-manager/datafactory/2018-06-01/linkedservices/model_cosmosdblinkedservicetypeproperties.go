package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CosmosDbLinkedServiceTypeProperties struct {
	AccountEndpoint                *interface{}            `json:"accountEndpoint,omitempty"`
	AccountKey                     SecretBase              `json:"accountKey"`
	AzureCloudType                 *interface{}            `json:"azureCloudType,omitempty"`
	ConnectionMode                 *CosmosDbConnectionMode `json:"connectionMode,omitempty"`
	ConnectionString               *interface{}            `json:"connectionString,omitempty"`
	Credential                     *CredentialReference    `json:"credential,omitempty"`
	Database                       *interface{}            `json:"database,omitempty"`
	EncryptedCredential            *string                 `json:"encryptedCredential,omitempty"`
	ServicePrincipalCredential     SecretBase              `json:"servicePrincipalCredential"`
	ServicePrincipalCredentialType *interface{}            `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalId             *interface{}            `json:"servicePrincipalId,omitempty"`
	Tenant                         *interface{}            `json:"tenant,omitempty"`
}

var _ json.Unmarshaler = &CosmosDbLinkedServiceTypeProperties{}

func (s *CosmosDbLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountEndpoint                *interface{}            `json:"accountEndpoint,omitempty"`
		AzureCloudType                 *interface{}            `json:"azureCloudType,omitempty"`
		ConnectionMode                 *CosmosDbConnectionMode `json:"connectionMode,omitempty"`
		ConnectionString               *interface{}            `json:"connectionString,omitempty"`
		Credential                     *CredentialReference    `json:"credential,omitempty"`
		Database                       *interface{}            `json:"database,omitempty"`
		EncryptedCredential            *string                 `json:"encryptedCredential,omitempty"`
		ServicePrincipalCredentialType *interface{}            `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId             *interface{}            `json:"servicePrincipalId,omitempty"`
		Tenant                         *interface{}            `json:"tenant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountEndpoint = decoded.AccountEndpoint
	s.AzureCloudType = decoded.AzureCloudType
	s.ConnectionMode = decoded.ConnectionMode
	s.ConnectionString = decoded.ConnectionString
	s.Credential = decoded.Credential
	s.Database = decoded.Database
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CosmosDbLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accountKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccountKey' for 'CosmosDbLinkedServiceTypeProperties': %+v", err)
		}
		s.AccountKey = impl
	}

	if v, ok := temp["servicePrincipalCredential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalCredential' for 'CosmosDbLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalCredential = impl
	}

	return nil
}
