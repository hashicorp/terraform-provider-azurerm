package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDatabricksLinkedServiceTypeProperties struct {
	AccessToken                 SecretBase           `json:"accessToken"`
	Authentication              *string              `json:"authentication,omitempty"`
	Credential                  *CredentialReference `json:"credential,omitempty"`
	Domain                      string               `json:"domain"`
	EncryptedCredential         *string              `json:"encryptedCredential,omitempty"`
	ExistingClusterId           *string              `json:"existingClusterId,omitempty"`
	InstancePoolId              *string              `json:"instancePoolId,omitempty"`
	NewClusterCustomTags        *map[string]string   `json:"newClusterCustomTags,omitempty"`
	NewClusterDriverNodeType    *string              `json:"newClusterDriverNodeType,omitempty"`
	NewClusterEnableElasticDisk *bool                `json:"newClusterEnableElasticDisk,omitempty"`
	NewClusterInitScripts       *[]string            `json:"newClusterInitScripts,omitempty"`
	NewClusterLogDestination    *string              `json:"newClusterLogDestination,omitempty"`
	NewClusterNodeType          *string              `json:"newClusterNodeType,omitempty"`
	NewClusterNumOfWorker       *string              `json:"newClusterNumOfWorker,omitempty"`
	NewClusterSparkConf         *map[string]string   `json:"newClusterSparkConf,omitempty"`
	NewClusterSparkEnvVars      *map[string]string   `json:"newClusterSparkEnvVars,omitempty"`
	NewClusterVersion           *string              `json:"newClusterVersion,omitempty"`
	PolicyId                    *string              `json:"policyId,omitempty"`
	WorkspaceResourceId         *string              `json:"workspaceResourceId,omitempty"`
}

var _ json.Unmarshaler = &AzureDatabricksLinkedServiceTypeProperties{}

func (s *AzureDatabricksLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Authentication              *string              `json:"authentication,omitempty"`
		Credential                  *CredentialReference `json:"credential,omitempty"`
		Domain                      string               `json:"domain"`
		EncryptedCredential         *string              `json:"encryptedCredential,omitempty"`
		ExistingClusterId           *string              `json:"existingClusterId,omitempty"`
		InstancePoolId              *string              `json:"instancePoolId,omitempty"`
		NewClusterCustomTags        *map[string]string   `json:"newClusterCustomTags,omitempty"`
		NewClusterDriverNodeType    *string              `json:"newClusterDriverNodeType,omitempty"`
		NewClusterEnableElasticDisk *bool                `json:"newClusterEnableElasticDisk,omitempty"`
		NewClusterInitScripts       *[]string            `json:"newClusterInitScripts,omitempty"`
		NewClusterLogDestination    *string              `json:"newClusterLogDestination,omitempty"`
		NewClusterNodeType          *string              `json:"newClusterNodeType,omitempty"`
		NewClusterNumOfWorker       *string              `json:"newClusterNumOfWorker,omitempty"`
		NewClusterSparkConf         *map[string]string   `json:"newClusterSparkConf,omitempty"`
		NewClusterSparkEnvVars      *map[string]string   `json:"newClusterSparkEnvVars,omitempty"`
		NewClusterVersion           *string              `json:"newClusterVersion,omitempty"`
		PolicyId                    *string              `json:"policyId,omitempty"`
		WorkspaceResourceId         *string              `json:"workspaceResourceId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Authentication = decoded.Authentication
	s.Credential = decoded.Credential
	s.Domain = decoded.Domain
	s.EncryptedCredential = decoded.EncryptedCredential
	s.ExistingClusterId = decoded.ExistingClusterId
	s.InstancePoolId = decoded.InstancePoolId
	s.NewClusterCustomTags = decoded.NewClusterCustomTags
	s.NewClusterDriverNodeType = decoded.NewClusterDriverNodeType
	s.NewClusterEnableElasticDisk = decoded.NewClusterEnableElasticDisk
	s.NewClusterInitScripts = decoded.NewClusterInitScripts
	s.NewClusterLogDestination = decoded.NewClusterLogDestination
	s.NewClusterNodeType = decoded.NewClusterNodeType
	s.NewClusterNumOfWorker = decoded.NewClusterNumOfWorker
	s.NewClusterSparkConf = decoded.NewClusterSparkConf
	s.NewClusterSparkEnvVars = decoded.NewClusterSparkEnvVars
	s.NewClusterVersion = decoded.NewClusterVersion
	s.PolicyId = decoded.PolicyId
	s.WorkspaceResourceId = decoded.WorkspaceResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDatabricksLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accessToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessToken' for 'AzureDatabricksLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessToken = impl
	}

	return nil
}
