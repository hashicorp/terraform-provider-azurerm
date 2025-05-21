package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDatabricksLinkedServiceTypeProperties struct {
	AccessToken                 SecretBase              `json:"accessToken"`
	Authentication              *interface{}            `json:"authentication,omitempty"`
	Credential                  *CredentialReference    `json:"credential,omitempty"`
	Domain                      interface{}             `json:"domain"`
	EncryptedCredential         *string                 `json:"encryptedCredential,omitempty"`
	ExistingClusterId           *interface{}            `json:"existingClusterId,omitempty"`
	InstancePoolId              *interface{}            `json:"instancePoolId,omitempty"`
	NewClusterCustomTags        *map[string]interface{} `json:"newClusterCustomTags,omitempty"`
	NewClusterDriverNodeType    *interface{}            `json:"newClusterDriverNodeType,omitempty"`
	NewClusterEnableElasticDisk *bool                   `json:"newClusterEnableElasticDisk,omitempty"`
	NewClusterInitScripts       *[]string               `json:"newClusterInitScripts,omitempty"`
	NewClusterLogDestination    *interface{}            `json:"newClusterLogDestination,omitempty"`
	NewClusterNodeType          *interface{}            `json:"newClusterNodeType,omitempty"`
	NewClusterNumOfWorker       *interface{}            `json:"newClusterNumOfWorker,omitempty"`
	NewClusterSparkConf         *map[string]interface{} `json:"newClusterSparkConf,omitempty"`
	NewClusterSparkEnvVars      *map[string]interface{} `json:"newClusterSparkEnvVars,omitempty"`
	NewClusterVersion           *interface{}            `json:"newClusterVersion,omitempty"`
	PolicyId                    *interface{}            `json:"policyId,omitempty"`
	WorkspaceResourceId         *interface{}            `json:"workspaceResourceId,omitempty"`
}

var _ json.Unmarshaler = &AzureDatabricksLinkedServiceTypeProperties{}

func (s *AzureDatabricksLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Authentication              *interface{}            `json:"authentication,omitempty"`
		Credential                  *CredentialReference    `json:"credential,omitempty"`
		Domain                      interface{}             `json:"domain"`
		EncryptedCredential         *string                 `json:"encryptedCredential,omitempty"`
		ExistingClusterId           *interface{}            `json:"existingClusterId,omitempty"`
		InstancePoolId              *interface{}            `json:"instancePoolId,omitempty"`
		NewClusterCustomTags        *map[string]interface{} `json:"newClusterCustomTags,omitempty"`
		NewClusterDriverNodeType    *interface{}            `json:"newClusterDriverNodeType,omitempty"`
		NewClusterEnableElasticDisk *bool                   `json:"newClusterEnableElasticDisk,omitempty"`
		NewClusterInitScripts       *[]string               `json:"newClusterInitScripts,omitempty"`
		NewClusterLogDestination    *interface{}            `json:"newClusterLogDestination,omitempty"`
		NewClusterNodeType          *interface{}            `json:"newClusterNodeType,omitempty"`
		NewClusterNumOfWorker       *interface{}            `json:"newClusterNumOfWorker,omitempty"`
		NewClusterSparkConf         *map[string]interface{} `json:"newClusterSparkConf,omitempty"`
		NewClusterSparkEnvVars      *map[string]interface{} `json:"newClusterSparkEnvVars,omitempty"`
		NewClusterVersion           *interface{}            `json:"newClusterVersion,omitempty"`
		PolicyId                    *interface{}            `json:"policyId,omitempty"`
		WorkspaceResourceId         *interface{}            `json:"workspaceResourceId,omitempty"`
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
