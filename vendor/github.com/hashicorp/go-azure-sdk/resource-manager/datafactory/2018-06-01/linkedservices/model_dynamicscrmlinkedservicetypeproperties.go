package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DynamicsCrmLinkedServiceTypeProperties struct {
	AuthenticationType             string               `json:"authenticationType"`
	Credential                     *CredentialReference `json:"credential,omitempty"`
	DeploymentType                 string               `json:"deploymentType"`
	Domain                         *string              `json:"domain,omitempty"`
	EncryptedCredential            *string              `json:"encryptedCredential,omitempty"`
	HostName                       *string              `json:"hostName,omitempty"`
	OrganizationName               *string              `json:"organizationName,omitempty"`
	Password                       SecretBase           `json:"password"`
	Port                           *int64               `json:"port,omitempty"`
	ServicePrincipalCredential     SecretBase           `json:"servicePrincipalCredential"`
	ServicePrincipalCredentialType *string              `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalId             *string              `json:"servicePrincipalId,omitempty"`
	ServiceUri                     *string              `json:"serviceUri,omitempty"`
	Username                       *string              `json:"username,omitempty"`
}

var _ json.Unmarshaler = &DynamicsCrmLinkedServiceTypeProperties{}

func (s *DynamicsCrmLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType             string               `json:"authenticationType"`
		Credential                     *CredentialReference `json:"credential,omitempty"`
		DeploymentType                 string               `json:"deploymentType"`
		Domain                         *string              `json:"domain,omitempty"`
		EncryptedCredential            *string              `json:"encryptedCredential,omitempty"`
		HostName                       *string              `json:"hostName,omitempty"`
		OrganizationName               *string              `json:"organizationName,omitempty"`
		Port                           *int64               `json:"port,omitempty"`
		ServicePrincipalCredentialType *string              `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId             *string              `json:"servicePrincipalId,omitempty"`
		ServiceUri                     *string              `json:"serviceUri,omitempty"`
		Username                       *string              `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.Credential = decoded.Credential
	s.DeploymentType = decoded.DeploymentType
	s.Domain = decoded.Domain
	s.EncryptedCredential = decoded.EncryptedCredential
	s.HostName = decoded.HostName
	s.OrganizationName = decoded.OrganizationName
	s.Port = decoded.Port
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.ServiceUri = decoded.ServiceUri
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DynamicsCrmLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'DynamicsCrmLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["servicePrincipalCredential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalCredential' for 'DynamicsCrmLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalCredential = impl
	}

	return nil
}
