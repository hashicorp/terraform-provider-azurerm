package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ODataLinkedServiceTypeProperties struct {
	AadResourceId                        *interface{}                            `json:"aadResourceId,omitempty"`
	AadServicePrincipalCredentialType    *ODataAadServicePrincipalCredentialType `json:"aadServicePrincipalCredentialType,omitempty"`
	AuthHeaders                          *map[string]string                      `json:"authHeaders,omitempty"`
	AuthenticationType                   *ODataAuthenticationType                `json:"authenticationType,omitempty"`
	AzureCloudType                       *interface{}                            `json:"azureCloudType,omitempty"`
	EncryptedCredential                  *string                                 `json:"encryptedCredential,omitempty"`
	Password                             SecretBase                              `json:"password"`
	ServicePrincipalEmbeddedCert         SecretBase                              `json:"servicePrincipalEmbeddedCert"`
	ServicePrincipalEmbeddedCertPassword SecretBase                              `json:"servicePrincipalEmbeddedCertPassword"`
	ServicePrincipalId                   *interface{}                            `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey                  SecretBase                              `json:"servicePrincipalKey"`
	Tenant                               *interface{}                            `json:"tenant,omitempty"`
	Url                                  interface{}                             `json:"url"`
	UserName                             *interface{}                            `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &ODataLinkedServiceTypeProperties{}

func (s *ODataLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AadResourceId                     *interface{}                            `json:"aadResourceId,omitempty"`
		AadServicePrincipalCredentialType *ODataAadServicePrincipalCredentialType `json:"aadServicePrincipalCredentialType,omitempty"`
		AuthHeaders                       *map[string]string                      `json:"authHeaders,omitempty"`
		AuthenticationType                *ODataAuthenticationType                `json:"authenticationType,omitempty"`
		AzureCloudType                    *interface{}                            `json:"azureCloudType,omitempty"`
		EncryptedCredential               *string                                 `json:"encryptedCredential,omitempty"`
		ServicePrincipalId                *interface{}                            `json:"servicePrincipalId,omitempty"`
		Tenant                            *interface{}                            `json:"tenant,omitempty"`
		Url                               interface{}                             `json:"url"`
		UserName                          *interface{}                            `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AadResourceId = decoded.AadResourceId
	s.AadServicePrincipalCredentialType = decoded.AadServicePrincipalCredentialType
	s.AuthHeaders = decoded.AuthHeaders
	s.AuthenticationType = decoded.AuthenticationType
	s.AzureCloudType = decoded.AzureCloudType
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.Url = decoded.Url
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ODataLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'ODataLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["servicePrincipalEmbeddedCert"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalEmbeddedCert' for 'ODataLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalEmbeddedCert = impl
	}

	if v, ok := temp["servicePrincipalEmbeddedCertPassword"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalEmbeddedCertPassword' for 'ODataLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalEmbeddedCertPassword = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'ODataLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
