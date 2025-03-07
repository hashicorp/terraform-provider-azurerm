package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestServiceLinkedServiceTypeProperties struct {
	AadResourceId                        *interface{}                  `json:"aadResourceId,omitempty"`
	AuthHeaders                          *interface{}                  `json:"authHeaders,omitempty"`
	AuthenticationType                   RestServiceAuthenticationType `json:"authenticationType"`
	AzureCloudType                       *interface{}                  `json:"azureCloudType,omitempty"`
	ClientId                             *interface{}                  `json:"clientId,omitempty"`
	ClientSecret                         SecretBase                    `json:"clientSecret"`
	Credential                           *CredentialReference          `json:"credential,omitempty"`
	EnableServerCertificateValidation    *bool                         `json:"enableServerCertificateValidation,omitempty"`
	EncryptedCredential                  *string                       `json:"encryptedCredential,omitempty"`
	Password                             SecretBase                    `json:"password"`
	Resource                             *interface{}                  `json:"resource,omitempty"`
	Scope                                *interface{}                  `json:"scope,omitempty"`
	ServicePrincipalCredentialType       *interface{}                  `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalEmbeddedCert         SecretBase                    `json:"servicePrincipalEmbeddedCert"`
	ServicePrincipalEmbeddedCertPassword SecretBase                    `json:"servicePrincipalEmbeddedCertPassword"`
	ServicePrincipalId                   *interface{}                  `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey                  SecretBase                    `json:"servicePrincipalKey"`
	Tenant                               *interface{}                  `json:"tenant,omitempty"`
	TokenEndpoint                        *interface{}                  `json:"tokenEndpoint,omitempty"`
	Url                                  interface{}                   `json:"url"`
	UserName                             *interface{}                  `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &RestServiceLinkedServiceTypeProperties{}

func (s *RestServiceLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AadResourceId                     *interface{}                  `json:"aadResourceId,omitempty"`
		AuthHeaders                       *interface{}                  `json:"authHeaders,omitempty"`
		AuthenticationType                RestServiceAuthenticationType `json:"authenticationType"`
		AzureCloudType                    *interface{}                  `json:"azureCloudType,omitempty"`
		ClientId                          *interface{}                  `json:"clientId,omitempty"`
		Credential                        *CredentialReference          `json:"credential,omitempty"`
		EnableServerCertificateValidation *bool                         `json:"enableServerCertificateValidation,omitempty"`
		EncryptedCredential               *string                       `json:"encryptedCredential,omitempty"`
		Resource                          *interface{}                  `json:"resource,omitempty"`
		Scope                             *interface{}                  `json:"scope,omitempty"`
		ServicePrincipalCredentialType    *interface{}                  `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId                *interface{}                  `json:"servicePrincipalId,omitempty"`
		Tenant                            *interface{}                  `json:"tenant,omitempty"`
		TokenEndpoint                     *interface{}                  `json:"tokenEndpoint,omitempty"`
		Url                               interface{}                   `json:"url"`
		UserName                          *interface{}                  `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AadResourceId = decoded.AadResourceId
	s.AuthHeaders = decoded.AuthHeaders
	s.AuthenticationType = decoded.AuthenticationType
	s.AzureCloudType = decoded.AzureCloudType
	s.ClientId = decoded.ClientId
	s.Credential = decoded.Credential
	s.EnableServerCertificateValidation = decoded.EnableServerCertificateValidation
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Resource = decoded.Resource
	s.Scope = decoded.Scope
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.TokenEndpoint = decoded.TokenEndpoint
	s.Url = decoded.Url
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RestServiceLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'RestServiceLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'RestServiceLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["servicePrincipalEmbeddedCert"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalEmbeddedCert' for 'RestServiceLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalEmbeddedCert = impl
	}

	if v, ok := temp["servicePrincipalEmbeddedCertPassword"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalEmbeddedCertPassword' for 'RestServiceLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalEmbeddedCertPassword = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'RestServiceLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
