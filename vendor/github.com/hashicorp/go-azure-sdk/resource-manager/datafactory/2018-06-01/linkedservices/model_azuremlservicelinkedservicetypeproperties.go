package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLServiceLinkedServiceTypeProperties struct {
	Authentication      *string    `json:"authentication,omitempty"`
	EncryptedCredential *string    `json:"encryptedCredential,omitempty"`
	MlWorkspaceName     string     `json:"mlWorkspaceName"`
	ResourceGroupName   string     `json:"resourceGroupName"`
	ServicePrincipalId  *string    `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey SecretBase `json:"servicePrincipalKey"`
	SubscriptionId      string     `json:"subscriptionId"`
	Tenant              *string    `json:"tenant,omitempty"`
}

var _ json.Unmarshaler = &AzureMLServiceLinkedServiceTypeProperties{}

func (s *AzureMLServiceLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Authentication      *string `json:"authentication,omitempty"`
		EncryptedCredential *string `json:"encryptedCredential,omitempty"`
		MlWorkspaceName     string  `json:"mlWorkspaceName"`
		ResourceGroupName   string  `json:"resourceGroupName"`
		ServicePrincipalId  *string `json:"servicePrincipalId,omitempty"`
		SubscriptionId      string  `json:"subscriptionId"`
		Tenant              *string `json:"tenant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Authentication = decoded.Authentication
	s.EncryptedCredential = decoded.EncryptedCredential
	s.MlWorkspaceName = decoded.MlWorkspaceName
	s.ResourceGroupName = decoded.ResourceGroupName
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.SubscriptionId = decoded.SubscriptionId
	s.Tenant = decoded.Tenant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureMLServiceLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureMLServiceLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
