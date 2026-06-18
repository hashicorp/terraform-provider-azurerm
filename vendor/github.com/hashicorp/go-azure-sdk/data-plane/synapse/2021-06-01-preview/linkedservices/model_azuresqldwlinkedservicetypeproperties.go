package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureSqlDWLinkedServiceTypeProperties struct {
	AzureCloudType      *interface{}                  `json:"azureCloudType,omitempty"`
	ConnectionString    interface{}                   `json:"connectionString"`
	EncryptedCredential *interface{}                  `json:"encryptedCredential,omitempty"`
	Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
	ServicePrincipalId  *interface{}                  `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey SecretBase                    `json:"servicePrincipalKey"`
	Tenant              *interface{}                  `json:"tenant,omitempty"`
}

var _ json.Unmarshaler = &AzureSqlDWLinkedServiceTypeProperties{}

func (s *AzureSqlDWLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AzureCloudType      *interface{}                  `json:"azureCloudType,omitempty"`
		ConnectionString    interface{}                   `json:"connectionString"`
		EncryptedCredential *interface{}                  `json:"encryptedCredential,omitempty"`
		Password            *AzureKeyVaultSecretReference `json:"password,omitempty"`
		ServicePrincipalId  *interface{}                  `json:"servicePrincipalId,omitempty"`
		Tenant              *interface{}                  `json:"tenant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AzureCloudType = decoded.AzureCloudType
	s.ConnectionString = decoded.ConnectionString
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Password = decoded.Password
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureSqlDWLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureSqlDWLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
